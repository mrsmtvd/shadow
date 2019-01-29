package handlers

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

type AssetFSHandler struct {
	dashboard.Handler

	registry  *sync.Map
	buildDate *time.Time
}

type assetFSList struct {
	IsDir   bool
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	Path    string
	Reader  io.Reader
}

type assetFSBreadcrumb struct {
	Name   string
	Path   string
	Active bool
}

func NewAssetFSHandler(registry *sync.Map, buildDate *time.Time) *AssetFSHandler {
	return &AssetFSHandler{
		registry:  registry,
		buildDate: buildDate,
	}
}

func (h *AssetFSHandler) getRoot() ([]assetFSList, error) {
	files := make([]assetFSList, 0, 0)

	modTime := time.Now()
	if h.buildDate != nil {
		modTime = *h.buildDate
	}

	h.registry.Range(func(key, value interface{}) bool {
		name := key.(string)

		files = append(files, assetFSList{
			IsDir:   true,
			Name:    name,
			Size:    0,
			Mode:    os.FileMode(0644) | os.ModeDir,
			ModTime: modTime,
			Path:    filepath.Join("/", name),
		})

		return true
	})

	return files, nil
}

func (h *AssetFSHandler) getComponentByPath(name, path string) ([]assetFSList, error) {
	var fs *assetfs.AssetFS

	h.registry.Range(func(key, value interface{}) bool {
		if key.(string) == name {
			fs = value.(*assetfs.AssetFS)
			return false
		}

		return true
	})

	if fs == nil {
		return nil, errors.New("directory " + name + " not found")
	}

	fileRoot, err := fs.Open(path)
	if err != nil {
		return nil, err
	}

	statRoot, err := fileRoot.Stat()
	if err != nil {
		return nil, err
	}

	if !statRoot.IsDir() {
		return []assetFSList{{
			IsDir:   false,
			Name:    statRoot.Name(),
			Size:    statRoot.Size(),
			Mode:    statRoot.Mode(),
			ModTime: statRoot.ModTime(),
			Path:    filepath.Join("/", name, path),
			Reader:  fileRoot,
		}}, nil
	}

	fsDirectory, ok := statRoot.(*assetfs.AssetDirectory)
	if !ok {
		return nil, errors.New("failed cast " + filepath.Join(path, statRoot.Name()) + " to assetfs.AssetDirectory")
	}

	files, err := fsDirectory.Readdir(0)
	if err != nil {
		return nil, err
	}

	ret := make([]assetFSList, 0, len(files))

	for _, file := range files {
		fileSub, err := fs.Open(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}

		statSub, err := fileSub.Stat()
		if err != nil {
			return nil, err
		}

		infoSub := assetFSList{
			IsDir:   statSub.IsDir(),
			Name:    statSub.Name(),
			Size:    statSub.Size(),
			Mode:    statSub.Mode(),
			ModTime: statSub.ModTime(),
			Path:    filepath.Join("/", name, path, statSub.Name()),
			Reader:  fileSub,
		}

		if statSub.IsDir() {
			if h.buildDate != nil {
				infoSub.ModTime = *h.buildDate
			}

			// TODO: directory size
		}

		ret = append(ret, infoSub)
	}

	return ret, nil
}

func (h *AssetFSHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	var (
		files      []assetFSList
		breadcrumb []assetFSBreadcrumb
		err        error
	)

	sep := string(os.PathSeparator)
	path := r.URL().Query().Get("path")
	path = strings.Trim(filepath.Clean(path), sep)

	if path == "" || path == "." {
		path = sep
	}

	if path == sep {
		files, err = h.getRoot()
	} else {
		dir, file := filepath.Split(path)
		if dir == "" {
			files, err = h.getComponentByPath(file, sep)
		} else {
			parts := strings.Split(dir, sep)
			if len(parts) > 1 {
				files, err = h.getComponentByPath(parts[0], filepath.Join(filepath.Join(parts[1:]...), file))
			}
		}
	}

	if err == nil {
		sort.SliceStable(files, func(i, j int) bool {
			return files[i].Name < files[j].Name
		})

		sort.SliceStable(files, func(i, j int) bool {
			return files[i].IsDir != files[j].IsDir
		})

		switch r.URL().Query().Get("mode") {
		case "raw":
			if len(files) == 1 {
				if _, err := io.Copy(w, files[0].Reader); err != nil {
					h.InternalError(w, r, err)
				}

				return
			}

		case "file":
			if len(files) == 1 {
				w.Header().Set("Content-Length", strconv.FormatInt(files[0].Size, 10))
				w.Header().Set("Content-Type", "application/x-gzip")
				w.Header().Set("Content-Disposition", "attachment; filename="+files[0].Name)

				if _, err := io.Copy(w, files[0].Reader); err != nil {
					h.InternalError(w, r, err)
				}

				return
			}
		}
	}

	// breadcrumbs
	parts := strings.Split(strings.TrimLeft(path, sep), sep)
	prefix := sep
	breadcrumb = make([]assetFSBreadcrumb, 0, len(parts)+1)
	breadcrumb = append(breadcrumb, assetFSBreadcrumb{
		Name: "Root",
		Path: prefix,
	})

	for _, name := range parts {
		prefix = filepath.Join(prefix, name)

		breadcrumb = append(breadcrumb, assetFSBreadcrumb{
			Name: name,
			Path: prefix,
		})
	}

	breadcrumb[len(breadcrumb)-1].Active = true

	h.Render(r.Context(), "assetfs", map[string]interface{}{
		"breadcrumb": breadcrumb,
		"files":      files,
		"error":      err,
	})
}
