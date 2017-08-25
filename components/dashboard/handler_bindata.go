package dashboard

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow"
)

type BindataHandler struct {
	Handler

	application shadow.Application
}

type bindataList struct {
	IsDir   bool
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	Path    string
	Reader  io.Reader
}

type bindataBreadcrumb struct {
	Name   string
	Path   string
	Active bool
}

func (h *BindataHandler) getRoot() ([]bindataList, error) {
	files := []bindataList{}

	components, err := h.application.GetComponents()
	if err != nil {
		return nil, err
	}

	modTime := time.Now()
	if h.application.GetBuildDate() != nil {
		modTime = *h.application.GetBuildDate()
	}

	for _, component := range components {
		if componentTemplate, ok := component.(hasTemplate); ok {
			fs := componentTemplate.GetTemplates()
			if fs == nil {
				continue
			}

			files = append(files, bindataList{
				IsDir:   true,
				Name:    component.GetName(),
				Size:    0,
				Mode:    os.FileMode(0644) | os.ModeDir,
				ModTime: modTime,
				Path:    filepath.Join("/", component.GetName()),
			})
		}
	}

	return files, nil
}

func (h *BindataHandler) getComponentByPath(name, path string) ([]bindataList, error) {
	if !h.application.HasComponent(name) {
		return nil, fmt.Errorf("Component %s not found", name)
	}

	componentTemplate, ok := h.application.GetComponent(name).(hasTemplate)
	if !ok {
		return nil, fmt.Errorf("Component %s haven't templates", name)
	}

	fs := componentTemplate.GetTemplates()
	if fs.Prefix != "" {
		fs.Prefix = ""
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
		return []bindataList{{
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
		return nil, fmt.Errorf("Failed cast %s to assetfs.AssetDirectory", filepath.Join(path, statRoot.Name()))
	}

	files, err := fsDirectory.Readdir(0)
	if err != nil {
		return nil, err
	}

	ret := make([]bindataList, 0, len(files))

	for _, file := range files {
		fileSub, err := fs.Open(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}

		statSub, err := fileSub.Stat()
		if err != nil {
			return nil, err
		}

		infoSub := bindataList{
			IsDir:   statSub.IsDir(),
			Name:    statSub.Name(),
			Size:    statSub.Size(),
			Mode:    statSub.Mode(),
			ModTime: statSub.ModTime(),
			Path:    filepath.Join("/", name, path, statSub.Name()),
			Reader:  fileSub,
		}

		if statSub.IsDir() {
			if h.application.GetBuildDate() != nil {
				infoSub.ModTime = *h.application.GetBuildDate()
			}

			// TODO: directory size
		}

		ret = append(ret, infoSub)
	}

	return ret, nil
}

func (h *BindataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		files      []bindataList
		breadcrumb []bindataBreadcrumb
		err        error
	)

	sep := string(os.PathSeparator)
	path := r.URL.Query().Get("path")
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

		switch r.URL.Query().Get("mode") {
		case "raw":
			if len(files) == 1 {
				io.Copy(w, files[0].Reader)
				return
			}
			break
		case "file":
			if len(files) == 1 {
				w.Header().Set("Content-Length", strconv.FormatInt(files[0].Size, 10))
				w.Header().Set("Content-Type", "application/x-gzip")
				w.Header().Set("Content-Disposition", "attachment; filename="+files[0].Name)

				io.Copy(w, files[0].Reader)
				return
			}
			break
		}
	}

	// breadcrumbs
	parts := strings.Split(strings.TrimLeft(path, sep), sep)
	prefix := sep
	breadcrumb = make([]bindataBreadcrumb, 0, len(parts)+1)
	breadcrumb = append(breadcrumb, bindataBreadcrumb{
		Name: "Root",
		Path: prefix,
	})

	for _, name := range parts {
		prefix = filepath.Join(prefix, name)

		breadcrumb = append(breadcrumb, bindataBreadcrumb{
			Name: name,
			Path: prefix,
		})
	}

	breadcrumb[len(breadcrumb)-1].Active = true

	h.Render(r.Context(), ComponentName, "bindata", map[string]interface{}{
		"breadcrumb": breadcrumb,
		"files":      files,
		"error":      err,
	})
}
