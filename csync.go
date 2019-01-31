package transcode

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"log"

	"github.com/VideoCoin/common/handle"
	pb "github.com/VideoCoin/common/proto"
	"github.com/VideoCoin/common/stream"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"

	"github.com/grafov/m3u8"
)

func (s *Service) getDuration(input string) (float64, error) {
	log.Printf("using input %s", input)
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
	//create playlist
	// wait for chunk
	// get chunk dir
	// append to playlist
	// upload chunk
	// upload playlist

	var q = new(JobQueue)

	playlist, err := m3u8.NewMediaPlaylist(10000, 10000)
	if err != nil {
		handle.Err(err)
		return
	}

	watcher, err := fsnotify.NewWatcher()
	handle.Fatal(err)

	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				chunk := path.Base(event.Name)

				if (event.Op&fsnotify.Create == fsnotify.Create) && !strings.Contains(chunk, "tmp") && !strings.Contains(chunk, ".m3u8") {
					log.Println("created file:", chunk)
					q.Push(Job{ChunkName: chunk, ChunksDir: dir, Playlist: playlist, Bitrate: bitrate})
					s.Work(workOrder, q)

				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				handle.Err(err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

// Work execute jobs only if at least two are in queue
// This prevents accidently working a chunk that ffmpeg has not finished writing yet
func (s *Service) Work(workOrder *pb.WorkOrder, jobs *JobQueue) {
	if jobs.Len() >= 3 {
		s.manager.ChunkCreated(s.ctx, &pb.ChunkCreatedRequest{})
		job := jobs.Pop()
		s.DoTheDamnThing(workOrder, &job)
	}
}

// DoTheDamnThing Appends to playlist, generates chunk id, calls verifier, uploads result
func (s *Service) DoTheDamnThing(workOrder *pb.WorkOrder, job *Job) error {

	var b = make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return err
	}

	chunkLoc := path.Join(s.cfg.OutputDir, fmt.Sprintf("%d", workOrder.StreamId), job.ChunksDir, job.ChunkName)
	uploadPath := path.Join(fmt.Sprintf("%d", workOrder.StreamId), job.ChunksDir)
	if job.ChunkName == "0.ts" {
		duration, err := s.getDuration(chunkLoc)
		handle.Err(err)

		job.Playlist.TargetDuration = duration
	}

	duration, err := s.getDuration(chunkLoc)
	if err != nil {
		return err
	}

	newChunkName := fmt.Sprintf("%x.ts", b)

	if err = job.Playlist.Append(newChunkName, duration, ""); err != nil {
		return err
	}

	chunk, err := os.Open(chunkLoc)
	if err != nil {
		return err
	}

	// Upload chunk
	if err = s.Upload(path.Join(uploadPath, newChunkName), chunk); err != nil {
		return err
	}

	// Upload playlist
	if err = s.Upload(path.Join(uploadPath, "index.m3u8"), job.Playlist.Encode()); err != nil {
		return err
	}

	inputChunkID, ok := new(big.Int).SetString(strings.TrimSuffix(job.ChunkName, ".ts"), 16)
	if !ok {
		return fmt.Errorf("failed to convert chunk to bigint")
	}

	outputChunkID := new(big.Int).SetBytes(b)

	_, err = s.sm.AddInputChunkId(s.bcAuth, big.NewInt(workOrder.StreamId), inputChunkID)
	if err != nil {
		return err
	}

	s.AddNonce()

	walletAddr := common.HexToAddress(workOrder.WalletAddress)

	err = s.SubmitProof(walletAddr, job.Bitrate, inputChunkID, outputChunkID)
	if err != nil {
		return err
	}

	s.AddNonce()

	return s.VerifyChunk(workOrder.Id, fmt.Sprintf("%s/%d-%s/%s", s.cfg.BaseStreamURL, workOrder.StreamId, workOrder.WalletAddress, job.ChunkName), fmt.Sprintf("https://storage.googleapis.com/%s/%d/%s/%s", s.cfg.Bucket, workOrder.StreamId, job.ChunksDir, newChunkName), job.Bitrate)

}

// AddNonce increment nonce by one, required for every blockcain interaction
func (s *Service) AddNonce() {
	newNonce, err := s.bcClient.PendingNonceAt(s.ctx, s.pkAddr)
	handle.Err(err)
	s.bcAuth.Nonce = big.NewInt(int64(newNonce))
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
func (s *Service) VerifyChunk(workOrderID uint32, src string, res string, bitrate uint32) error {
	form := url.Values{}
	form.Add("source_chunk_url", src)
	form.Add("result_chunk_url", res)
	form.Add("job_id", fmt.Sprintf("%d", workOrderID))

	resp, err := http.PostForm(s.cfg.VerifierURL+"/api/v1/verify", form)
	if err != nil {
		return err
	}

	log.Printf("verifier response: code [ %d ]", resp.StatusCode)

	return nil
}

// Upload uploads an object to gcs with publicread acl
func (s *Service) Upload(output string, r io.Reader) error {
	client, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		return err
	}

	svc, err := storage.New(client)
	if err != nil {
		return err
	}

	object := &storage.Object{
		Name:         output,
		CacheControl: "public, max-age=0",
	}

	if _, err := svc.Objects.Insert(s.cfg.Bucket, object).Media(r).PredefinedAcl("publicread").Do(); err != nil {
		return err
	}

	return nil
}

// Pop returns item in FIFO
func (q *JobQueue) Pop() (job Job) {
	if len(q.Jobs) == 0 {
		return Job{}
	}

	job, q.Jobs = q.Jobs[0], q.Jobs[1:]

	return

}

// Push item to end of array
func (q *JobQueue) Push(job Job) {
	q.Jobs = append(q.Jobs, job)
}

// Len returns length of job queue
func (q *JobQueue) Len() int {
	return len(q.Jobs)
}
