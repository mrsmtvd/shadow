package logging

import (
	"github.com/mrsmtvd/shadow"
)

type Component interface {
	shadow.Component

	Logger() Logger
}
