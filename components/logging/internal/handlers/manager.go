package handlers

import (
	"net/http"
	"strconv"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logging/internal/wrapper"
	"go.uber.org/zap/zapcore"
)

const (
	idSeparator = "^"
)

type loggerView struct {
	Parent    string
	Name      string
	Level     int8
	Logger    *wrapper.Wrapper
	Dependent int
}

type levelView struct {
	Name  string
	Value int8
	Level zapcore.Level
}

type ManagerHandler struct {
	dashboard.Handler

	wrapper *wrapper.Wrapper
	levels  []levelView
}

func NewManagerHandler(wrapper *wrapper.Wrapper, levels [][]interface{}) *ManagerHandler {
	lvl := make([]levelView, len(levels))
	for i, v := range levels {
		lvl[i] = levelView{
			Value: v[0].(int8),
			Level: zapcore.Level(v[0].(int8)),
			Name:  v[1].(string),
		}
	}

	return &ManagerHandler{
		wrapper: wrapper,
		levels:  lvl,
	}
}

func (h *ManagerHandler) collect(w *wrapper.Wrapper) map[string]loggerView {
	name := w.Name()
	level := int8(wrapper.DebugLevel) - 1
	result := make(map[string]loggerView)
	dependent := 0

	if lvl := w.LevelEnabler(); lvl != nil {
		for _, v := range h.levels {
			if lvl.Enabled(v.Level) {
				level = v.Value
				break
			}
		}
	}

	for _, logger := range w.Tree() {
		dependent++

		for wrapName, wrapLogger := range h.collect(logger) {
			if wrapLogger.Parent == "" {
				wrapLogger.Parent = name
			}

			result[name+idSeparator+wrapName] = wrapLogger
		}
	}

	result[name] = loggerView{
		Name:      name,
		Level:     level,
		Logger:    w,
		Dependent: dependent,
	}

	return result
}

func (h *ManagerHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	var err error
	loggers := h.collect(h.wrapper)

	if r.IsPost() {
		err = r.Original().ParseForm()

		if err == nil {
			for key, value := range r.Original().PostForm {
				logger, ok := loggers[key]
				if !ok || len(value) == 0 {
					continue
				}

				levelValue, e := strconv.ParseInt(value[0], 10, 64)
				if e != nil {
					continue
				}

				if logger.Level == int8(levelValue) {
					continue
				}

				logger.Logger.SetLevelEnabler(false, zapcore.Level(levelValue))
				logger.Level = int8(levelValue)
			}
		}

		h.Redirect(r.URL().String(), http.StatusFound, w, r)
		return
	}

	if err != nil {
		r.Session().FlashBag().Error(err.Error())
	}

	h.Render(r.Context(), "manager", map[string]interface{}{
		"loggers": loggers,
		"levels":  h.levels,
	})
}
