package handlers

import (
	"encoding/hex"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/ota"
	"github.com/kihamo/shadow/components/ota/release"
	"github.com/kihamo/shadow/components/ota/repository"
)

type releaseView struct {
	ID           string
	Version      string
	Size         int64
	Checksum     string
	IsCurrent    bool
	IsRemovable  bool
	Path         string
	Architecture string
	UploadedAt   *time.Time
	DownloadURL  string
}

type response struct {
	Result  string `json:"result"`
	Message string `json:"message,omitempty"`
}

type ReleasesHandler struct {
	dashboard.Handler

	Updater           *ota.Updater
	UploadRepository  *repository.Directory
	UpgradeRepository ota.Repository
	CurrentRelease    ota.Release
}

func (h *ReleasesHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()

	releases, err := h.UpgradeRepository.Releases("")
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
			ID:           release.GenerateReleaseID(rl),
			Version:      rl.Version(),
			Size:         rl.Size(),
			Checksum:     hex.EncodeToString(rl.Checksum()),
			IsCurrent:    rl == h.CurrentRelease,
			Architecture: rl.Architecture(),
			Path:         rl.Path(),
		}
		rView.DownloadURL = "/ota/repository/" + rView.ID + "/" + generateFileName(rl)

		if releaseFile, ok := rl.(*release.LocalFile); ok {
			rView.UploadedAt = &[]time.Time{releaseFile.FileInfo().ModTime()}[0]
			rView.IsRemovable = true
		}

		releasesView = append(releasesView, rView)
	}

	h.Render(r.Context(), "releases", map[string]interface{}{
		"releases":    releasesView,
		"currentArch": runtime.GOARCH,
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
			if rlID := release.GenerateReleaseID(rl); rlID == id {
				if rl == h.CurrentRelease {
					_ = w.SendJSON(response{
						Result:  "failed",
						Message: "can't remove current release",
					})
					return
				}

				h.UploadRepository.Remove(rl)
				info := []interface{}{"version", rl.Version()}

				if releaseFile, ok := rl.(*release.LocalFile); ok {
					os.Remove(releaseFile.Path())
					info = append(info, "path", releaseFile.Path())
				}

				logging.Log(r.Context()).Info("Remove release", info...)

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
			if rlID := release.GenerateReleaseID(rl); rlID == id {
				err = h.Updater.Update(rl)
				if err != nil {
					r.Session().FlashBag().Error(err.Error())
				} else {
					info := []interface{}{"version", rl.Version()}

					if releaseFile, ok := rl.(*release.LocalFile); ok {
						info = append(info, "path", releaseFile.Path())
					}

					logging.Log(r.Context()).Info("Release upgrade", info...)
				}

				if r.URL().Query().Get("restart") != "" {
					err = h.Updater.Restart()
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
