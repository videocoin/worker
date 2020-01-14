package hlswatcher

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
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
	Source    string  `json:"source"`
	Num       uint64  `json:"num"`
	Angle     uint32  `json:"angle"`
	Name      string  `json:"name"`
	Duration  float64 `json:"duration"`
	IsLast    bool    `json:"is_last"`
	VariantID string  `json:"variant_id"`
}

type Watcher struct {
	SegmentsCh chan *SegmentInfo
	ErrCh      chan error

	wg          *sync.WaitGroup
	mu          *sync.Mutex
	running     bool
	files       []string
	segmentsNum map[string]uint64
	period      time.Duration
}

func New(period time.Duration) *Watcher {
	return &Watcher{
		SegmentsCh:  make(chan *SegmentInfo, 1),
		ErrCh:       make(chan error, 1),
		wg:          &sync.WaitGroup{},
		mu:          new(sync.Mutex),
		files:       make([]string, 0),
		segmentsNum: make(map[string]uint64),
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
	segmentsNum := make(map[string]uint64)

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
			segmentsNum[path] = uint64(len(segmentList))
		}

		for path, ss := range segments {
			for _, s := range ss {
				if s.Num <= w.segmentsNum[path] {
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
	w.segmentsNum[name] = 0

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
	segments := []*SegmentInfo{}

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

			pathParts := strings.Split(name, "/")
			plFileName := pathParts[len(pathParts)-1]
			variantID := plFileName[0 : len(plFileName)-6]

			ts := &SegmentInfo{
				Num:       uint64(idx) + 1,
				Angle:     extractAngle(segment.URI),
				Name:      segment.URI,
				Duration:  segment.Duration,
				Source:    path.Join(filepath.Dir(name), segment.URI),
				VariantID: variantID,
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

func extractAngle(name string) uint32 {
	angle := uint64(0)
	p1 := strings.Split(name, ".")
	if len(p1) == 2 {
		p2 := strings.Split(p1[0], "-")
		if len(p2) >= 5 {
			angleStr := p2[len(p2)-2]
			angle, _ = strconv.ParseUint(angleStr, 10, 32)
		}
	}

	return uint32(angle)
}
