package handlers

import (
	"encoding/hex"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/ota"
	"github.com/kihamo/shadow/components/ota/repository"
)

type RepositoryHandler struct {
	dashboard.Handler

	Repository ota.Repository
}

func (h *RepositoryHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	arch := strings.TrimSpace(r.URL().Query().Get("architecture"))

	releases, err := h.Repository.Releases(arch)
	if err != nil {
		h.InternalError(w, r, err)
		return
	}

	if id := r.URL().Query().Get(":id"); id != "" {
		for _, rl := range releases {
			if rlID := ota.GenerateReleaseID(rl); rlID == id {
				releaseBinFile, err := rl.File()
				if err != nil {
					logging.Log(r.Context()).Error("Get release download file failed ",
						"version", rl.Version(),
						"path", rl.Path(),
						"error", err.Error(),
					)

					h.NotFound(w, r)
					return
				}

				fileName := r.URL().Query().Get(":file")
				if fileName == "" {
					fileName = ota.GenerateFileName(rl)
				}

				w.Header().Set("Content-Length", strconv.FormatInt(rl.Size(), 10))
				w.Header().Set("Content-Type", rl.Type().MIME())
				w.Header().Set("Content-Disposition", "attachment; filename="+fileName)

				if !r.IsHead() {
					io.Copy(w, releaseBinFile)
				}

				releaseBinFile.Close()
				return
			}
		}

		h.NotFound(w, r)
		return
	}

	if !r.Config().Bool(ota.ConfigRepositoryServerEnabled) {
		h.NotFound(w, r)
		return
	}

	records := make([]repository.ShadowRecord, 0, len(releases))
	for _, rl := range releases {
		fileURL := &url.URL{
			Scheme: "http",
			Host:   r.Original().Host,
			Path:   "/ota/repository/" + ota.GenerateReleaseID(rl) + "/" + ota.GenerateFileName(rl),
		}

		if r.Original().TLS != nil {
			fileURL.Scheme = "https"
		}

		records = append(records, repository.ShadowRecord{
			Architecture: rl.Architecture(),
			Checksum:     hex.EncodeToString(rl.Checksum()),
			Size:         rl.Size(),
			Version:      rl.Version(),
			File:         fileURL.String(),
			CreatedAt:    rl.CreatedAt(),
		})
	}

	_ = w.SendJSON(records)
}
