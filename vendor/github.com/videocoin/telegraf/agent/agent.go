package agent

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/videocoin/telegraf"
	"github.com/videocoin/telegraf/core"
	"github.com/videocoin/telegraf/core/config"
	"github.com/videocoin/telegraf/core/models"
	"github.com/videocoin/telegraf/plugins/serializers/influx"
	"github.com/sirupsen/logrus"
)

// Agent runs a set of plugins.
type Agent struct {
	Config *config.Config
	Logger *logrus.Entry
}

// NewAgent returns an Agent for the given Config.
func NewAgent(config *config.Config, logger *logrus.Entry) (*Agent, error) {
	a := &Agent{
		Config: config,
		Logger: logger,
	}
	return a, nil
}

// Run starts and runs the Agent until the context is done.
func (a *Agent) Run(ctx context.Context) error {
	a.Logger.Infof("config: interval:%s, quiet:%#v, hostname:%#v, "+
		"flush interval:%s",
		a.Config.Agent.Interval.Duration, a.Config.Agent.Quiet,
		a.Config.Agent.Hostname, a.Config.Agent.FlushInterval.Duration)

	if ctx.Err() != nil {
		return ctx.Err()
	}

	a.Logger.Debug("connecting outputs")
	err := a.connectOutputs(ctx)
	if err != nil {
		return err
	}

	inputC := make(chan telegraf.Metric, 100)
	procC := make(chan telegraf.Metric, 100)
	outputC := make(chan telegraf.Metric, 100)

	startTime := time.Now()

	a.Logger.Debug("starting service inputs")
	err = a.startServiceInputs(ctx, inputC)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	src := inputC
	dst := inputC

	wg.Add(1)
	go func(dst chan telegraf.Metric) {
		defer wg.Done()

		err := a.runInputs(ctx, startTime, dst)
		if err != nil {
			a.Logger.Errorf("error running inputs: %v", err)
		}

		a.Logger.Debugf("stopping service inputs")
		a.stopServiceInputs()

		close(dst)
		a.Logger.Debugf("input channel closed")
	}(dst)

	src = dst

	if len(a.Config.Processors) > 0 {
		dst = procC

		wg.Add(1)
		go func(src, dst chan telegraf.Metric) {
			defer wg.Done()

			err := a.runProcessors(src, dst)
			if err != nil {
				a.Logger.Errorf("rrror running processors: %v", err)
			}
			close(dst)
			a.Logger.Debug("processor channel closed")
		}(src, dst)

		src = dst
	}

	if len(a.Config.Aggregators) > 0 {
		dst = outputC

		wg.Add(1)
		go func(src, dst chan telegraf.Metric) {
			defer wg.Done()

			err := a.runAggregators(startTime, src, dst)
			if err != nil {
				a.Logger.Errorf("error running aggregators: %v", err)
			}
			close(dst)
			a.Logger.Debug("output channel closed")
		}(src, dst)

		src = dst
	}

	wg.Add(1)
	go func(src chan telegraf.Metric) {
		defer wg.Done()

		err := a.runOutputs(startTime, src)
		if err != nil {
			a.Logger.Errorf("error running outputs: %v", err)
		}
	}(src)

	wg.Wait()

	a.Logger.Debugf("closing outputs")
	a.closeOutputs()

	a.Logger.Debug("stopped successfully")
	return nil
}

