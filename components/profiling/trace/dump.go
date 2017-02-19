package trace

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/kardianos/osext"
)

const (
	DumpStatusPrepare = iota
	DumpStatusFinished
	DumpStatusReading
	DumpStatusError
)

type Dump struct {
	mutex sync.RWMutex

	id        string
	file      string
	size      int64
	startedAt time.Time
	stoppedAt time.Time
	profiles  []Profile
	status    int
}

var filePattern = regexp.MustCompile(`^[0-9]{14}.tar.gz$`)
var tracePattern = regexp.MustCompile(`^(.*?).trace$`)

func (d *Dump) Delete() error {
	return os.Remove(d.file)
}

func (d *Dump) GetId() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.id
}

func (d *Dump) SetId(id string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.id = id
}

func (d *Dump) GetFile() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.file
}

func (d *Dump) GetSize() int64 {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.size
}

func (d *Dump) SetSize(size int64) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.size = size
}

func (d *Dump) GetStartedAt() time.Time {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.startedAt
}

func (d *Dump) GetStoppedAt() time.Time {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.stoppedAt
}

func (d *Dump) GetProfiles() []Profile {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.profiles
}

func (d *Dump) AddProfile(profile Profile) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.profiles = append(d.profiles, profile)
}

func (d *Dump) GetStatus() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.status
}

func (d *Dump) SetStatus(status int) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.status = status
}

var dumps struct {
	mutex sync.RWMutex
	dumps map[string]*Dump
}

func init() {
	dumps.dumps = make(map[string]*Dump)
}

func LoadDumps(path string) error {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if !filePattern.MatchString(f.Name()) {
			continue
		}

		filePath := filepath.Join(path, f.Name())

		hash := md5.New()
		io.WriteString(hash, filePath)

		dump := &Dump{
			id:        hex.EncodeToString(hash.Sum(nil)),
			file:      filePath,
			size:      f.Size(),
			startedAt: f.ModTime(),
			stoppedAt: f.ModTime(),
			profiles:  []Profile{},
			status:    DumpStatusReading,
		}

		dumps.mutex.Lock()
		dumps.dumps[dump.GetId()] = dump
		dumps.mutex.Unlock()

		go func() {
			if err := read(dump); err != nil {
				log.Printf("Error read %s with error %s", dump.GetFile(), err.Error())
				dump.SetStatus(DumpStatusError)
			} else {
				dump.SetStatus(DumpStatusFinished)
			}
		}()
	}

	return nil
}

func GetDumps() []Dump {
	dumps.mutex.RLock()
	defer dumps.mutex.RUnlock()

	result := make([]Dump, 0, len(dumps.dumps))
	for _, dump := range dumps.dumps {
		dump.mutex.RLock()
		result = append(result, *dump)
		dump.mutex.RUnlock()
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].stoppedAt.Sub(result[j].stoppedAt) > 0
	})

	return result
}

func GetDump(id string) *Dump {
	dumps.mutex.RLock()
	defer dumps.mutex.RUnlock()

	if dump, ok := dumps.dumps[id]; ok {
		return dump
	}

	return nil
}

func DeleteDump(id string) error {
	dumps.mutex.Lock()
	defer dumps.mutex.Unlock()

	dump, ok := dumps.dumps[id]
	if !ok {
		return fmt.Errorf("Dump \"%s\" not found", id)
	}

	if err := dump.Delete(); err != nil {
		return err
	}

	delete(dumps.dumps, id)
	return nil
}

func read(dump *Dump) error {
	file, err := os.Open(dump.GetFile())
	if err != nil {
		return err
	}
	defer file.Close()

	archiveReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer archiveReader.Close()

	tarReader := tar.NewReader(archiveReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		m := tracePattern.FindStringSubmatch(header.Name)
		if len(m) != 2 {
			continue
		}

		if profile := GetProfile(m[1]); profile != nil {
			dump.AddProfile(*profile)
		}

	}

	return nil
}

func saveArchive(dump *Dump) {
	file, err := os.Create(dump.GetFile())
	if err != nil {
		return
	}

	defer file.Close()
	defer func() {
		stat, _ := file.Stat()
		dump.SetSize(stat.Size())
	}()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// copy binary
	binary, err := osext.Executable()
	if err != nil {
		return
	}

	if err := writeFile(tarWriter, binary, filepath.Base(binary)); err != nil {
		return
	}

	// copy trace
	for _, profile := range dump.GetProfiles() {
		if profile.Buffer.Len() == 0 {
			continue
		}

		tmpfile, err := ioutil.TempFile("", "trace_")
		if err != nil {
			return
		}

		io.Copy(tmpfile, profile.Buffer)

		if err := writeFile(tarWriter, tmpfile.Name(), fmt.Sprintf("%s.trace", profile.Id)); err != nil {
			fmt.Println(err)

			tmpfile.Close()
			return
		}

		tmpfile.Close()
		profile.Buffer.Reset()
	}

	dump.SetStatus(DumpStatusFinished)
}

func writeFile(archive *tar.Writer, filePath string, name string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		return err
	}

	header.Name = name

	if err := archive.WriteHeader(header); err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(archive, file); err != nil {
		return err
	}

	return nil
}
