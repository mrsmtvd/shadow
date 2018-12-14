package trace

import (
	"bytes"
	"runtime/pprof"
	"sort"
	"sync"
)

type Profile struct {
	mutex sync.RWMutex

	id          string
	description string
	started     bool
	buffer      *bytes.Buffer
}

var profiles struct {
	mutex    sync.RWMutex
	once     sync.Once
	profiles map[string]*Profile
}

func init() {
	profiles.profiles = make(map[string]*Profile)
}

func (p *Profile) GetId() string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.id
}

func (p *Profile) GetDescription() string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.description
}

func (p *Profile) GetStarted() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.started
}

func (p *Profile) SetStarted(started bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.started = started
}

func (p *Profile) Write(b []byte) (n int, err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.buffer.Write(b)
}

func loadProfilesOnce() {
	profiles.mutex.Lock()
	defer profiles.mutex.Unlock()

	profiles.profiles = map[string]*Profile{
		ProfileCpu: {
			id:          ProfileCpu,
			description: "CPU profiling for the current process",
			buffer:      bytes.NewBuffer(nil),
		},
		ProfileTrace: {
			id:          ProfileTrace,
			description: "tracing for the current program",
			buffer:      bytes.NewBuffer(nil),
		},
	}

	for _, p := range pprof.Profiles() {
		profile := Profile{
			id:     p.Name(),
			buffer: bytes.NewBuffer(nil),
		}

		switch profile.id {
		case "goroutine":
			profile.description = "stack traces of all current goroutines"
		case "heap":
			profile.description = "a sampling of all heap allocations"
		case "threadcreate":
			profile.description = "stack traces that led to the creation of new OS threads"
		case "block":
			profile.description = "stack traces that led to blocking on synchronization primitives"
		case "mutex":
			profile.description = "stack traces of holders of contended mutexes"
		default:
			profile.description = profile.id
		}

		profiles.profiles[profile.id] = &profile
	}
}

func LoadProfiles() {
	profiles.once.Do(loadProfilesOnce)
}

func GetProfiles() []*Profile {
	LoadProfiles()

	profiles.mutex.RLock()
	list := profiles.profiles
	profiles.mutex.RUnlock()

	result := make([]*Profile, 0, len(list))
	for _, profile := range list {
		result = append(result, profile)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].GetId() < result[j].GetId()
	})

	return result
}

func GetProfile(id string) *Profile {
	LoadProfiles()

	profiles.mutex.RLock()
	defer profiles.mutex.RUnlock()

	if profile, ok := profiles.profiles[id]; ok {
		return profile
	}

	return nil
}
