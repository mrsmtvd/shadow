package wrapper

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Wrapper struct {
	lock sync.RWMutex

	encoder      zapcore.Encoder
	writeSyncer  zapcore.WriteSyncer
	levelEnabler LevelEnabler
	options      []zap.Option

	name   string
	logger *zap.Logger
	sugar  *zap.SugaredLogger

	tree map[string]*Wrapper
}

func NewNop(name string) *Wrapper {
	return &Wrapper{
		options: make([]zap.Option, 0),
		name:    name,
		logger:  zap.L(),
		sugar:   zap.S(),
		tree:    make(map[string]*Wrapper),
	}
}

func New(name string, enc zapcore.Encoder, ws zapcore.WriteSyncer, level LevelEnabler, options ...zap.Option) *Wrapper {
	w := &Wrapper{
		name: name,
		tree: make(map[string]*Wrapper),
	}

	w.InitFull(enc, ws, level, options...)
	return w
}

func (w *Wrapper) Init(enc zapcore.Encoder, ws zapcore.WriteSyncer, level LevelEnabler, options ...zap.Option) {
	w.init(false, enc, ws, level, options...)
}

func (w *Wrapper) InitFull(enc zapcore.Encoder, ws zapcore.WriteSyncer, level LevelEnabler, options ...zap.Option) {
	w.init(true, enc, ws, level, options...)
}

func (w *Wrapper) init(full bool, enc zapcore.Encoder, ws zapcore.WriteSyncer, level LevelEnabler, options ...zap.Option) {
	var core zapcore.Core

	if enc == nil || ws == nil || level == nil {
		core = zapcore.NewNopCore()
	} else {
		core = zapcore.NewCore(enc, ws, level)
	}

	l := zap.New(core, options...).Named(w.name)

	w.lock.Lock()
	w.encoder = enc
	w.writeSyncer = ws
	w.levelEnabler = level
	w.options = options

	w.logger = l
	w.sugar = l.Sugar()
	w.lock.Unlock()

	if full {
		for _, node := range w.Tree() {
			node.init(full, enc, ws, level, options...)
		}
	}
}

func (w *Wrapper) Encoder() zapcore.Encoder {
	w.lock.RLock()
	defer w.lock.RUnlock()

	return w.encoder
}

func (w *Wrapper) SetEncoder(full bool, enc zapcore.Encoder) {
	w.init(full, enc, w.WriteSyncer(), w.LevelEnabler())
}

func (w *Wrapper) WriteSyncer() zapcore.WriteSyncer {
	w.lock.RLock()
	defer w.lock.RUnlock()

	return w.writeSyncer
}

func (w *Wrapper) SetWriteSyncer(full bool, ws zapcore.WriteSyncer) {
	w.init(full, w.Encoder(), ws, w.LevelEnabler())
}

func (w *Wrapper) LevelEnabler() LevelEnabler {
	w.lock.RLock()
	defer w.lock.RUnlock()

	return w.levelEnabler
}

func (w *Wrapper) SetLevelEnabler(full bool, level LevelEnabler) {
	w.init(full, w.Encoder(), w.WriteSyncer(), level)
}

func (w *Wrapper) Options() []zap.Option {
	w.lock.RLock()
	defer w.lock.RUnlock()

	opts := make([]zap.Option, len(w.options))
	copy(opts, w.options)

	return opts
}

func (w *Wrapper) WithOptions(full bool, options ...zap.Option) {
	w.lock.RLock()
	opts := append(w.Options(), options...)
	w.lock.RUnlock()

	w.init(full, w.Encoder(), w.WriteSyncer(), w.LevelEnabler(), opts...)
}

func (w *Wrapper) Name() string {
	return w.name
}

func (w *Wrapper) Named(name string) Logger {
	return w.LoadOrStore(name)
}

func (w *Wrapper) LoadOrStore(name string) *Wrapper {
	w.lock.RLock()
	exists, ok := w.tree[name]
	w.lock.RUnlock()

	if ok {
		return exists
	}

	l := New(name, w.Encoder(), w.WriteSyncer(), w.LevelEnabler(), w.Options()...)

	w.lock.Lock()
	w.tree[name] = l
	w.lock.Unlock()

	return l
}

func (w *Wrapper) Tree() map[string]*Wrapper {
	w.lock.RLock()
	defer w.lock.RUnlock()

	result := make(map[string]*Wrapper, len(w.tree))
	for k, v := range w.tree {
		result[k] = v
	}

	return result
}

func (w *Wrapper) Logger() *zap.Logger {
	w.lock.RLock()
	defer w.lock.RUnlock()

	return w.logger
}

func (w *Wrapper) Sugar() *zap.SugaredLogger {
	w.lock.RLock()
	defer w.lock.RUnlock()

	return w.sugar
}

func (w *Wrapper) Debug(message string, args ...interface{}) {
	w.Sugar().Debugw(message, args...)
}

func (w *Wrapper) Debugf(template string, args ...interface{}) {
	w.Sugar().Debugf(template, args...)
}

func (w *Wrapper) Info(message string, args ...interface{}) {
	w.Sugar().Infow(message, args...)
}

func (w *Wrapper) Infof(template string, args ...interface{}) {
	w.Sugar().Infof(template, args...)
}

func (w *Wrapper) Warn(message string, args ...interface{}) {
	w.Sugar().Warnw(message, args...)
}

func (w *Wrapper) Warnf(template string, args ...interface{}) {
	w.Sugar().Warnf(template, args...)
}

func (w *Wrapper) Error(message string, args ...interface{}) {
	w.Sugar().Errorw(message, args...)
}

func (w *Wrapper) Errorf(template string, args ...interface{}) {
	w.Sugar().Errorf(template, args...)
}

func (w *Wrapper) Panic(message string, args ...interface{}) {
	w.Sugar().Panicw(message, args...)
}

func (w *Wrapper) Panicf(template string, args ...interface{}) {
	w.Sugar().Panicf(template, args...)
}

func (w *Wrapper) Fatal(message string, args ...interface{}) {
	w.Sugar().Fatalw(message, args...)
}

func (w *Wrapper) Fatalf(template string, args ...interface{}) {
	w.Sugar().Fatalf(template, args...)
}
