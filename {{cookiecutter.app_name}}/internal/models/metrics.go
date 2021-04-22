package models


import (
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "{{cookiecutter.app_name}}_"

var info uint64

// Metrics implements the prometheus.Metrics interface and
// exposes gorush metrics for prometheus
type Metrics struct {
	TotalInfoCount *prometheus.Desc
}

// NewMetrics returns a new Metrics with all prometheus.Desc initialized
func NewMetrics() Metrics {
	return Metrics{
		TotalInfoCount: prometheus.NewDesc(
			namespace+"total_info_count",
			"Number of info count",
			nil, nil,
		),
	}
}

// Collect returns the metrics with values
func (c Metrics) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.TotalInfoCount,
		prometheus.GaugeValue,
		float64(info),
	)
}

// Describe returns all possible prometheus.Desc
func (c Metrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.TotalInfoCount
}

func InfoCountInc() {
	atomic.AddUint64(&info, 1)
}
