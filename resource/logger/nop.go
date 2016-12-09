package logger

import (
	"os"

	"github.com/rs/xlog"
)

type nop struct{}

var NopLogger = &nop{}

func (n nop) Write(p []byte) (_ int, _ error) { return len(p), nil }

func (n nop) SetField(_ string, _ interface{}) {}

func (n nop) Debug(_ ...interface{}) {}

func (n nop) Debugf(_ string, _ ...interface{}) {}

func (n nop) Info(_ ...interface{}) {}

func (n nop) Infof(_ string, _ ...interface{}) {}

func (n nop) Warn(_ ...interface{}) {}

func (n nop) Warnf(_ string, _ ...interface{}) {}

func (n nop) Error(_ ...interface{}) {}

func (n nop) Errorf(_ string, _ ...interface{}) {}

func (n nop) Fatal(_ ...interface{}) { os.Exit(0) }

func (n nop) Fatalf(_ string, _ ...interface{}) { os.Exit(0) }

func (n nop) Output(_ int, _ string) error { return nil }

func (n nop) OutputF(_ xlog.Level, _ int, _ string, _ map[string]interface{}) {}

func (n nop) Printf(_ string, _ ...interface{}) {}

func (n nop) Print(_ ...interface{}) {}

func (n nop) Println(_ ...interface{}) {}

func (n nop) Panic(_ ...interface{}) { os.Exit(0) }

func (n nop) Panicf(_ string, _ ...interface{}) { os.Exit(0) }

func (n nop) Panicln(_ ...interface{}) { os.Exit(0) }

func (n nop) Log(_ ...interface{}) {}
