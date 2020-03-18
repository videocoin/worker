package capacity

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

const testSamplePath = "/opt/test.mp4"
const mbCount = 1920 * 1080 / (16 * 16) * (10 * 30)

type Capacitor struct {
	isInternal   bool
	lastCapacity int
	lastPeformed time.Time
}

func NewCapacitor(isInternal bool) *Capacitor {
	capacitor := &Capacitor{
		isInternal: isInternal,
	}

	// force sync capacity update on init
	_, _ = capacitor.getCapacity()

	return capacitor
}

func (c *Capacitor) getCapacity() (int, error) {
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
		return 0, err
	}

	if !cmd.ProcessState.Success() {
		return 0, fmt.Errorf("ffmpeg process failed")
	}

	c.lastPeformed = time.Now()
	c.lastCapacity = mbCount / int(time.Since(start).Seconds())

	return c.lastCapacity, nil
}

func (c *Capacitor) GetInfo() (map[string]interface{}, []byte, error) {
	info := map[string]interface{}{}
	info["encode"] = 0

	if !c.isInternal && time.Since(c.lastPeformed.Add(time.Hour)) >= time.Hour {
		go c.getCapacity()
	}

	if !c.isInternal {
		info["encode"] = c.lastCapacity
	}

	b, err := json.Marshal(info)
	if err != nil {
		return nil, nil, err
	}

	return info, b, nil
}
