package capacity

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/videocoin/transcode/transcoder"
)

const testSamplePath = "/opt/capacity_test.mp4"
const mbCount = 1920 * 1080 / (16 * 16) * (10 * 30)

type Capacitor struct {
	transcoder *transcoder.Transcoder

	isInternal   bool
	lastCapacity int
	lastPeformed time.Time
}

func NewCapacitor(isInternal bool, transcoder *transcoder.Transcoder) *Capacitor {
	capacitor := &Capacitor{
		transcoder: transcoder,
		isInternal: isInternal,
	}

	// force sync capacity update on init
	_ = capacitor.getCapacity()

	return capacitor
}

func (c *Capacitor) getCapacity() error {
	c.transcoder.Stop()
	defer c.transcoder.Start()

	args := []string{
		"-i", testSamplePath,
		"-y",
		"-c:a", "copy",
		"-c:v", "libx264",
		"-b:v", "4800k",
		"/tmp/out.mp4",
	}

	start := time.Now()

	cmd := exec.Command("ffmpeg", args...)
	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}

	if !cmd.ProcessState.Success() {
		return fmt.Errorf("ffmpeg process failed")
	}

	c.lastPeformed = time.Now()
	c.lastCapacity = mbCount / int(time.Since(start).Seconds())

	return nil
}

func (c *Capacitor) IsUpdateTime() bool {
	return time.Since(c.lastPeformed.Add(time.Hour)) >= time.Hour
}

func (c *Capacitor) GetInfo() (map[string]interface{}, []byte, error) {
	if !c.isInternal && !c.transcoder.IsWorking() && time.Since(c.lastPeformed.Add(time.Hour)) >= time.Hour {
		go c.getCapacity()
	}

	info := map[string]interface{}{}
	info["encode"] = 0
	if !c.isInternal {
		info["encode"] = c.lastCapacity
	}

	b, err := json.Marshal(info)
	if err != nil {
		return nil, nil, err
	}

	return info, b, nil
}
