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
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/VideoCoin/go-videocoin/core/types"
	"github.com/fsnotify/fsnotify"

	"github.com/grafov/m3u8"
)

// SyncDir watches file system and processes chunks as they are written
func (s *Service) SyncDir(stop chan bool, workOrder *pb.WorkOrder, dir string, bitrate uint32) {
	var jobChan = make(chan Job, 10)
	go s.process(jobChan, workOrder)

	playlist, err := m3u8.NewMediaPlaylist(100, 200)
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

				if event.Op&fsnotify.Create == fsnotify.Create &&
					!strings.Contains(chunk, "tmp") &&
					!strings.Contains(chunk, ".m3u8") &&
					!written[chunk] {

					written[chunk] = true
					chunkNum := getChunkNum(chunk)

					_, err = s.streamManager.AddInputChunkId(s.bcAuth, big.NewInt(workOrder.StreamId), chunkNum)
					if err != nil {
						s.log.Errorf("failed to add input chunk: [ %d ] bitrate: [ %d ] err: [ %s ]", chunkNum, bitrate, err.Error())
					}

					randomID := RandomBigInt()

					// Add job to the job channel to be worked on later
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

			case abort := <-stop:
				if abort {
					watcher.Close()
					watcher.Remove(dir)
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

// handleChunk Appends to playlist, generates chunk id, calls verifier, uploads result
func (s *Service) handleChunk(job *Job) error {
	chunkLoc := path.Join(job.ChunksDir, job.InputChunkName)
	uploadPath := fmt.Sprintf("%d/%d", job.StreamID, job.Bitrate)

	if job.InputChunkName == "0.ts" {
		duration, err := s.Duration(chunkLoc)
		if err != nil {
			duration = 10.0
		}

		job.Playlist.TargetDuration = duration
	}

	duration, err := s.Duration(chunkLoc)
	if err != nil {
		s.log.Warnf("failed to get duration: %s location: %s", err.Error(), chunkLoc)
		return err
	}

	if err = job.Playlist.Append(job.OutputChunkName, duration, ""); err != nil {
		return err
	}

	chunk, err := os.Open(chunkLoc)
	if err != nil {
		return err
	}

	// Upload chunk
	if err = s.Upload(path.Join(uploadPath, job.OutputChunkName), chunk); err != nil {
		return err
	}

	// Upload playlist
	if err = s.Upload(path.Join(uploadPath, "index.m3u8"), job.Playlist.Encode()); err != nil {
		return err
	}

	tx, err := s.submitProof(job.Bitrate, job.InputID, job.OutputID)
	if err != nil {
		return err
	}

	localFile := fmt.Sprintf("%s/%d-%x/%s", s.cfg.BaseStreamURL, job.StreamID, job.Wallet, job.InputChunkName)
	outputURL := fmt.Sprintf("https://storage.googleapis.com/%s/%d/%d/%s", s.cfg.Bucket, job.StreamID, job.Bitrate, job.OutputChunkName)

	if err = s.VerifyChunk(tx, job.StreamID, localFile, outputURL, job.Bitrate, job.InputID, job.OutputID); err != nil {
		return err
	}
	return nil
}

// SubmitProof registers work (output chunk)
func (s *Service) submitProof(bitrate uint32, inputChunkID *big.Int, outputChunkID *big.Int) (*types.Transaction, error) {
	tx, err := s.streamInstance.SubmitProof(s.bcAuth, big.NewInt(int64(bitrate)), inputChunkID, big.NewInt(0), outputChunkID)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// VerifyChunk blahg
func (s *Service) VerifyChunk(tx *types.Transaction, streamID *big.Int, src string, res string, bitrate uint32, inputID *big.Int, resultID *big.Int) error {
	s.log.Infof("calling verifier with: src: %s \nres: %s \ninput_id: %d \noutput_id: %d \nstream_id: %d \nbitrate: %d", src, res, inputID, resultID, streamID, bitrate)

	_, err := s.verifier.Verify(context.Background(), &pb.VerifyRequest{
		TxHash:         tx.Hash().Hex(),
		StreamId:       streamID.Uint64(),
		Bitrate:        bitrate,
		InputId:        inputID.Uint64(),
		OutputId:       resultID.Uint64(),
		SourceChunkUrl: src,
		ResultChunkUrl: res,
	})

	return err
}

func (s *Service) process(jobChan chan Job, workOrder *pb.WorkOrder) {
	s.updateStatus(workOrder.StreamId, pb.WorkOrderStatusTranscoding.String())

	for len(jobChan) < 2 {
		time.Sleep(1 * time.Second)
	}

	s.updateStatus(workOrder.StreamId, pb.WorkOrderStatusReady.String())

	for {
		for len(jobChan) < 2 {
			time.Sleep(1 * time.Second)
		}
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
