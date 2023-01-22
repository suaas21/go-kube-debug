package context

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Summary export quantile data
func (p *Prom) Summary(obj map[float64]float64) *Prom {
	if obj == nil {
		obj = map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.95: 0.005,
			0.99: 0.001,
		}
	}

	opts := prometheus.SummaryOpts{
		Namespace:  p.prefix,
		Subsystem:  p.app,
		Objectives: obj,
		Name:       "request_summary_latency",
	}

	s := prometheus.NewSummaryVec(opts, []string{"pattern", "method", "status_code"})
	prometheus.MustRegister(s)

	p.context.Summary = s
	return p
}

// Histogram export histogram data
func (p *Prom) Histogram(bucket []float64) *Prom {
	if bucket == nil {
		bucket = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	}

	opts := prometheus.HistogramOpts{
		Namespace: p.prefix,
		Subsystem: p.app,
		Buckets:   bucket,
		Name:      "request_histogram_latency",
	}

	h := prometheus.NewHistogramVec(opts, []string{"pattern", "method", "status_code"})
	prometheus.MustRegister(h)

	p.context.Histogram = h
	return p
}

func (p *Prom) count() {
	opts := prometheus.CounterOpts{
		Namespace: p.prefix,
		Subsystem: p.app,
		Name:      "request_count",
	}
	c := prometheus.NewCounterVec(opts, []string{"pattern", "method", "status_code"})
	prometheus.MustRegister(c)

	p.context.Count = c
}
