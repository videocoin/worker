package transcode

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/exec"
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
func (s *Service) syncDir(
	stop chan struct{},
	cmd *exec.Cmd,
	workOrder *pb.WorkOrder,
	dir string,
	bitrate uint32,
) {

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

	done := make(chan struct{})

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

					randomID := randomBigInt(8)

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
						ContractAddr:    workOrder.ContractAddress,
						cmd:             cmd,
						stopChan:        stop,
					}

				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				if err != nil {
					s.log.Errorf("event watcher error: %s", err.Error())
				}

			case <-stop:
				watcher.Close()
				watcher.Remove(dir)
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
		duration, err := s.duration(chunkLoc)
		if err != nil {
			duration = 10.0
		}

		job.Playlist.TargetDuration = duration
	}

	duration, err := s.duration(chunkLoc)
	if err != nil {
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
	if err = s.upload(path.Join(uploadPath, job.OutputChunkName), chunk); err != nil {
		return err
	}

	// Upload playlist
	if err = s.upload(path.Join(uploadPath, "index.m3u8"), job.Playlist.Encode()); err != nil {
		return err
	}

	tx, err := s.submitProof(job.Bitrate, job.InputID, job.OutputID)
	if err != nil {
		return err
	}

	inputChunk := fmt.Sprintf("%s/%d-%x/%s", s.cfg.BaseStreamURL, job.StreamID, job.Wallet, job.InputChunkName)
	outputChunk := fmt.Sprintf("https://storage.googleapis.com/%s/%d/%d/%s", s.cfg.Bucket, job.StreamID, job.Bitrate, job.OutputChunkName)

	if err = s.verify(tx, job, inputChunk, outputChunk); err != nil {
		return err
	}
	return nil
}

// SubmitProof registers work (output chunk)
func (s *Service) submitProof(
	bitrate uint32,
	inputChunkID *big.Int,
	outputChunkID *big.Int,
) (*types.Transaction, error) {
	tx, err := s.streamInstance.SubmitProof(s.bcAuth, big.NewInt(int64(bitrate)), inputChunkID, big.NewInt(0), outputChunkID)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// verifyChunk calls verifier with input and output chunk urls
func (s *Service) verify(
	tx *types.Transaction,
	job *Job,
	localFile string,
	outputURL string,
) error {

	s.log.Infof("calling verifier with: src: %s \nres: %s \ninput_id: %d \noutput_id: %d \nstream_id: %d \nbitrate: %d", localFile, outputURL, job.InputID, job.OutputID, job.StreamID, job.Bitrate)

	_, err := s.verifier.Verify(context.Background(), &pb.VerifyRequest{
		TxHash:         tx.Hash().Hex(),
		StreamId:       job.StreamID.Int64(),
		Bitrate:        job.Bitrate,
		InputId:        job.InputID.Uint64(),
		OutputId:       job.OutputID.Uint64(),
		SourceChunkUrl: localFile,
		ResultChunkUrl: outputURL,
	})

	balance, err := s.manager.CheckBalance(context.Background(), &pb.CheckBalanceRequest{ContractAddress: job.ContractAddr})
	if err != nil {
		s.log.Warnf("failed to check balance, allowing work")
	}

	s.log.Infof("current balance at address %s is %f", job.ContractAddr, balance.Balance)

	if balance.Balance <= 0 {
		job.cmd.Process.Kill()
		job.stopChan <- struct{}{}
	}

	return err
}

func (s *Service) process(
	jobChan chan Job,
	workOrder *pb.WorkOrder,
) {
	s.updateStatus(workOrder.StreamId, pb.WorkOrderStatusTranscoding.String())

	for len(jobChan) < 2 {
		time.Sleep(1 * time.Second)
	}

	s.updateStatus(workOrder.StreamId, pb.WorkOrderStatusReady.String())

	for {
		for len(jobChan) < 2 {
			time.Sleep(500 * time.Millisecond)
		}
		select {
		case j := <-jobChan:
			if err := s.handleChunk(&j); err != nil {
				s.log.Errorf("failed to handle chunk: %s", err.Error())
			}
		}
	}
}

func (s *Service) updateStatus(
	streamID int64,
	status string,
) {
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
