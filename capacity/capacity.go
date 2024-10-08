package capacity

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/worker/transcoder"
)

const testSamplePath = "/opt/capacity_test.mp4"
const mbCount = 1920 * 1080 / (16 * 16) * (10 * 30)

type Capacitor struct {
	transcoder *transcoder.Transcoder
	logger     *logrus.Entry

	isInternal         bool
	lastEncodeCapacity int
	lastCPUCapacity    int
	lastPeformed       time.Time
}

func NewCapacitor(isInternal bool, transcoder *transcoder.Transcoder, logger *logrus.Entry) *Capacitor {
	capacitor := &Capacitor{
		transcoder: transcoder,
		isInternal: isInternal,
		logger:     logger,
	}

	// force sync capacity update on init
	if !isInternal {
		if err := capacitor.getCPUCapacity(); err != nil {
			logger.WithError(err).Errorf("failed to get cpu capacity")
		}
		if err := capacitor.getEncodeCapacity(); err != nil {
			logger.WithError(err).Errorf("failed to get encode capacity")
		}
	}

	return capacitor
}

func (c *Capacitor) getEncodeCapacity() error {
	if c.transcoder.IsRunning() {
		err := c.transcoder.Stop()
		if err != nil {
			return err
		}
		defer func() {
			go func() {
				err := c.transcoder.Start()
				if err != nil {
					c.logger.Error(err)
				}

			}()
		}()
	}

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
	c.lastEncodeCapacity = mbCount / int(time.Since(start).Seconds())

	c.logger.Infof("encode capacity is %d", c.lastEncodeCapacity)

	return nil
}

func (c *Capacitor) getCPUCapacity() error {
	p, err := cpu.Percent(5*time.Second, false)
	if err != nil {
		return err
	}

	if len(p) < 1 {
		return fmt.Errorf("failed to get cpu usage: no results")
	}

	c.lastCPUCapacity = 100 - int(p[0])

	c.logger.Infof("cpu capacity is %d", c.lastCPUCapacity)

	return nil
}

func (c *Capacitor) IsUpdateTime() bool {
	return time.Since(c.lastPeformed.Add(time.Hour)) >= time.Hour
}

func (c *Capacitor) GetInfo() (map[string]interface{}, []byte, error) {
	defer func() {
		// if !c.isInternal && !c.transcoder.IsWorking() && time.Since(c.lastPeformed) >= time.Minute*15 {
		// 	c.logger.Info("get encode capacity")
		// 	if err := c.getEncodeCapacity(); err != nil {
		// 		c.logger.WithError(err).Errorf("failed to get encode capacity")
		// 	}
		// }

		c.logger.Info("get cpu capacity")
		if err := c.getCPUCapacity(); err != nil {
			c.logger.WithError(err).Errorf("failed to get cpu capacity")
		}
	}()

	info := map[string]interface{}{}
	info["encode"] = 0
	info["cpu"] = 0

	if !c.isInternal {
		info["encode"] = c.lastEncodeCapacity
		info["cpu"] = c.lastCPUCapacity
	}

	b, err := json.Marshal(info)
	if err != nil {
		return nil, nil, err
	}

	return info, b, nil
}
