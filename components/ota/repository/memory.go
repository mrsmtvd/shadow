package repository

import (
	"sync"

	"github.com/mrsmtvd/shadow/components/ota"
)

type Memory struct {
	lock     sync.RWMutex
	releases []ota.Release
}

func NewMemory(releases ...ota.Release) *Memory {
	return &Memory{
		releases: releases,
	}
}

func (r *Memory) Add(release ota.Release) {
	r.lock.Lock()
	r.releases = append(r.releases, release)
	r.lock.Unlock()
}

func (r *Memory) Remove(release ota.Release) error {
	r.lock.Lock()

	for i, rl := range r.releases {
		if release == rl {
			r.releases = append(r.releases[:i], r.releases[i+1:]...)
			break
		}
	}

	r.lock.Unlock()

	return nil
}

func (r *Memory) CanRemove(release ota.Release) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, rl := range r.releases {
		if release == rl {
			return true
		}
	}

	return false
}

func (r *Memory) Releases(arch string) ([]ota.Release, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	releases := make([]ota.Release, 0, len(r.releases))
	for _, release := range r.releases {
		if arch != "" && release.Architecture() != arch && release.Architecture() != ota.ArchitectureUnknown {
			continue
		}

		releases = append(releases, release)
	}

	return releases, nil
}

func (r *Memory) Update() error {
	return nil
}

func (r *Memory) Clean() {
	r.lock.Lock()
	r.releases = nil
	r.lock.Unlock()
}
