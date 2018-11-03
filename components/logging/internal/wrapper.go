package internal

import (
	"go.uber.org/zap"
)

type wrapper interface {
	Sugar() *zap.SugaredLogger
	Logger() *zap.Logger
	SetLogger(l *zap.Logger)
}
