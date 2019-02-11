package transcode

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"path"
	"strings"
	"time"

	pb "github.com/VideoCoin/common/proto"
	"github.com/VideoCoin/common/stream"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/VideoCoin/go-videocoin/core/types"
	"github.com/fsnotify/fsnotify"

	"github.com/grafov/m3u8"
)

// SyncDir watches file system and processes chunks as they are written
func (s *Service) SyncDir(workOrder *pb.WorkOrder, dir string, bitrate uint32) {
	var jobChan = make(chan Job, 10)
	go s.process(jobChan, workOrder)

	playlist, err := m3u8.NewMediaPlaylist(10000, 10000)
	if err != nil {
		s.log.Errorf("failed to generate new media playlist: %s", err.Error())
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		s.log.Errorf("failed to make watcher: %s", err.Error())
		return
	}

	defer watcher.Close()

	done := make(chan bool)

	walletHex := common.HexToAddress(workOrder.WalletAddress)

	written := make(map[string]bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				chunk := path.Base(event.Name)

				if (event.Op&fsnotify.Create == fsnotify.Create) &&
					!strings.Contains(chunk, "tmp") &&
					!strings.Contains(chunk, ".m3u8") &&
					!written[chunk] {

					balance, err := s.manager.CheckBalance(context.Background(), &pb.CheckBalanceRequest{ContractAddress: workOrder.ContractAddress})
					if err != nil {
						s.log.Warnf("failed to check balance, allowing work")
					}
					if balance.Balance <= 0 {
						return
					}

					written[chunk] = true
					chunkNum := getChunkNum(chunk)

					_, err = s.streamManager.AddInputChunkId(s.bcAuth, big.NewInt(workOrder.StreamId), chunkNum)
					if err != nil {
						s.log.Errorf("failed to add input chunk: [ %d ] bitrate: [ %d ] err: [ %s ]", chunkNum, bitrate, err.Error())
					}

					randomID := RandomBigInt()

					s.log.Infof("created file: %s generated name: %d", chunk, randomID)

					jobChan <- Job{
						ChunksDir:       dir,
						InputChunkName:  chunk,
						Bitrate:         bitrate,
						Playlist:        playlist,
						OutputID:        randomID,
						InputID:         chunkNum,
						OutputChunkName: fmt.Sprintf("%d.ts", randomID),
						Wallet:          walletHex,
						StreamID:        big.NewInt(workOrder.StreamId),
					}

				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				if err != nil {
					s.log.Errorf("event watcher error: %s", err.Error())
				}
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		s.log.Fatalf("water failure: %s", err.Error())
	}

	<-done
}

// DoTheDamnThing Appends to playlist, generates chunk id, calls verifier, uploads result
func (s *Service) handleChunk(job *Job) error {

	chunkLoc := path.Join(job.ChunksDir, job.InputChunkName)
	uploadPath := fmt.Sprintf("%d/%d", job.StreamID, job.Bitrate)

	if job.InputChunkName == "0.ts" {
		duration, err := s.Duration(chunkLoc)
		if err != nil {
			s.log.Warnf("failed to get chunk duration: %s", err.Error())
			duration = 10.0
		}

		job.Playlist.TargetDuration = duration
	}

	duration, err := s.Duration(chunkLoc)
	if err != nil {
		s.log.Errorf("failed to get chunk duration: %s", err.Error())
		return err
	}

	if err = job.Playlist.Append(job.OutputChunkName, duration, ""); err != nil {
		s.log.Errorf("failed to append to playlist: %s", err.Error())
		return err
	}

	chunk, err := os.Open(chunkLoc)
	if err != nil {
		s.log.Errorf("failed to open chunk: %s", err.Error())
		return err
	}

	// Upload chunk
	if err = s.Upload(path.Join(uploadPath, job.OutputChunkName), chunk); err != nil {
		s.log.Errorf("failed to upload chunk: %s", err.Error())
		return err
	}

	// Upload playlist
	if err = s.Upload(path.Join(uploadPath, "index.m3u8"), job.Playlist.Encode()); err != nil {
		s.log.Errorf("failed to upload playlist: %s", err.Error())
		return err
	}

	tx, err := s.SubmitProof(job.Wallet, job.Bitrate, job.InputID, job.OutputID)
	if err != nil {
		s.log.Errorf("failed to submit proof: %s", err.Error())
		return err
	}

	localFile := fmt.Sprintf("%s/%d-%x/%s", s.cfg.BaseStreamURL, job.StreamID, job.Wallet, job.InputChunkName)
	outputURL := fmt.Sprintf("https://storage.googleapis.com/%s/%d/%d/%s", s.cfg.Bucket, job.StreamID, job.Bitrate, job.OutputChunkName)

	if err = s.VerifyChunk(tx, job.StreamID, localFile, outputURL, job.Bitrate, job.InputID, job.OutputID); err != nil {
		s.log.Errorf("failed to verify chunk: %s", err.Error())
		return err
	}
	return nil
}

// SubmitProof registers work (output chunk)
func (s *Service) SubmitProof(address common.Address, bitrate uint32, inputChunkID *big.Int, outputChunkID *big.Int) (*types.Transaction, error) {
	streamInstance, err := stream.NewStream(address, s.bcClient)
	if err != nil {
		s.log.Errorf("failed to create stream instance: %s", err.Error())
		return nil, err
	}

	s.log.Infof("submitting proof: addr: %x\nbitrate: %d\ninput_id: %d\noutput_id: %d", address.Hex(), bitrate, inputChunkID, outputChunkID)

	tx, err := streamInstance.SubmitProof(s.bcAuth, big.NewInt(int64(bitrate)), inputChunkID, big.NewInt(0), outputChunkID)
	if err != nil {
		s.log.Errorf("error on submit proof %s", err.Error())
		return nil, err
	}

	return tx, nil
}

// VerifyChunk blahg
func (s *Service) VerifyChunk(tx *types.Transaction, streamID *big.Int, src string, res string, bitrate uint32, inputID *big.Int, resultID *big.Int) error {
	s.log.Infof("calling verifier with: src: %s \nres: %s \ninput_id: %d \noutput_id: %d \nstream_id: %d \nbitrate: %d", src, res, inputID, resultID, streamID, bitrate)

	_, err := s.verifier.Verify(context.Background(), &pb.VerifyRequest{
		TxHash:         tx.Hash().Hex(),
		StreamId:       streamID.Int64(),
		Bitrate:        bitrate,
		InputId:        inputID.Int64(),
		OutputId:       resultID.Int64(),
		SourceChunkUrl: src,
		ResultChunkUrl: res,
	})

	if err != nil {
		s.log.Errorf("failed to call verify: %s", err.Error())
	}

	return nil
}

func (s *Service) process(jobChan chan Job, workOrder *pb.WorkOrder) {

	s.updateStatus(workOrder.StreamId, pb.WorkOrderStatusTranscoding.String())

	for len(jobChan) < 2 {
		time.Sleep(1 * time.Second)
	}

	s.updateStatus(workOrder.StreamId, pb.WorkOrderStatusReady.String())

	for {
		select {
		case j := <-jobChan:
			if err := s.handleChunk(&j); err != nil {
				s.log.Errorf("failed to handle chunk: %s", err.Error())
			}
		}
	}
}

func (s *Service) updateStatus(streamID int64, status string) {
	_, err := s.manager.UpdateStreamStatus(s.ctx, &pb.UpdateStreamStatusRequest{
		StreamId: streamID,
		Status:   status,
	})

	if err != nil {
		s.log.Errorf("failed to update stream status: %s", err.Error())
	}
}

func (s *Service) chunkCreated(j *Job) error {
	_, err := s.manager.ChunkCreated(s.ctx, &pb.ChunkCreatedRequest{
		StreamId:      j.StreamID.Int64(),
		SourceChunkId: j.InputID.Int64(),
		ResultChunkId: j.OutputID.Int64(),
		Bitrate:       j.Bitrate,
	})

	return err
}
