package transcode

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	pb "github.com/videocoin/common/proto"
	"golang.org/x/oauth2/google"

	storage "google.golang.org/api/storage/v1"

	"github.com/grafov/m3u8"
)

// CSync syncer struct

// CSyncInit Returns initialized csync object
func CSyncInit(cfg *Config) *CSync {
	return &CSync{
		log: log.WithField("name", "csync"),
		cfg: cfg,
	}
}

func (c *CSync) getDuration(input string) (float64, error) {
	c.log.Infof("using input %s", input)
	args := []string{"-v", "panic", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", input}
	stdout, err := exec.Command("ffprobe", args...).CombinedOutput()
	if err != nil {
		c.log.Warnf("ffprobe output: %s", string(stdout))
		return 0.0, err
	}

	cleanOut := strings.TrimSpace(string(stdout))

	return strconv.ParseFloat(cleanOut, 64)
}

// SyncDir watches file system and processes chunks as they are written
func (c *CSync) SyncDir(workOrder *pb.WorkOrder, chunkDir string) {
	//create playlist
	// wait for chunk
	// get chunk dir
	// append to playlist
	// upload chunk
	// upload playlist

	fullPath := path.Join(c.cfg.OutputDir, workOrder.StreamHash, chunkDir)

	var q = new(JobQueue)

	playlist, err := m3u8.NewMediaPlaylist(10000, 10000)
	if err != nil {
		c.log.Errorf("failed to generate playlist: %s", err.Error())
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
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
					c.log.Println("created file:", chunk)
					q.Push(Job{ChunkName: chunk, ChunksDir: chunkDir, Playlist: playlist})
					c.Work(workOrder, q)

				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				c.log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(fullPath)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

// Work execute jobs only if at least two are in queue
// This prevents accidently working a chunk that ffmpeg has not finished writing yet
func (c *CSync) Work(workOrder *pb.WorkOrder, jobs *JobQueue) {
	if jobs.Len() >= 3 {
		job := jobs.Pop()
		c.DoTheDamnThing(workOrder, &job)
	}
}

// DoTheDamnThing Appends to playlist, generates chunk id, calls verifier, uploads result
func (c *CSync) DoTheDamnThing(workOrder *pb.WorkOrder, job *Job) error {

	var b = make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return err
	}

	chunkLoc := path.Join(c.cfg.OutputDir, workOrder.StreamHash, job.ChunksDir, job.ChunkName)
	uploadPath := path.Join(workOrder.StreamHash, job.ChunksDir)
	if job.ChunkName == "0.ts" {
		duration, err := c.getDuration(chunkLoc)
		if err != nil {
			c.log.Errorf("failed to get duration: %s from: %s ", err.Error(), chunkLoc)
		}

		job.Playlist.TargetDuration = duration
	}

	duration, err := c.getDuration(chunkLoc)
	if err != nil {
		return err
	}

	newChunkName := fmt.Sprintf("%x.ts", b)

	if err = job.Playlist.Append(newChunkName, duration, ""); err != nil {
		return err
	}

	chunk, err := os.Open(chunkLoc)
	if err != nil {
		c.log.Errorf("failed to open chunk: %s", err.Error())
	}

	// Upload chunk
	if err = c.Upload(path.Join(uploadPath, newChunkName), chunk); err != nil {
		return err
	}

	// Upload playlist
	if err = c.Upload(path.Join(uploadPath, "index.m3u8"), job.Playlist.Encode()); err != nil {
		return err
	}

	c.VerifyChunk(workOrder.Id, fmt.Sprintf("%s/%s-%s/%s", c.cfg.BaseStreamURL, workOrder.UserId, workOrder.ApplicationId, job.ChunkName), fmt.Sprintf("https://storage.googleapis.com/%s/%s/%s/%s", c.cfg.Bucket, workOrder.StreamHash, job.ChunksDir, newChunkName))

	return nil
}

// VerifyChunk blahg
func (c *CSync) VerifyChunk(workOrderID uint32, src string, res string) error {
	form := url.Values{}
	form.Add("source_chunk_url", src)
	form.Add("result_chunk_url", res)
	form.Add("job_id", fmt.Sprintf("%d", workOrderID))

	resp, err := http.PostForm(c.cfg.VerifierURL+"/api/v1/verify", form)
	if err != nil {
		return err
	}

	c.log.Infof("verifier response: code [ %d ]", resp.StatusCode)

	return nil
}

// Upload uploads an object to gcs with publicread acl
func (c *CSync) Upload(output string, r io.Reader) error {
	c.log.Infof("uploading chunk: %s to %s: ", output, c.cfg.Bucket)
	client, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		c.log.Errorf("failed to create client: %s", err.Error())
		return err
	}

	svc, err := storage.New(client)
	if err != nil {
		c.log.Errorf("failed to create new service: %s", err.Error())
		return err
	}

	object := &storage.Object{
		Name:         output,
		CacheControl: "public, max-age=0",
	}

	if _, err := svc.Objects.Insert(c.cfg.Bucket, object).Media(r).PredefinedAcl("publicread").Do(); err != nil {
		c.log.Warnf("failed to upload object: %s", err.Error())
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
