package pinger

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	v1 "github.com/videocoin/cloud-api/dispatcher/v1"
	minersv1 "github.com/videocoin/cloud-api/miners/v1"
	"github.com/videocoin/transcode/capacity"
	"github.com/videocoin/transcode/sysinfo"
)

type Pinger struct {
	logger     *logrus.Entry
	dispatcher v1.DispatcherServiceClient
	capacitor  *capacity.Capacitor
	clientID   string
	timeout    time.Duration
	ticker     *time.Ticker
	appVersion string
}

func NewPinger(
	dispatcher v1.DispatcherServiceClient,
	capacitor *capacity.Capacitor,
	clientID string,
	timeout time.Duration,
	appVersion string,
	logger *logrus.Entry,
) (*Pinger, error) {
	ticker := time.NewTicker(timeout)
	return &Pinger{
		logger:     logger,
		capacitor:  capacitor,
		dispatcher: dispatcher,
		clientID:   clientID,
		timeout:    timeout,
		ticker:     ticker,
		appVersion: appVersion,
	}, nil
}

func (p *Pinger) Start() {
	p.logger.Debugf("starting pinger")

	for range p.ticker.C {
		ctx := context.Background()
		si := &sysinfo.SystemInfo{AppVersion: p.appVersion, Logger: p.logger}
		_, systemInfo, _ := si.GetInfo()
		_, capacityInfo, _ := p.capacitor.GetInfo()

		req := &minersv1.PingRequest{
			ClientID:     p.clientID,
			SystemInfo:   systemInfo,
			CapacityInfo: capacityInfo,
		}
		_, err := p.dispatcher.Ping(ctx, req)
		if err != nil {
			p.logger.Errorf("failed to ping: %s", err)
			continue
		}

		p.logger.Debugf("ping")
	}
}

func (p *Pinger) Stop() error {
	p.logger.Debugf("stopping pinger")
	p.ticker.Stop()
	return nil
}