// Test runs the inputs once and prints the output to stdout in line protocol.
func (a *Agent) Test(ctx context.Context) error {
	var wg sync.WaitGroup
	metricC := make(chan telegraf.Metric)
	nulC := make(chan telegraf.Metric)
	defer func() {
		close(metricC)
		close(nulC)
		wg.Wait()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		s := influx.NewSerializer()
		s.SetFieldSortOrder(influx.SortFields)
		for metric := range metricC {
			octets, err := s.Serialize(metric)
			if err == nil {
				fmt.Print("> ", string(octets))

			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range nulC {
		}
	}()

	for _, input := range a.Config.Inputs {
		select {
		case <-ctx.Done():
			return nil
		default:
			if _, ok := input.Input.(telegraf.ServiceInput); ok {
				a.Logger.Warningf("skipping plugin [[%s]]: service inputs not supported in --test mode",
					input.Name())
				continue
			}

			acc := NewAccumulator(input, metricC, a.Logger)
			acc.SetPrecision(a.Precision())
			input.SetDefaultTags(a.Config.Tags)

			// Special instructions for some inputs. cpu, for example, needs to be
			// run twice in order to return cpu usage percentages.
			switch input.Name() {
			case "inputs.cpu", "inputs.mongodb", "inputs.procstat":
				nulAcc := NewAccumulator(input, nulC, a.Logger)
				nulAcc.SetPrecision(a.Precision())
				if err := input.Input.Gather(nulAcc); err != nil {
					return err
				}

				time.Sleep(500 * time.Millisecond)
				if err := input.Input.Gather(acc); err != nil {
					return err
				}
			default:
				if err := input.Input.Gather(acc); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// runInputs starts and triggers the periodic gather for Inputs.
//
// When the context is done the timers are stopped and this function returns
// after all ongoing Gather calls complete.
func (a *Agent) runInputs(
	ctx context.Context,
	startTime time.Time,
	dst chan<- telegraf.Metric,
) error {
	var wg sync.WaitGroup
	for _, input := range a.Config.Inputs {
		interval := a.Config.Agent.Interval.Duration
		jitter := a.Config.Agent.CollectionJitter.Duration

		// Overwrite agent interval if this plugin has its own.
		if input.Config.Interval != 0 {
			interval = input.Config.Interval
		}

		acc := NewAccumulator(input, dst, a.Logger)
		acc.SetPrecision(a.Precision())

		wg.Add(1)
		go func(input *models.RunningInput) {
			defer wg.Done()

			if a.Config.Agent.RoundInterval {
				err := core.SleepContext(
					ctx, core.AlignDuration(startTime, interval))
				if err != nil {
					return
				}
			}

			a.gatherOnInterval(ctx, acc, input, interval, jitter)
		}(input)
	}
	wg.Wait()

	return nil
}

// gather runs an input's gather function periodically until the context is
// done.
func (a *Agent) gatherOnInterval(
	ctx context.Context,
	acc telegraf.Accumulator,
	input *models.RunningInput,
	interval time.Duration,
	jitter time.Duration,
) {
	defer panicRecover(input)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		err := core.SleepContext(ctx, core.RandomDuration(jitter))
		if err != nil {
			return
		}

		err = a.gatherOnce(acc, input, interval)
		if err != nil {
			acc.AddError(err)
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}

// gatherOnce runs the input's Gather function once, logging a warning each
// interval it fails to complete before.
func (a *Agent) gatherOnce(
	acc telegraf.Accumulator,
	input *models.RunningInput,
	timeout time.Duration,
) error {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	done := make(chan error)
	go func() {
		done <- input.Gather(acc)
	}()

	for {
		select {
		case err := <-done:
			return err
		case <-ticker.C:
			a.Logger.Warningf("input %q did not complete within its interval",
				input.Name())
		}
	}
}

// runProcessors applies processors to metrics.
func (a *Agent) runProcessors(
	src <-chan telegraf.Metric,
	agg chan<- telegraf.Metric,
) error {
	for metric := range src {
		metrics := a.applyProcessors(metric)

		for _, metric := range metrics {
			agg <- metric
		}
	}

	return nil
}

// applyProcessors applies all processors to a metric.
func (a *Agent) applyProcessors(m telegraf.Metric) []telegraf.Metric {
	metrics := []telegraf.Metric{m}
	for _, processor := range a.Config.Processors {
		metrics = processor.Apply(metrics...)
	}

	return metrics
}

func updateWindow(start time.Time, roundInterval bool, period time.Duration) (time.Time, time.Time) {
	var until time.Time
	if roundInterval {
		until = core.AlignTime(start, period)
		if until == start {
			until = core.AlignTime(start.Add(time.Nanosecond), period)
		}
	} else {
		until = start.Add(period)
	}

	since := until.Add(-period)

	return since, until
}

// runAggregators adds metrics to the aggregators and triggers their periodic
// push call.
//
// Runs until src is closed and all metrics have been processed.  Will call
// push one final time before returning.
func (a *Agent) runAggregators(
	startTime time.Time,
	src <-chan telegraf.Metric,
	dst chan<- telegraf.Metric,
) error {
	ctx, cancel := context.WithCancel(context.Background())

	// Before calling Add, initialize the aggregation window.  This ensures
	// that any metric created after start time will be aggregated.
	for _, agg := range a.Config.Aggregators {
		since, until := updateWindow(startTime, a.Config.Agent.RoundInterval, agg.Period())
		agg.UpdateWindow(since, until)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for metric := range src {
			var dropOriginal bool
			for _, agg := range a.Config.Aggregators {
				if ok := agg.Add(metric); ok {
					dropOriginal = true
				}
			}

			if !dropOriginal {
				dst <- metric
			} else {
				metric.Drop()
			}
		}
		cancel()
	}()

	aggregations := make(chan telegraf.Metric, 100)
	wg.Add(1)
	go func() {
		defer wg.Done()

		var aggWg sync.WaitGroup
		for _, agg := range a.Config.Aggregators {
			aggWg.Add(1)
			go func(agg *models.RunningAggregator) {
				defer aggWg.Done()

				acc := NewAccumulator(agg, aggregations, a.Logger)
				acc.SetPrecision(a.Precision())
				a.push(ctx, agg, acc)
			}(agg)
		}

		aggWg.Wait()
		close(aggregations)
	}()

	for metric := range aggregations {
		metrics := a.applyProcessors(metric)
		for _, metric := range metrics {
			dst <- metric
		}
	}

	wg.Wait()
	return nil
}

// push runs the push for a single aggregator every period.
func (a *Agent) push(
	ctx context.Context,
	aggregator *models.RunningAggregator,
	acc telegraf.Accumulator,
) {
	for {
		// Ensures that Push will be called for each period, even if it has
		// already elapsed before this function is called.  This is guaranteed
		// because so long as only Push updates the EndPeriod.  This method
		// also avoids drift by not using a ticker.
		until := time.Until(aggregator.EndPeriod())

		select {
		case <-time.After(until):
			aggregator.Push(acc)
			break
		case <-ctx.Done():
			aggregator.Push(acc)
			return
		}
	}
}

// runOutputs triggers the periodic write for Outputs.
//

// Runs until src is closed and all metrics have been processed.  Will call
// Write one final time before returning.
func (a *Agent) runOutputs(
	startTime time.Time,
	src <-chan telegraf.Metric,
) error {
	interval := a.Config.Agent.FlushInterval.Duration
	jitter := a.Config.Agent.FlushJitter.Duration

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	for _, output := range a.Config.Outputs {
		interval := interval
		// Overwrite agent flush_interval if this plugin has its own.
		if output.Config.FlushInterval != 0 {
			interval = output.Config.FlushInterval
		}

		wg.Add(1)
		go func(output *models.RunningOutput) {
			defer wg.Done()

			if a.Config.Agent.RoundInterval {
				err := core.SleepContext(
					ctx, core.AlignDuration(startTime, interval))
				if err != nil {
					return
				}
			}

			a.flush(ctx, output, interval, jitter)
		}(output)
	}

	for metric := range src {
		for i, output := range a.Config.Outputs {
			if i == len(a.Config.Outputs)-1 {
				output.AddMetric(metric)
			} else {
				output.AddMetric(metric.Copy())
			}
		}
	}

	a.Logger.Info("hang on, flushing any cached metrics before shutdown")
	cancel()
	wg.Wait()

	return nil
}

// flush runs an output's flush function periodically until the context is
// done.
func (a *Agent) flush(
	ctx context.Context,
	output *models.RunningOutput,
	interval time.Duration,
	jitter time.Duration,
) {
	// since we are watching two channels we need a ticker with the jitter
	// integrated.
	ticker := NewTicker(interval, jitter)
	defer ticker.Stop()

	logError := func(err error) {
		if err != nil {
			a.Logger.Errorf("error writing to output [%s]: %v", output.Name, err)
		}
	}

	for {
		// Favor shutdown over other methods.
		select {
		case <-ctx.Done():
			logError(a.flushOnce(output, interval, output.Write))
			return
		default:
		}

		select {
		case <-ticker.C:
			logError(a.flushOnce(output, interval, output.Write))
		case <-output.BatchReady:
			// Favor the ticker over batch ready
			select {
			case <-ticker.C:
				logError(a.flushOnce(output, interval, output.Write))
			default:
				logError(a.flushOnce(output, interval, output.WriteBatch))
			}
		case <-ctx.Done():
			logError(a.flushOnce(output, interval, output.Write))
			return
		}
	}
}

// flushOnce runs the output's Write function once, logging a warning each
// interval it fails to complete before.
func (a *Agent) flushOnce(
	output *models.RunningOutput,
	timeout time.Duration,
	writeFunc func() error,
) error {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	done := make(chan error)
	go func() {
		done <- writeFunc()
	}()

	for {
		select {
		case err := <-done:
			output.LogBufferStatus()
			return err
		case <-ticker.C:
			a.Logger.Warningf("output %q did not complete within its flush interval", output.Name)
			output.LogBufferStatus()
		}
	}

}

// connectOutputs connects to all outputs.
func (a *Agent) connectOutputs(ctx context.Context) error {
	for _, output := range a.Config.Outputs {
		a.Logger.Debugf("attempting connection to output: %s", output.Name)
		err := output.Output.Connect()
		if err != nil {
			a.Logger.Errorf("failed to connect to output %s, retrying in 15s, "+
				"error was '%s' \n", output.Name, err)

			err := core.SleepContext(ctx, 15*time.Second)
			if err != nil {
				return err
			}

			err = output.Output.Connect()
			if err != nil {
				return err
			}
		}
		a.Logger.Debugf("successfully connected to output: %s", output.Name)
	}
	return nil
}

// closeOutputs closes all outputs.
func (a *Agent) closeOutputs() {
	for _, output := range a.Config.Outputs {
		output.Close()
	}
}

// startServiceInputs starts all service inputs.
func (a *Agent) startServiceInputs(
	ctx context.Context,
	dst chan<- telegraf.Metric,
) error {
	started := []telegraf.ServiceInput{}

	for _, input := range a.Config.Inputs {
		if si, ok := input.Input.(telegraf.ServiceInput); ok {
			// Service input plugins are not subject to timestamp rounding.
			// This only applies to the accumulator passed to Start(), the
			// Gather() accumulator does apply rounding according to the
			// precision agent setting.
			acc := NewAccumulator(input, dst, a.Logger)
			acc.SetPrecision(time.Nanosecond)

			err := si.Start(acc)
			if err != nil {
				a.Logger.Errorf("service for input %s failed to start: %v",
					input.Name(), err)

				for _, si := range started {
					si.Stop()
				}

				return err
			}

			started = append(started, si)
		}
	}

	return nil
}

// stopServiceInputs stops all service inputs.
func (a *Agent) stopServiceInputs() {
	for _, input := range a.Config.Inputs {
		if si, ok := input.Input.(telegraf.ServiceInput); ok {
			si.Stop()
		}
	}
}

// Returns the rounding precision for metrics.
func (a *Agent) Precision() time.Duration {
	precision := a.Config.Agent.Precision.Duration
	interval := a.Config.Agent.Interval.Duration

	if precision > 0 {
		return precision
	}

	switch {
	case interval >= time.Second:
		return time.Second
	case interval >= time.Millisecond:
		return time.Millisecond
	case interval >= time.Microsecond:
		return time.Microsecond
	default:
		return time.Nanosecond
	}
}

// panicRecover displays an error if an input panics.
func panicRecover(input *models.RunningInput) {
	if err := recover(); err != nil {
		trace := make([]byte, 2048)
		runtime.Stack(trace, true)
	}
}
