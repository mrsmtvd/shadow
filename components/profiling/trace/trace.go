package trace

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"log"
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
		return errors.New("trace already started")
	}

	LoadProfiles()

	runProfiles := make([]*Profile, 0, len(list))

	profiles.mutex.Lock()
	defer profiles.mutex.Unlock()

	for _, id := range list {
		profile, ok := profiles.profiles[id]
		if !ok {
			return errors.New("profile \"" + id + "\" not found")
		}

		runProfiles = append(runProfiles, profile)
	}

	if len(runProfiles) == 0 {
		return nil
	}

	now := time.Now()
	for i := range runProfiles {
		switch runProfiles[i].GetId() {
		case ProfileCpu:
			if err := pprof.StartCPUProfile(runProfiles[i]); err != nil {
				return err
			}
		case ProfileTrace:
			if err := trace.Start(runProfiles[i]); err != nil {
				return err
			}
		}

		runProfiles[i].SetStarted(true)
	}

	startAt = &now
	return nil
}

func StopProfiles(path string) error {
	if !IsStarted() {
		return errors.New("trace already stopped")
	}

	dir, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !dir.IsDir() {
		return errors.New(path + " isn't directory")
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
		if !profile.GetStarted() {
			continue
		}

		switch profile.GetId() {
		case ProfileCpu:
			pprof.StopCPUProfile()
		case ProfileTrace:
			trace.Stop()
		default:
			pprof.Lookup(profile.GetId()).WriteTo(profile, 0)
		}

		dump.AddProfile(*profile)
		profile.SetStarted(false)
	}

	dumps.dumps[dump.GetId()] = dump
	startAt = nil

	go func() {
		if err := saveDump(dump); err != nil {
			log.Printf("Error save dump file %s with error %s", dump.GetFile(), err.Error())
			dump.SetStatus(DumpStatusError)
		} else {
			dump.SetStatus(DumpStatusFinished)
		}
	}()

	return nil
}
