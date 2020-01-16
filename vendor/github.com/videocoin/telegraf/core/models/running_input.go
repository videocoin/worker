package models

import (
	"time"

	"github.com/videocoin/telegraf"
)

type RunningInput struct {
	Input  telegraf.Input
	Config *InputConfig

	defaultTags map[string]string
}

func NewRunningInput(input telegraf.Input, config *InputConfig) *RunningInput {
	return &RunningInput{
		Input:  input,
		Config: config,
	}
}

// InputConfig is the common config for all inputs.
type InputConfig struct {
	Name     string
	Interval time.Duration

	NameOverride      string
	MeasurementPrefix string
	MeasurementSuffix string
	Tags              map[string]string
	Filter            Filter
}

func (r *RunningInput) Name() string {
	return "inputs." + r.Config.Name
}

func (r *RunningInput) metricFiltered(metric telegraf.Metric) {
	metric.Drop()
}

func (r *RunningInput) MakeMetric(metric telegraf.Metric) telegraf.Metric {
	if ok := r.Config.Filter.Select(metric); !ok {
		r.metricFiltered(metric)
		return nil
	}

	m := makemetric(
		metric,
		r.Config.NameOverride,
		r.Config.MeasurementPrefix,
		r.Config.MeasurementSuffix,
		r.Config.Tags,
		r.defaultTags)

	r.Config.Filter.Modify(metric)
	if len(metric.FieldList()) == 0 {
		r.metricFiltered(metric)
		return nil
	}

	return m
}

func (r *RunningInput) Gather(acc telegraf.Accumulator) error {
	err := r.Input.Gather(acc)
	return err
}

func (r *RunningInput) SetDefaultTags(tags map[string]string) {
	r.defaultTags = tags
}
