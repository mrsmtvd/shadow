package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/profiling"
	"github.com/kihamo/shadow/components/profiling/trace"
)

type TraceHandler struct {
	dashboard.Handler
}

func (h *TraceHandler) actionStart(_ *dashboard.Response, r *dashboard.Request) error {
	if trace.IsStarted() {
		return fmt.Errorf("Trace already started")
	}

	if err := r.Original().ParseForm(); err != nil {
		return err
	}

	runProfiles := []string{}
	for _, profile := range trace.GetProfiles() {
		id := profile.GetId()

		if r.Original().PostForm.Get("profile_"+id) != "" {
			runProfiles = append(runProfiles, id)
			r.Logger().Infof("Run trace \"%s\"", id)
		}
	}

	if len(runProfiles) == 0 {
		return fmt.Errorf("Nothing to start")
	}

	err := trace.StartProfiles(runProfiles)
	r.Logger().Infof("Run trace: %s", strings.Join(runProfiles, ", "))

	return err
}

func (h *TraceHandler) actionStop(_ *dashboard.Response, r *dashboard.Request) error {
	if !trace.IsStarted() {
		return fmt.Errorf("Trace already stoped")
	}

	err := trace.StopProfiles(r.Config().String(profiling.ConfigDumpDirectory))
	r.Logger().Info("Stop trace")

	return err
}

func (h *TraceHandler) actionDownload(w *dashboard.Response, r *dashboard.Request) error {
	id := r.URL().Query().Get("id")
	if id == "" {
		return fmt.Errorf("Dump \"%s\" not found", id)
	}

	dump := trace.GetDump(id)
	if dump == nil {
		return fmt.Errorf("Dump \"%s\" not found", id)
	}

	file, err := os.Open(dump.GetFile())
	if err != nil {
		return err
	}
	defer file.Close()

	w.Header().Set("Content-Length", strconv.FormatInt(dump.GetSize(), 10))
	w.Header().Set("Content-Type", "application/x-gzip")
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(dump.GetFile()))

	io.Copy(w, file)

	return nil
}

func (h *TraceHandler) actionDelete(_ *dashboard.Response, r *dashboard.Request) error {
	id := r.URL().Query().Get("id")
	if id == "" {
		return fmt.Errorf("Dump \"%s\" not found", id)
	}

	if id == "all" {
		for _, dump := range trace.GetDumps() {
			if err := trace.DeleteDump(dump.GetId()); err != nil {
				return err
			}

			r.Logger().Infof("Remove %s dump from file %s", dump.GetId(), dump.GetFile())
		}

		return nil
	}

	dump := trace.GetDump(id)
	if dump == nil {
		return fmt.Errorf("Dump \"%s\" not found", id)
	}

	err := trace.DeleteDump(id)
	r.Logger().Infof("Remove %s dump from file %s", id, dump.GetFile())

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
			redirectUrl := &url.URL{}
			*redirectUrl = *r.Original().URL
			redirectUrl.RawQuery = ""

			h.Redirect(redirectUrl.String(), http.StatusFound, w, r)
			return
		}

	} else if action == "download" {
		if err = h.actionDownload(w, r); err != nil {
			r.Logger().Error("Error in download trace: %s", err.Error())
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
		"error":      err,
	}

	if started != nil {
		context["duration"] = time.Now().Sub(*started)
	}

	for _, dump := range dumps {
		if dump.GetStatus() == trace.DumpStatusPrepare {
			context["remove_all"] = false
		}
	}

	h.Render(r.Context(), profiling.ComponentName, "trace", context)
}
