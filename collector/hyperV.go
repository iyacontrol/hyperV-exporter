package collector

import "github.com/prometheus/client_golang/prometheus"

// HyperVCollector is a Prometheus collector for hyper-v
type HyperVCollector struct {
}

// NewHyperVCollector ...
func NewHyperVCollector() (Collector, error) {
	return &HyperVCollector{}, nil
}

// Collect sends the metric values for each metric
// to the provided prometheus Metric channel.
func (c *HyperVCollector) Collect(ch chan<- prometheus.Metric) error {
	return nil
}
