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

func getDuration(input string) (float64, error) {
	args := []string{"-v", "panic", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", input}
	stdout, err := exec.Command("ffprobe", args...).CombinedOutput()
	if err != nil {
		return 0.0, err
	}

	cleanOut := strings.TrimSpace(string(stdout))

	return strconv.ParseFloat(cleanOut, 64)
}

// SyncDir watches file system and processes chunks as they are written
func (c *CSync) SyncDir(userID string, appID string, workOrderID uint32, bucket string, folder string, streamHash string, inputURL string) {
	//create playlist
	// wait for chunk
	// get chunk dir
	// append to playlist
	// upload chunk
	// upload playlist

	var q = new(JobQueue)

	playlist, err := m3u8.NewMediaPlaylist(0, 0)
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

				if chunk == "0.ts" {
					duration, err := getDuration(path.Join(folder, event.Name))
					if err != nil {
						c.log.Errorf("failed to get duration: %s", err.Error())
					}

					playlist.TargetDuration = duration
				}

				if (event.Op&fsnotify.Create == fsnotify.Create) && !strings.Contains(chunk, "tmp") && !strings.Contains(chunk, ".m3u8") {
					c.log.Println("created file:", chunk)
					q.Push(Job{ChunkName: chunk, Folder: folder, Playlist: playlist})
					c.Work(userID, appID, workOrderID, q)

				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				c.log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(folder)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

// Work execute jobs only if at least two are in queue
// This prevents accidently working a chunk that ffmpeg has not finished writing yet
func (c *CSync) Work(userID string, appID string, workOrderID uint32, jobs *JobQueue) {
	if jobs.Len() >= 2 {
		job := jobs.Pop()
		c.DoTheDamnThing(userID, appID, workOrderID, job.ChunkName, job.Folder, job.Playlist)
	}
}

// DoTheDamnThing Appends to playlist, generates chunk id, calls verifier, uploads result
func (c *CSync) DoTheDamnThing(
	userID string,
	appID string,
	workOrderID uint32,
	chunkname string,
	folder string,
	playlist *m3u8.MediaPlaylist) error {

	var b = make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return err
	}

	duration, err := getDuration(path.Join(cfg.OutputDir, folder, chunkname))
	if err != nil {
		return err
	}

	newChunkName := fmt.Sprintf("%x.ts", b)

	if err = playlist.Append(newChunkName, duration, ""); err != nil {
		return err
	}

	chunk, err := os.Open(path.Join(cfg.OutputDir, folder, chunkname))
	if err != nil {
		c.log.Errorf("failed to open chunk: %s", err.Error())
	}

	// Upload chunk
	if err = c.Upload(path.Join(folder, newChunkName), chunk); err != nil {
		return err
	}

	// Upload playlist
	if err = c.Upload(path.Join(folder, "index.m3u8"), playlist.Encode()); err != nil {
		return err
	}

	c.VerifyChunk(workOrderID, fmt.Sprintf("%s/%s-%s/%s", c.cfg.BaseStreamURL, userID, appID, chunkname), fmt.Sprintf("https://storage.googleapis.com/%s/%s/%s", c.cfg.Bucket, folder, newChunkName))

	return nil
}

// VerifyChunk blahg
func (c *CSync) VerifyChunk(workOrderID uint32, src string, res string) error {
	client := &http.Client{}

	form := url.Values{}
	form.Add("source_chunk_url", src)
	form.Add("result_chunk_url", res)
	form.Add("job_id", fmt.Sprintf("%d", workOrderID))

	request, err := http.NewRequest("POST", "URL", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	c.log.Infof("verifier response: code [ %d ]", resp.StatusCode)

	return nil
}

// Upload uploads an object to gcs with publicread acl
func (c *CSync) Upload(output string, r io.Reader) error {
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
		CacheControl: "public, max-age=315360000",
	}

	if _, err := svc.Objects.Insert(c.cfg.Bucket, object).Media(r).PredefinedAcl("publicread").Do(); err != nil {

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
