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

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/fsnotify/fsnotify"
	jobs_v1 "github.com/videocoin/cloud-api/jobs/v1"
	manager_v1 "github.com/videocoin/cloud-api/manager/v1"
	verifier_v1 "github.com/videocoin/cloud-api/verifier/v1"

	"github.com/grafov/m3u8"
)

// SyncDir watches file system and processes chunks as they are written
func (s *Service) syncDir(stop chan struct{}, cmd *exec.Cmd, job *jobs_v1.Job, dir string, bitrate uint32) {
	var ch = make(chan Task, 10)
	go s.process(ch, job)

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

	walletHex := common.HexToAddress(job.StreamAddress)

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
					ch <- Task{
						Id:              job.Id,
						ChunksDir:       dir,
						InputChunkName:  chunk,
						Bitrate:         bitrate,
						Playlist:        playlist,
						OutputID:        chunkNum,
						InputID:         chunkNum,
						OutputChunkName: fmt.Sprintf("%d.ts", chunkNum),
						Wallet:          walletHex,
						StreamID:        big.NewInt(job.StreamId),
						StreamAddress:   job.StreamAddress,
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
				err = watcher.Remove(dir)
				if err != nil {
					s.log.Errorf("failed to remove dir [ %s ]: %s", dir, err.Error())
				}
				watcher.Close()
				close(ch)

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
func (s *Service) handleChunk(task *Task) error {
	chunkLoc := path.Join(task.ChunksDir, task.InputChunkName)
	uploadPath := fmt.Sprintf("%s/%d", task.Id, task.Bitrate)

	if task.InputChunkName == "0.ts" {
		duration, err := s.duration(chunkLoc)
		if err != nil {
			duration = 2.0
		}

		task.Playlist.TargetDuration = duration
	}

	duration, err := s.duration(chunkLoc)
	if err != nil {
		return err
	}

	if err = task.Playlist.Append(task.OutputChunkName, duration, ""); err != nil {
		return err
	}

	chunk, err := os.Open(chunkLoc)
	if err != nil {
		return err
	}

	// Upload chunk
	if err = s.upload(path.Join(uploadPath, task.OutputChunkName), chunk); err != nil {
		return err
	}

	// Upload playlist
	if err = s.upload(path.Join(uploadPath, "index.m3u8"), task.Playlist.Encode()); err != nil {
		return err
	}

	addInputTx, err := s.streamManager.AddInputChunkId(s.bcAuth, task.StreamID, task.InputID)
	if err != nil {
		return err
	}

	go func() {
		s.log.Infof("waiting for AddInputChunkTX: [ %x ] to be mined", addInputTx.Hash())
		rcp, err := bind.WaitMined(context.Background(), nil, addInputTx)
		if err != nil {
			s.log.Errorf("failed to wait for tx: %s", err.Error())
		} else {
			s.log.Infof("WaitBind() receipt: %v", rcp)
		}
	}()

	tx, err := s.submitProof(task.StreamAddress, task.Bitrate, task.InputID, task.OutputID)
	if err != nil {
		return err
	}

	inputChunk := fmt.Sprintf("%s/%s/%s", s.cfg.BaseStreamURL, task.Id, task.InputChunkName)
	outputChunk := fmt.Sprintf("https://%s/%s/%d/%s", s.cfg.Bucket, task.Id, task.Bitrate, task.OutputChunkName)

	go s.verify(tx, task, inputChunk, outputChunk)

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
func (s *Service) verify(tx *types.Transaction, task *Task, localFile, outputURL string) error {
	_, err := s.verifier.Verify(context.Background(), &verifier_v1.VerifyRequest{
		TxHash:         tx.Hash().Hex(),
		StreamId:       task.StreamID.Int64(),
		Bitrate:        task.Bitrate,
		InputId:        task.InputID.Uint64(),
		OutputId:       task.OutputID.Uint64(),
		SourceChunkUrl: localFile,
		ResultChunkUrl: outputURL,
	})

	if err != nil {
		s.log.Errorf("failed to call verifier: %s", err.Error())
	}

	balance, err := s.manager.CheckBalance(context.Background(), &manager_v1.CheckBalanceRequest{ContractAddress: task.StreamAddress})
	if err != nil {
		s.log.Warnf("failed to check balance, allowing work")
	}

	resp, err := s.manager.Get(context.Background(), &manager_v1.JobRequest{Id: task.Id})
	if err != nil {
		s.log.Warnf("failed to get current job status: %s", err.Error())
	}

	if balance.Balance <= 0 || resp.Status == jobs_v1.JobStatusCompleted /* job has been reset */ {
		_ = task.cmd.Process.Kill()
		task.stopChan <- struct{}{}
	}

	return err
}

func (s *Service) process(ch chan Task, job *jobs_v1.Job) {
	for len(ch) < 2 {
		time.Sleep(1 * time.Second)
	}

	s.updateStatus(job.Id, jobs_v1.JobStatusReady)

	for {
		for len(ch) < 2 {
			time.Sleep(500 * time.Millisecond)
		}

		t := <-ch

		go func() {
			if err := s.chunkCreated(&t); err != nil {
				s.log.Errorf("failed to register chunk: %s", err.Error())
			}
		}()

		if err := s.handleChunk(&t); err != nil {
			s.log.Errorf("failed to handle chunk: %s", err.Error())
		}
	}
}

func (s *Service) updateStatus(id string, status jobs_v1.JobStatus) {
	_, err := s.manager.UpdateStatus(s.ctx, &manager_v1.UpdateJobRequest{
		Id:     id,
		Status: status,
	})

	if err != nil {
		s.log.Errorf("failed to update stream status: %s with id: %s", err.Error(), id)
	}
}

func (s *Service) chunkCreated(t *Task) error {
	_, err := s.manager.ChunkCreated(s.ctx, &manager_v1.ChunkCreatedRequest{
		StreamId:      t.StreamID.Int64(),
		SourceChunkId: t.InputID.Int64(),
		ResultChunkId: t.OutputID.Int64(),
		Bitrate:       t.Bitrate,
	})

	return err
}
