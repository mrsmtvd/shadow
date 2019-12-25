package handlers

import (
	"encoding/hex"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/ota"
	"github.com/kihamo/shadow/components/ota/repository"
)

type releaseView struct {
	ID            string
	Version       string
	Size          int64
	Checksum      string
	IsCurrent     bool
	IsRemovable   bool
	IsUpgradeable bool
	Path          string
	Architecture  string
	UploadedAt    *time.Time
	DownloadURL   string
}

type response struct {
	Result  string `json:"result"`
	Message string `json:"message,omitempty"`
}

type ReleasesHandler struct {
	dashboard.Handler

	Installer      *ota.Installer
	AllRepository  *repository.Merge
	CurrentRelease ota.Release
}

func (h *ReleasesHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()

	if q.Get("update") != "" {
		if err := h.AllRepository.Update(); err != nil {
			r.Session().FlashBag().Error(err.Error())
		}
	}

	releases, err := h.AllRepository.Releases("")
	if err != nil {
		r.Session().FlashBag().Error(err.Error())
	} else {
		switch q.Get(":action") {
		case "remove":
			h.actionRemove(w, r, releases)
			return

		case "upgrade":
			h.actionUpgrade(w, r, releases)
			return
		}
	}

	releasesView := make([]releaseView, 0, len(releases))
	for _, rl := range releases {
		rView := releaseView{
			ID:            ota.GenerateReleaseID(rl),
			Version:       rl.Version(),
			Size:          rl.Size(),
			Checksum:      hex.EncodeToString(rl.Checksum()),
			IsCurrent:     rl == h.CurrentRelease,
			IsRemovable:   rl != h.CurrentRelease && h.AllRepository.CanRemove(rl),
			IsUpgradeable: rl != h.CurrentRelease && rl.Architecture() == runtime.GOARCH,
			Architecture:  rl.Architecture(),
			Path:          rl.Path(),
			UploadedAt:    rl.CreatedAt(),
		}
		rView.DownloadURL = "/ota/repository/" + rView.ID + "/" + ota.GenerateFileName(rl)

		releasesView = append(releasesView, rView)
	}

	h.Render(r.Context(), "releases", map[string]interface{}{
		"releases": releasesView,
	})
}

func (h *ReleasesHandler) actionRemove(w *dashboard.Response, r *dashboard.Request, releases []ota.Release) {
	if !r.IsPost() {
		h.MethodNotAllowed(w, r)
		return
	}

	id := strings.TrimSpace(r.URL().Query().Get(":id"))
	if id != "" {
		for _, rl := range releases {
			if rlID := ota.GenerateReleaseID(rl); rlID == id {
				if rl == h.CurrentRelease {
					_ = w.SendJSON(response{
						Result:  "failed",
						Message: "can't remove current release",
					})
					return
				}

				if h.AllRepository.CanRemove(rl) {
					if err := h.AllRepository.Remove(rl); err != nil {
						_ = w.SendJSON(response{
							Result:  "failed",
							Message: fmt.Sprintf("remove release failed %v", err),
						})
						return
					}

					logging.Log(r.Context()).Info("Remove release",
						"version", rl.Version(),
						"path", rl.Path(),
					)

					go h.AllRepository.Update()
				}

				_ = w.SendJSON(response{
					Result: "success",
				})
				return
			}
		}
	}

	h.NotFound(w, r)
}

func (h *ReleasesHandler) actionUpgrade(w *dashboard.Response, r *dashboard.Request, releases []ota.Release) {
	if !r.IsPost() {
		h.MethodNotAllowed(w, r)
		return
	}

	id := strings.TrimSpace(r.URL().Query().Get(":id"))
	if id != "" {
		var err error

		for _, rl := range releases {
			if rlID := ota.GenerateReleaseID(rl); rlID == id {
				err = h.Installer.Install(rl)
				if err != nil {
					r.Session().FlashBag().Error(err.Error())
				} else {
					logging.Log(r.Context()).Info("Release upgrade",
						"version", rl.Version(),
						"path", rl.Path(),
					)
				}

				if r.URL().Query().Get("restart") != "" {
					err = h.Installer.Restart()
				}

				if err != nil {
					_ = w.SendJSON(response{
						Result:  "failed",
						Message: err.Error(),
					})
				} else {
					_ = w.SendJSON(response{
						Result: "success",
					})
				}

				return
			}
		}
	}

	h.NotFound(w, r)
}
