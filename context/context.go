package context

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Prom struct {
	context *Context
	prefix  string
	app     string
}

func NewProm(app string) *Prom {
	prom := &Prom{
		context: &Context{},
		prefix:  "hubble",
		app:     app,
	}
	prom.count()
	return prom
}

type Context struct {
	Summary   *prometheus.SummaryVec
	Histogram *prometheus.HistogramVec
	Count     *prometheus.CounterVec
}

func (p *Prom) GetContext() *Context {
	return p.context
}

func (c *Context) Push(pattern, method string, code int, then time.Time) {

	values := []string{pattern, method, fmt.Sprintf("%v", code)}

	c.Count.WithLabelValues(values...).Add(1)

	dur := float64(time.Since(then)/time.Millisecond) / 1000.0
	if c.Summary != nil {
		c.Summary.WithLabelValues(values...).Observe(dur)
	}
	if c.Histogram != nil {
		c.Histogram.WithLabelValues(values...).Observe(dur)
	}
}
