package transcode

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	pb "github.com/VideoCoin/common/proto"
	"github.com/VideoCoin/common/stream"
	"github.com/VideoCoin/go-videocoin/common"
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
					!strings.Contains(chunk, ".m3u8") {

					randomID := RandomBigInt()

					s.log.Infof("created file: %s generated name: %d", chunk, randomID)

					jobChan <- Job{
						ChunksDir:       dir,
						InputChunkName:  chunk,
						Bitrate:         bitrate,
						Playlist:        playlist,
						OutputID:        randomID,
						InputID:         ChunkNum(chunk),
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

	err = s.SubmitProof(job.Wallet, job.Bitrate, job.InputID, job.OutputID)
	if err != nil {
		s.log.Errorf("failed to submit proof: %s", err.Error())
		return err
	}

	localFile := fmt.Sprintf("%s/%d-%s/%s", s.cfg.BaseStreamURL, job.StreamID, job.Wallet, job.InputChunkName)
	outputURL := fmt.Sprintf("https://storage.googleapis.com/%s/%d/%d/%s", s.cfg.Bucket, job.StreamID, job.Bitrate, job.OutputChunkName)

	if err = s.VerifyChunk(job.StreamID, localFile, outputURL, job.Bitrate, job.InputID, job.OutputID); err != nil {
		s.log.Errorf("failed to verify chunk: %s", err.Error())
		return err
	}
	return nil
}

// SubmitProof registers work (output chunk)
func (s *Service) SubmitProof(address common.Address, bitrate uint32, inputChunkID *big.Int, outputChunkID *big.Int) error {
	streamInstance, err := stream.NewStream(address, s.bcClient)
	if err != nil {
		s.log.Errorf("failed to create stream instance: %s", err.Error())
		return err
	}

	for i := 0; i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
		_, err = streamInstance.SubmitProof(s.bcAuth, big.NewInt(int64(bitrate)), inputChunkID, big.NewInt(0), outputChunkID)
		if err == nil {
			s.log.Info("success")
			break
		} else if err.Error() == errNonceTooLow {
			s.log.Errorf("nonce too low")
			continue
		} else if err != nil {
			s.log.Errorf("error on submit proof %s", err.Error())
			break
		}
	}

	return nil
}

// VerifyChunk blahg
func (s *Service) VerifyChunk(streamID *big.Int, src string, res string, bitrate uint32, inputID *big.Int, resultID *big.Int) error {
	s.log.Infof("calling verifier with: src: %s \nres: %s \ninput_id: %d \noutput_id: %d \nstream_id: %d \nbitrate: %d", src, res, inputID, resultID, streamID, bitrate)
	resp, err := http.PostForm(s.cfg.VerifierURL+"/api/v1/verify", url.Values{
		"source_chunk_url": {src},
		"result_chunk_url": {res},
		"stream_id":        {fmt.Sprintf("%d", streamID)},
		"bitrate":          {fmt.Sprintf("%d", bitrate)},
		"input_id":         {fmt.Sprintf("%d", inputID)},
		"result_id":        {fmt.Sprintf("%d", resultID)},
	})

	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	s.log.Infof("verifier response: code [ %d ] msg [ %s ]", resp.StatusCode, string(body))

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
			if err := s.chunkCreated(&j); err != nil {
				s.log.Errorf("failed to report chunk created: %s", err.Error())
			}

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
