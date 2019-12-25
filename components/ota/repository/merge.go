package repository

import (
	"sync"

	"github.com/kihamo/shadow/components/ota"
	"go.uber.org/multierr"
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

func (r *Merge) Remove(release ota.Release) (err error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, repo := range r.repositories {
		if remover, ok := repo.(ota.RepositoryRemover); ok {
			err = remover.Remove(release)
			if err != nil {
				break
			}
		}
	}

	return err
}

func (r *Merge) CanRemove(release ota.Release) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, repo := range r.repositories {
		if remover, ok := repo.(ota.RepositoryRemover); ok {
			if remover.CanRemove(release) {
				return true
			}
		}
	}

	return false
}

func (r *Merge) Releases(arch string) ([]ota.Release, error) {
	releases := make([]ota.Release, 0)

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, repo := range r.repositories {
		items, err := repo.Releases(arch)
		if err != nil {
			return nil, err
		}
		releases = append(releases, items...)
	}

	return releases, nil
}

func (r *Merge) Update() (err error) {
	var (
		wg   sync.WaitGroup
		lock sync.Mutex
	)

	r.mutex.RLock()
	for _, repo := range r.repositories {
		wg.Add(1)

		go func(rp ota.Repository) {
			defer wg.Done()

			if e := rp.Update(); e != nil {
				lock.Lock()
				err = multierr.Append(err, e)
				lock.Unlock()
			}
		}(repo)
	}
	r.mutex.RUnlock()

	wg.Wait()

	return err
}
