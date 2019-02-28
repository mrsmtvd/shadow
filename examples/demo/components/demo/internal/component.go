package internal

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/examples/demo/components/demo"
)

type Component struct {
}

func (c *Component) Name() string {
	return demo.ComponentName
}

func (c *Component) Version() string {
	return demo.ComponentVersion
}

func (c *Component) Run(shadow.Application, chan<- struct{}) error {
	return nil
}
