package trace

import (
	"bytes"
	"runtime/pprof"
	"sort"
	"sync"
)

type Profile struct {
	Id          string
	Description string
	Started     bool
	Buffer      *bytes.Buffer
}

var profiles struct {
	mutex    sync.RWMutex
	once     sync.Once
	profiles map[string]*Profile
}

func init() {
	profiles.profiles = make(map[string]*Profile)
}

func loadProfilesOnce() {
	profiles.mutex.Lock()
	defer profiles.mutex.Unlock()

	profiles.profiles = map[string]*Profile{
		ProfileCpu: {
			Id:          ProfileCpu,
			Description: "CPU profiling for the current process",
			Buffer:      bytes.NewBuffer(nil),
		},
		ProfileTrace: {
			Id:          ProfileTrace,
			Description: "tracing for the current program",
			Buffer:      bytes.NewBuffer(nil),
		},
	}

	for _, p := range pprof.Profiles() {
		profile := Profile{
			Id:     p.Name(),
			Buffer: bytes.NewBuffer(nil),
		}

		switch profile.Id {
		case "goroutine":
			profile.Description = "stack traces of all current goroutines"
		case "heap":
			profile.Description = "a sampling of all heap allocations"
		case "threadcreate":
			profile.Description = "stack traces that led to the creation of new OS threads"
		case "block":
			profile.Description = "stack traces that led to blocking on synchronization primitives"
		case "mutex":
			profile.Description = "stack traces of holders of contended mutexes"
		default:
			profile.Description = profile.Id
		}

		profiles.profiles[profile.Id] = &profile
	}
}

func LoadProfiles() {
	profiles.once.Do(loadProfilesOnce)
}

func GetProfiles() []Profile {
	LoadProfiles()

	profiles.mutex.RLock()
	defer profiles.mutex.RUnlock()

	result := make([]Profile, 0, len(profiles.profiles))
	for _, profile := range profiles.profiles {
		result = append(result, *profile)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Id < result[j].Id
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
