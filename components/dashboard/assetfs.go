package dashboard

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow"
)

const (
	AssetFSPrefixRoute     = "assets"
	AssetFSPrefixTemplates = "templates"
)

type AssetsHandler struct {
	Handler

	root http.FileSystem
	path string
}

func NewAssetsHandler(root http.FileSystem) *AssetsHandler {
	return &AssetsHandler{
		root: root,
	}
}

func NewAssetsHandlerByPath(root http.FileSystem, path string) *AssetsHandler {
	return &AssetsHandler{
		root: root,
		path: path,
	}
}

func (h *AssetsHandler) ServeHTTP(w http.ResponseWriter, r *Request) {
	var path string
	if h.path != "" {
		path = h.path
	} else {
		path = r.URL().Query().Get(":filepath")
	}

	if path == "" {
		h.NotFound(w, r)
		return
	}

	f, err := h.root.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			h.NotFound(w, r)
			return
		}

		h.InternalError(w, r, err)

		return
	}

	d, err := f.Stat()
	if err != nil {
		h.InternalError(w, r, err)
		return
	}

	if d.IsDir() {
		h.NotFound(w, r)
		return
	}

	ctype := mime.TypeByExtension(filepath.Ext(d.Name()))
	if ctype == "" {
		var buf [512]byte
		n, _ := io.ReadFull(f, buf[:])
		ctype = http.DetectContentType(buf[:n])
		_, err := f.Seek(0, io.SeekStart)

		if err != nil {
			h.InternalError(w, r, err)
			return
		}
	}

	if ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}

	w.Header().Set("Cache-Control", "max-age=315360000, public, immutable")
	w.Header().Set("Last-Modified", d.ModTime().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Length", strconv.FormatInt(d.Size(), 10))
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, f)

	if err != nil {
		h.InternalError(w, r, err)
	}
}

type HasAssetFS interface {
	AssetFS() *assetfs.AssetFS
}

func RouteFromAssetFS(component HasAssetFS) Route {
	fs := component.AssetFS()
	fs.Prefix = AssetFSPrefixRoute

	return NewRoute("/"+component.(shadow.Component).Name()+"/assets/*filepath", NewAssetsHandler(fs)).
		WithMethods([]string{http.MethodGet})
}

func TemplatesFromAssetFS(component HasAssetFS) *assetfs.AssetFS {
	fs := component.AssetFS()
	fs.Prefix = AssetFSPrefixTemplates

	return fs
}
