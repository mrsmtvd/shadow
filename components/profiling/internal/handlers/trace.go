package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/logging"
	"github.com/mrsmtvd/shadow/components/profiling"
	"github.com/mrsmtvd/shadow/components/profiling/trace"
)

type TraceHandler struct {
	dashboard.Handler
}

func (h *TraceHandler) actionStart(_ *dashboard.Response, r *dashboard.Request) error {
	if trace.IsStarted() {
		return errors.New("trace already started")
	}

	if err := r.Original().ParseForm(); err != nil {
		return err
	}

	profiles := trace.GetProfiles()
	runProfiles := make([]string, 0, len(profiles))

	for _, profile := range profiles {
		id := profile.GetID()

		if r.Original().PostForm.Get("profile_"+id) != "" {
			runProfiles = append(runProfiles, id)
			logging.Log(r.Context()).Info("Run trace " + id)
		}
	}

	if len(runProfiles) == 0 {
		return errors.New("nothing to start")
	}

	err := trace.StartProfiles(runProfiles)
	logging.Log(r.Context()).Info("Run trace: " + strings.Join(runProfiles, ", "))

	return err
}

func (h *TraceHandler) actionStop(_ *dashboard.Response, r *dashboard.Request) error {
	if !trace.IsStarted() {
		return errors.New("trace already stoped")
	}

	err := trace.StopProfiles(r.Config().String(profiling.ConfigDumpDirectory))
	logging.Log(r.Context()).Info("Stop trace")

	return err
}

func (h *TraceHandler) actionDownload(w *dashboard.Response, r *dashboard.Request) error {
	id := r.URL().Query().Get("id")
	if id == "" {
		return errors.New("dump \"" + id + "\" not found")
	}

	dump := trace.GetDump(id)
	if dump == nil {
		return errors.New("dump \"" + id + "\" not found")
	}

	file, err := os.Open(dump.GetFile())
	if err != nil {
		return err
	}
	defer file.Close()

	w.Header().Set("Content-Length", strconv.FormatInt(dump.GetSize(), 10))
	w.Header().Set("Content-Type", "application/x-gzip")
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(dump.GetFile()))

	_, _ = io.Copy(w, file)

	return nil
}

func (h *TraceHandler) actionDelete(_ *dashboard.Response, r *dashboard.Request) error {
	id := r.URL().Query().Get("id")
	if id == "" {
		return errors.New("dump \"" + id + "\" not found")
	}

	if id == "all" {
		dumps := trace.GetDumps()
		for _, dump := range dumps {
			if err := trace.DeleteDump(dump.GetID()); err != nil {
				return err
			}

			logging.Log(r.Context()).Info("Remove " + dump.GetID() + " dump from file " + dump.GetFile())
		}

		return nil
	}

	dump := trace.GetDump(id)
	if dump == nil {
		return errors.New("dump \"" + id + "\" not found")
	}

	err := trace.DeleteDump(id)
	logging.Log(r.Context()).Info("Remove " + id + " dump from file " + dump.GetFile())

	return err
}

func (h *TraceHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if !r.Config().Bool(config.ConfigDebug) {
		h.NotFound(w, r)
		return
	}

	var err error

	action := r.Original().URL.Query().Get("action")

	if r.IsPost() {
		switch action {
		case "start":
			err = h.actionStart(w, r)
		case "stop":
			err = h.actionStop(w, r)
		case "delete":
			err = h.actionDelete(w, r)
		}

		if err == nil {
			redirectURL := &url.URL{}
			*redirectURL = *r.Original().URL
			redirectURL.RawQuery = ""

			h.Redirect(redirectURL.String(), http.StatusFound, w, r)

			return
		}
	} else if action == "download" {
		if err = h.actionDownload(w, r); err != nil {
			logging.Log(r.Context()).Error("Error in download trace: %s", err.Error())
		}

		return
	}

	dumps := trace.GetDumps()
	started := trace.GetStarted()
	context := map[string]interface{}{
		"dumps":      dumps,
		"profiles":   trace.GetProfiles(),
		"started":    started,
		"duration":   0,
		"remove_all": len(dumps) != 0,
	}

	if err != nil {
		r.Session().FlashBag().Error(err.Error())
	}

	if started != nil {
		context["duration"] = time.Since(*started)
	}

	for _, dump := range dumps {
		if dump.GetStatus() == trace.DumpStatusPrepare {
			context["remove_all"] = false
			break
		}
	}

	h.Render(r.Context(), "trace", context)
}
