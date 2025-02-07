package repository

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mrsmtvd/shadow/components/ota"
	"github.com/mrsmtvd/shadow/components/ota/release"
)

type Directory struct {
	*Memory

	dir string
}

func NewDirectory() *Directory {
	return &Directory{
		Memory: NewMemory(),
	}
}

func (r *Directory) SetPath(dir string) {
	r.lock.Lock()
	r.dir = dir
	r.lock.Unlock()
}

func (r *Directory) Remove(release ota.Release) (err error) {
	err = os.Remove(release.Path())

	if err == nil {
		err = r.Memory.Remove(release)
	}

	return err
}

func (r *Directory) CanRemove(release ota.Release) bool {
	return r.Memory.CanRemove(release)
}

func (r *Directory) Update() error {
	r.Memory.Clean()

	return filepath.Walk(r.dir, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}

		if !strings.HasPrefix(info.Name(), "release-file-") {
			return nil
		}

		rl, err := release.NewLocalFile(path, "")
		if err == nil {
			r.Add(release.NewCompress(rl))
		}

		return err
	})
}
