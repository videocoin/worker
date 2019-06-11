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

	manager_v1 "github.com/VideoCoin/cloud-api/manager/v1"
	verifier_v1 "github.com/VideoCoin/cloud-api/verifier/v1"
	workorder_v1 "github.com/VideoCoin/cloud-api/workorder/v1"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/VideoCoin/go-videocoin/core/types"
	"github.com/fsnotify/fsnotify"

	"github.com/grafov/m3u8"
)

// SyncDir watches file system and processes chunks as they are written
func (s *Service) syncDir(stop chan struct{}, cmd *exec.Cmd, workOrder *workorder_v1.WorkOrder, dir string, bitrate uint32) {
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

	walletHex := common.HexToAddress(workOrder.StreamAddress)

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

					// Add job to the job channel to be worked on later
					jobChan <- Job{
						ChunksDir:       dir,
						InputChunkName:  chunk,
						Bitrate:         bitrate,
						Playlist:        playlist,
						OutputID:        chunkNum,
						InputID:         chunkNum,
						OutputChunkName: fmt.Sprintf("%d.ts", chunkNum),
						Wallet:          walletHex,
						StreamID:        big.NewInt(workOrder.StreamId),
						StreamAddress:   workOrder.StreamAddress,
						StreamHash:      workOrder.StreamHash,
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
				err = watcher.Remove(dir)
				if err != nil {
					s.log.Errorf("failed to remove dir [ %s ]: %s", dir, err.Error())
				}

			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		s.log.Fatalf("failed to watch directory: %s", err.Error())
	}

	<-done
}

// handleChunk Appends to playlist, generates chunk id, calls verifier, uploads result
func (s *Service) handleChunk(job *Job) error {
	chunkLoc := path.Join(job.ChunksDir, job.InputChunkName)
	uploadPath := fmt.Sprintf("%s/%d", job.StreamHash, job.Bitrate)

	if job.InputChunkName == "0.ts" {
		duration, err := s.duration(chunkLoc)
		if err != nil {
			duration = 2.0
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

	tx, err := s.submitProof(job.StreamAddress, job.Bitrate, job.InputID, job.OutputID)
	if err != nil {
		return err
	}

	inputChunk := fmt.Sprintf("%s/%s/%s", s.cfg.BaseStreamURL, job.StreamHash, job.InputChunkName)
	outputChunk := fmt.Sprintf("https://%s/%s/%d/%s", s.cfg.Bucket, job.StreamHash, job.Bitrate, job.OutputChunkName)

	go s.verify(tx, job, inputChunk, outputChunk)

	return nil
}

// SubmitProof registers work (output chunk)
func (s *Service) submitProof(contractAddress string, bitrate uint32, inputChunkID *big.Int, outputChunkID *big.Int) (*types.Transaction, error) {
	streamInstance, err := s.createStreamInstance(contractAddress)
	if err != nil {
		return nil, err
	}

	return streamInstance.SubmitProof(s.bcAuth, big.NewInt(int64(bitrate)), inputChunkID, big.NewInt(0), outputChunkID)
}

// verifyChunk calls verifier with input and output chunk urls
func (s *Service) verify(tx *types.Transaction, job *Job, localFile, outputURL string) error {
	_, err := s.verifier.Verify(context.Background(), &verifier_v1.VerifyRequest{
		TxHash:         tx.Hash().Hex(),
		StreamId:       job.StreamID.Int64(),
		Bitrate:        job.Bitrate,
		InputId:        job.InputID.Uint64(),
		OutputId:       job.OutputID.Uint64(),
		SourceChunkUrl: localFile,
		ResultChunkUrl: outputURL,
	})

	if err != nil {
		s.log.Errorf("failed to call verifier: %s", err.Error())
	}

	balance, err := s.manager.CheckBalance(context.Background(), &manager_v1.CheckBalanceRequest{ContractAddress: job.StreamAddress})
	if err != nil {
		s.log.Warnf("failed to check balance, allowing work")
	}

	if balance.Balance <= 0 {
		_ = job.cmd.Process.Kill()
		job.stopChan <- struct{}{}
	}

	return err
}

func (s *Service) process(jobChan chan Job, workOrder *workorder_v1.WorkOrder) {
	for len(jobChan) < 2 {
		time.Sleep(1 * time.Second)
	}

	s.updateStatus(workOrder.StreamHash, workorder_v1.WorkOrderStatusReady)

	for {
		for len(jobChan) < 2 {
			time.Sleep(500 * time.Millisecond)
		}

		j := <-jobChan

		go func() {
			if err := s.chunkCreated(&j); err != nil {
				s.log.Errorf("failed to register chunk: %s", err.Error())
			}
		}()

		if err := s.handleChunk(&j); err != nil {
			s.log.Errorf("failed to handle chunk: %s", err.Error())
		}

	}
}

func (s *Service) updateStatus(streamHash string, status workorder_v1.WorkOrderStatus) {
	_, err := s.manager.UpdateStreamStatus(s.ctx, &manager_v1.StreamStatusRequest{
		StreamHash: streamHash,
		Status:     status,
	})

	if err != nil {
		s.log.Errorf("failed to update stream status: %s", err.Error())
	}
}

func (s *Service) chunkCreated(j *Job) error {
	_, err := s.manager.ChunkCreated(s.ctx, &manager_v1.ChunkCreatedRequest{
		StreamId:      j.StreamID.Int64(),
		SourceChunkId: j.InputID.Int64(),
		ResultChunkId: j.OutputID.Int64(),
		Bitrate:       j.Bitrate,
	})

	return err
}
