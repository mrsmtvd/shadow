package trace

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

const (
	ProfileCpu   = "cpu"
	ProfileTrace = "trace"
)

var startAt *time.Time

func GetStarted() *time.Time {
	return startAt
}

func IsStarted() bool {
	return startAt != nil
}

func StartProfiles(list []string) error {
	if IsStarted() {
		return fmt.Errorf("Trace already started")
	}

	LoadProfiles()

	runProfiles := make([]*Profile, 0, len(list))

	profiles.mutex.RLock()
	defer profiles.mutex.RUnlock()

	for _, id := range list {
		profile, ok := profiles.profiles[id]
		if !ok {
			return fmt.Errorf("Profile \"%s\" not found", id)
		}

		runProfiles = append(runProfiles, profile)
	}

	if len(runProfiles) == 0 {
		return nil
	}

	now := time.Now()
	for i := range runProfiles {
		switch runProfiles[i].Id {
		case ProfileCpu:
			if err := pprof.StartCPUProfile(runProfiles[i].Buffer); err != nil {
				return err
			}
		case ProfileTrace:
			if err := trace.Start(runProfiles[i].Buffer); err != nil {
				return err
			}
		}

		runProfiles[i].Started = true
	}

	startAt = &now
	return nil
}

func StopProfiles(path string) error {
	if !IsStarted() {
		return fmt.Errorf("Trace already stoped")
	}

	dir, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !dir.IsDir() {
		return fmt.Errorf("%s isn't directory", path)
	}

	profiles.mutex.RLock()
	defer profiles.mutex.RUnlock()

	hash := md5.New()
	io.WriteString(hash, startAt.String())

	dump := &Dump{
		id:        hex.EncodeToString(hash.Sum(nil)),
		file:      filepath.Join(path, startAt.Format("20060102150405.tar.gz")),
		startedAt: *startAt,
		stoppedAt: time.Now(),
		profiles:  make([]Profile, 0, len(profiles.profiles)),
		status:    DumpStatusPrepare,
	}

	for _, profile := range profiles.profiles {
		if !profile.Started {
			continue
		}

		switch profile.Id {
		case ProfileCpu:
			pprof.StopCPUProfile()
		case ProfileTrace:
			trace.Stop()
		default:
			pprof.Lookup(profile.Id).WriteTo(profile.Buffer, 0)
		}

		dump.AddProfile(*profile)
		profile.Started = false
	}

	dumps.dumps[dump.GetId()] = dump
	startAt = nil

	go saveArchive(dump)

	return nil
}
