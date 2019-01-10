package transcode

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"

	storage "google.golang.org/api/storage/v1"

	"github.com/grafov/m3u8"
)

// CSync syncer struct
type CSync struct {
	ctx context.Context
	cfg *Config
	log *logrus.Entry
}

func CSyncInit() *CSync {

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

func syncDirectory(bucket string, folder string, streamHash string, inputURL string) {
	//create playlist
	// wait for chunk
	// get chunk dir
	// append to playlist
	// upload chunk
	// upload playlist

	playlist, err := m3u8.NewMediaPlaylist(0, 0)
	if err != nil {
		_ = err
	}
	_ = playlist

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

				if path.Base(event.Name) == "0.ts" {
					duration, err := getDuration(path.Join(folder, event.Name))
					if err != nil {
						log.Errorf("failed to get duration: %s", err.Error())
					}

					playlist.TargetDuration = duration
				}

				if (event.Op&fsnotify.Create == fsnotify.Create) && !strings.Contains(event.Name, "tmp") && !strings.Contains(event.Name, ".m3u8") {
					log.Println("created file:", path.Base(event.Name))
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(folder)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func (c *CSync) doWork(chunkname string, folder string, playlist *m3u8.MediaPlaylist) error {
	// create playlist
	// wait for chunk
	// get chunk dir
	// append to playlist
	// upload chunk
	// upload playlist
	// call verifier

	var b = make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return err
	}

	duration, err := getDuration(path.Join(folder, chunkname))
	if err != nil {
		return err
	}

	newName := fmt.Sprintf("%x.ts", b)

	if err = playlist.Append(newName, duration, ""); err != nil {
		return err
	}

	return nil
}

func (c *CSync) upload(filename string, output string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

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

	if _, err := svc.Objects.Insert(c.cfg.Bucket, object).Media(f).PredefinedAcl("publicread").Do(); err != nil {
		return err
	}

	return nil
}
