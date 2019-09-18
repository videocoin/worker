package telegraf

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/telegraf/agent"
	"github.com/videocoin/telegraf/core"
	"github.com/videocoin/telegraf/core/config"
	"github.com/videocoin/telegraf/core/models"
	_ "github.com/videocoin/telegraf/plugins/aggregators/all"
	"github.com/videocoin/telegraf/plugins/inputs"
	_ "github.com/videocoin/telegraf/plugins/inputs/all"
	"github.com/videocoin/telegraf/plugins/outputs"
	_ "github.com/videocoin/telegraf/plugins/outputs/all"
	_ "github.com/videocoin/telegraf/plugins/processors/all"
	"github.com/videocoin/telegraf/plugins/serializers"
)

var stop chan struct{}

func Run(logger *logrus.Entry, machineID string) error {
	stop = make(chan struct{})
	run(stop, logger, machineID)
	return nil
}

func run(stop chan struct{}, logger *logrus.Entry, machineID string) {
	reload := make(chan bool, 1)
	reload <- true
	for <-reload {
		reload <- false

		ctx, cancel := context.WithCancel(context.Background())

		signals := make(chan os.Signal)
		signal.Notify(
			signals,
			os.Interrupt,
			syscall.SIGHUP,
			syscall.SIGTERM,
			syscall.SIGINT,
		)

		go func() {
			select {
			case sig := <-signals:
				if sig == syscall.SIGHUP {
					log.Info("reloading Telegraf config")
					<-reload
					reload <- true
				}
				cancel()
			case <-stop:
				cancel()
			}
		}()

		err := runAgent(ctx, logger, machineID)
		if err != nil {
			log.Fatalf("error running agent: %v", err)
		}
	}
}

func runAgent(ctx context.Context, logger *logrus.Entry, machineID string) error {
	c := config.NewConfig()
	c.Agent.Hostname, _ = os.Hostname()
	c.Agent.MetricBatchSize = 0
	c.Agent.MetricBufferLimit = 0
	c.Agent.Interval = core.Duration{Duration: time.Second * 5}
	c.Agent.FlushInterval = core.Duration{Duration: time.Second * 5}
	c.Logger = logger

	c.Tags = map[string]string{
		"machine_id": machineID,
	}

	// Inputs

	err := addInput(c, "cpu")
	if err != nil {
		return err
	}

	err = addInput(c, "system")
	if err != nil {
		return err
	}

	err = addInput(c, "mem")
	if err != nil {
		return err
	}

	err = addInput(c, "disk")
	if err != nil {
		return err
	}

	err = addOutput(c, "graphite")
	if err != nil {
		return err
	}

	ag, err := agent.NewAgent(c, logger)
	if err != nil {
		return err
	}

	return ag.Run(ctx)
}

func addInput(c *config.Config, inputFilter string) error {
	inputCreator, ok := inputs.Inputs[inputFilter]
	if !ok {
		return fmt.Errorf("Undefined but requested input: %s", inputFilter)
	}

	input := inputCreator()
	inputConfig := &models.InputConfig{
		Name:              inputFilter,
		Interval:          c.Agent.Interval.Duration,
		MeasurementPrefix: fmt.Sprintf("miners.agents."),
	}

	rp := models.NewRunningInput(input, inputConfig)
	rp.SetDefaultTags(c.Tags)

	// Add input
	c.Inputs = append(c.Inputs, rp)

	return nil
}

func addOutput(c *config.Config, outputFilter string) error {
	outputCreator, ok := outputs.Outputs[outputFilter]
	if !ok {
		return fmt.Errorf("Undefined but requested output: %s", outputFilter)
	}
	output := outputCreator()

	switch t := output.(type) {
	case serializers.SerializerOutput:
		serializerConfig := &serializers.Config{
			TimestampUnits:     time.Duration(1 * time.Second),
			DataFormat:         "graphite",
			GraphiteTagSupport: true,
			//Prefix:             c.Agent.Hostname,
		}
		serializer, err := serializers.NewSerializer(serializerConfig)
		if err != nil {
			return err
		}
		t.SetSerializer(serializer)
	}

	fltr := models.Filter{}
	outputConfig := &models.OutputConfig{
		Name:   outputFilter,
		Filter: fltr,
	}

	ro := models.NewRunningOutput(
		outputFilter,
		output,
		outputConfig,
		c.Agent.MetricBatchSize,
		c.Agent.MetricBufferLimit,
		c.Logger,
	)

	c.Outputs = append(c.Outputs, ro)

	return nil
}
