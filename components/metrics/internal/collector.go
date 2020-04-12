package internal

import (
	"github.com/kihamo/snitch"
)

type OnOffCollector struct {
	original snitch.Collector
	check    func() bool
}

func NewOnOffCollector(original snitch.Collector, check func() bool) *OnOffCollector {
	return &OnOffCollector{
		original: original,
		check:    check,
	}
}

func (c *OnOffCollector) Describe(ch chan<- *snitch.Description) {
	c.original.Describe(ch)
}

func (c *OnOffCollector) Collect(ch chan<- snitch.Metric) {
	if c.check() {
		c.original.Collect(ch)
	}
}
