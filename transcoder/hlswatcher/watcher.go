package hlswatcher

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/grafov/m3u8"
)

var (
	ErrDurationTooShort = errors.New("hlswatcher: duration is less than 1ns")
	ErrWatcherRunning   = errors.New("hlswatcher: watcher is already running")
)

type SegmentInfo struct {
	Index    uint64  `json:"index"`
	Source   string  `json:"source"`
	Name     string  `json:"name"`
	Num      uint64  `json:"num"`
	Duration float64 `json:"duration"`
	IsLast   bool    `json:"is_last"`
	IsVOD    bool    `json:"is_vod"`
}

type Watcher struct {
	SegmentsCh chan *SegmentInfo
	ErrCh      chan error

	wg          *sync.WaitGroup
	mu          *sync.Mutex
	running     bool
	files       []string
	segmentsNum map[string]int64
	period      time.Duration
}

func New(period time.Duration) *Watcher {
	return &Watcher{
		SegmentsCh:  make(chan *SegmentInfo, 1),
		ErrCh:       make(chan error, 1),
		wg:          &sync.WaitGroup{},
		mu:          new(sync.Mutex),
		files:       make([]string, 0),
		segmentsNum: make(map[string]int64),
		period:      period,
	}
}

func (w *Watcher) Start() error {
	if w.period < time.Nanosecond {
		return ErrDurationTooShort
	}

	w.wg.Add(1)
	w.mu.Lock()

	if w.running {
		w.mu.Unlock()
		return ErrWatcherRunning
	}
	w.running = true
	w.SegmentsCh = make(chan *SegmentInfo, 1)
	w.ErrCh = make(chan error, 1)

	segments := make(map[string][]*SegmentInfo)
	segmentsNum := make(map[string]int64)

	w.mu.Unlock()
	w.wg.Done()

	for {
		if !w.running {
			return nil
		}

		files := w.Files()

		for _, path := range files {
			_, err := os.Stat(path)
			if err != nil {
				continue
			}

			segmentList, err := w.extractSegments(path)
			if err != nil {
				go func() {
					if !w.running {
						return
					}
					w.ErrCh <- err
				}()
				continue
			}

			segments[path] = segmentList
			segmentsNum[path] = int64(segmentList[len(segmentList)-1].Index)
		}

		for path, ss := range segments {
			for _, s := range ss {
				if int64(s.Index) <= w.segmentsNum[path] {
					continue
				}

				if !w.running {
					return nil
				}

				w.SegmentsCh <- s
			}

			w.segmentsNum[path] = segmentsNum[path]
		}

		time.Sleep(w.period)
	}
}

func (w *Watcher) Wait() {
	w.wg.Wait()
}

func (w *Watcher) Add(name string) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if filepath.Ext(strings.TrimSpace(name)) != ".m3u8" {
		return fmt.Errorf("incorrect filename type")
	}

	w.files = append(w.files, name)
	w.segmentsNum[name] = -1

	return nil
}

func (w *Watcher) Remove(name string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	removeIndex := -1
	for i := range w.files {
		if w.files[i] == name {
			removeIndex = i
			break
		}
	}

	if removeIndex == -1 {
		return nil
	}

	w.files = append(w.files[:removeIndex], w.files[removeIndex+1:]...)

	return nil
}

func (w *Watcher) Files() []string {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.files
}

func (w *Watcher) Stop() {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return
	}
	defer w.mu.Unlock()

	w.running = false
	w.files = make([]string, 0)
}

func (w *Watcher) ExtractSegments(name string) ([]*SegmentInfo, error) {
	return w.extractSegments(name)
}

func (w *Watcher) extractSegments(name string) ([]*SegmentInfo, error) {
	var segments []*SegmentInfo

	content, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to open playlist: %s", err)
	}
	buf := bytes.NewBuffer(content)

	p, listType, err := m3u8.Decode(*buf, false)
	if err != nil {
		return nil, fmt.Errorf("failed to decode playlist: %s", err)
	}

	if listType == m3u8.MEDIA {
		pl := p.(*m3u8.MediaPlaylist)

		for idx, segment := range pl.Segments {
			if segment == nil {
				continue
			}

			ts := &SegmentInfo{
				Index:    uint64(idx),
				Num:      uint64(idx) + 1,
				Name:     segment.URI,
				Duration: segment.Duration,
				Source:   path.Join(filepath.Dir(name), segment.URI),
			}

			segments = append(segments, ts)
		}

		if pl.Closed {
			segments[len(segments)-1].IsLast = true
		}
	} else {
		return nil, fmt.Errorf("hlswatcher: wrong playlist type")
	}

	return segments, nil
}
