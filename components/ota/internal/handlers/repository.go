package handlers

import (
	"encoding/hex"
	"io"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/ota"
	"github.com/kihamo/shadow/components/ota/release"
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
			if rlID := release.GenerateReleaseID(rl); rlID == id {
				releaseBinFile, err := rl.BinFile()
				if err != nil {
					h.InternalError(w, r, err)
					return
				}

				fileName := r.URL().Query().Get(":file")
				if fileName == "" {
					fileName = generateFileName(rl)
				}

				w.Header().Set("Content-Length", strconv.FormatInt(rl.Size(), 10))
				w.Header().Set("Content-Type", "application/x-binary")
				w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
				io.Copy(w, releaseBinFile)
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

	type item struct {
		Architecture string `json:"architecture"`
		Checksum     string `json:"checksum"`
		Size         int64  `json:"size"`
		Version      string `json:"version"`
		File         string `json:"file"`
	}

	items := make([]item, 0, len(releases))
	for _, rl := range releases {
		fileURL := &url.URL{
			Scheme: "http",
			Host:   r.Original().Host,
			Path:   "/ota/repository/" + release.GenerateReleaseID(rl) + "/" + generateFileName(rl),
		}

		if r.Original().TLS != nil {
			fileURL.Scheme = "https"
		}

		items = append(items, item{
			Architecture: rl.Architecture(),
			Checksum:     hex.EncodeToString(rl.Checksum()),
			Size:         rl.Size(),
			Version:      rl.Version(),
			File:         fileURL.String(),
		})
	}

	_ = w.SendJSON(items)
}

func generateFileName(rl ota.Release) string {
	basePath := filepath.Base(rl.Path())

	if strings.HasSuffix(basePath, ".bin") {
		return basePath
	}

	return strings.ReplaceAll(basePath, " ", "_") +
		"." + strings.ReplaceAll(rl.Version(), " ", ".") +
		"." + rl.Architecture() + ".bin"
}
