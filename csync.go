package transcode

import (
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"log"

	pb "github.com/VideoCoin/common/proto"
	"github.com/VideoCoin/common/stream"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/fsnotify/fsnotify"

	"github.com/grafov/m3u8"
)

func (s *Service) getDuration(input string) (float64, error) {
	s.log.Infof("using input %s", input)
	args := []string{"-v", "panic", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", input}
	stdout, err := exec.Command("ffprobe", args...).CombinedOutput()
	if err != nil {
		return 0.0, err
	}

	cleanOut := strings.TrimSpace(string(stdout))

	return strconv.ParseFloat(cleanOut, 64)
}

// SyncDir watches file system and processes chunks as they are written
func (s *Service) SyncDir(workOrder *pb.WorkOrder, dir string, bitrate uint32) {
	var jobChan = make(chan Job, 10)
	go s.process(jobChan)

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

					randomID := randomBigInt()

					s.log.Infof("created file: %s generated name: %d", chunk, randomID)

					jobChan <- Job{
						ChunksDir:       dir,
						InputChunkName:  chunk,
						Bitrate:         bitrate,
						Playlist:        playlist,
						OutputID:        randomID,
						InputID:         getChunkNum(chunk),
						OutputChunkName: fmt.Sprintf("%d.ts"),
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
		log.Fatal(err)
	}
	<-done
}

// DoTheDamnThing Appends to playlist, generates chunk id, calls verifier, uploads result
func (s *Service) handleChunk(job *Job) error {

	chunkLoc := path.Join(job.ChunksDir, job.InputChunkName)
	uploadPath := fmt.Sprintf("%d/%d", job.StreamID, job.Bitrate)
	if job.InputChunkName == "0.ts" {
		duration, err := s.getDuration(chunkLoc)
		if err != nil {
			s.log.Warnf("failed to get chunk duration: %s", err.Error())
			duration = 10.0
		}

		job.Playlist.TargetDuration = duration
	}

	duration, err := s.getDuration(chunkLoc)
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
	if err = s.Upload(path.Join(uploadPath, job.OutputChunkName), chunk); err != nil {
		return err
	}

	// Upload playlist
	if err = s.Upload(path.Join(uploadPath, "index.m3u8"), job.Playlist.Encode()); err != nil {
		return err
	}

	err = s.SubmitProof(job.Wallet, job.Bitrate, job.InputID, job.OutputID)
	if err != nil {
		return err
	}

	s.addNonce()

	localFile := fmt.Sprintf("%s/%d-%s/%s", s.cfg.BaseStreamURL, job.StreamID, job.Wallet, job.InputChunkName)
	outputURL := fmt.Sprintf("https://storage.googleapis.com/%s/%d/%s/%s", s.cfg.Bucket, job.StreamID, job.ChunksDir, job.OutputChunkName)

	return s.VerifyChunk(job.StreamID, localFile, outputURL, job.Bitrate)

}

// SubmitProof registers work (output chunk)
func (s *Service) SubmitProof(address common.Address, bitrate uint32, inputChunkID *big.Int, outputChunkID *big.Int) error {
	streamInstance, err := stream.NewStream(address, s.bcClient)
	if err != nil {
		return err
	}

	_, err = streamInstance.SubmitProof(s.bcAuth, big.NewInt(int64(bitrate)), inputChunkID, big.NewInt(0), outputChunkID)
	if err != nil {
		return err
	}

	return nil
}

// VerifyChunk blahg
func (s *Service) VerifyChunk(streamID *big.Int, src string, res string, bitrate uint32) error {
	form := url.Values{}
	form.Add("source_chunk_url", src)
	form.Add("result_chunk_url", res)
	form.Add("stream_id", fmt.Sprintf("%d", streamID))

	resp, err := http.PostForm(s.cfg.VerifierURL+"/api/v1/verify", form)
	if err != nil {
		return err
	}

	s.log.Infof("verifier response: code [ %d ]", resp.StatusCode)

	return nil
}

func (s *Service) process(jobChan chan Job) {
	for len(jobChan) < 2 {
		sleep()
	}

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

func (s *Service) chunkCreated(j *Job) error {
	_, err := s.manager.ChunkCreated(s.ctx, &pb.ChunkCreatedRequest{
		StreamId:      j.StreamID.Int64(),
		SourceChunkId: j.InputID.Int64(),
		ResultChunkId: j.OutputID.Int64(),
		Bitrate:       j.Bitrate,
	})

	return err
}

// AddNonce increment nonce by one, required for every blockcain interaction
func (s *Service) addNonce() {
	newNonce, err := s.bcClient.PendingNonceAt(s.ctx, s.pkAddr)
	if err != nil {
		s.log.Errorf("failed to increase nonce: %s", err.Error())
		return
	}
	s.bcAuth.Nonce = big.NewInt(int64(newNonce))
}
