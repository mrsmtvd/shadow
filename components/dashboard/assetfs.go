package dashboard

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow"
)

const (
	AssetFSPrefixRoute     = "assets"
	AssetFSPrefixTemplates = "templates"
)

type AssetsHandler struct {
	Handler

	root http.FileSystem
}

func NewAssetsHandler(root http.FileSystem) *AssetsHandler {
	return &AssetsHandler{
		root: root,
	}
}

func (h *AssetsHandler) ServeHTTP(w *Response, r *Request) {
	path := r.URL().Query().Get(":filepath")

	f, err := h.root.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			h.NotFound(w, r)
			return
		}

		panic(err.Error())
	}

	d, err := f.Stat()
	if err != nil {
		panic(err.Error())
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
			panic(err)
		}
	}
	w.Header().Set("Content-Type", ctype)
	w.Header().Set("Cache-Control", "max-age=315360000, public, immutable")
	w.Header().Set("Last-Modified", d.ModTime().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Length", strconv.FormatInt(d.Size(), 10))
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, f)

	if err != nil {
		panic(err.Error())
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
