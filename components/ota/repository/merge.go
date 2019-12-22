package repository

import (
	"errors"
	"sync"

	"github.com/kihamo/shadow/components/ota"
)

type Merge struct {
	mutex        sync.RWMutex
	repositories []ota.Repository
}

func NewMerge(repositories ...ota.Repository) *Merge {
	return &Merge{
		repositories: repositories,
	}
}

func (r *Merge) Merge(repositories ...ota.Repository) *Merge {
	r.mutex.Lock()
	r.repositories = append(r.repositories, repositories...)
	r.mutex.Unlock()

	return r
}

func (r *Merge) Releases(arch string) ([]ota.Release, error) {
	releases := make([]ota.Release, 0)

	r.mutex.RLock()
	for _, repo := range r.repositories {
		items, err := repo.Releases(arch)
		if err != nil {
			return nil, err
		}
		releases = append(releases, items...)
	}
	r.mutex.RUnlock()

	return releases, nil
}

func (r *Merge) ReleaseLatest(arch string) (ota.Release, error) {
	releases, err := r.Releases(arch)
	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, errors.New("latest release not found")
	}

	return releases[len(releases)-1], nil
}
