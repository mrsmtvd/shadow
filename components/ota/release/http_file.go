package release

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/kihamo/shadow/components/ota"
)

type HTTPFile struct {
	mutex sync.RWMutex

	u            *url.URL
	version      string
	checksum     []byte
	size         int64
	architecture string
	fileType     *ota.FileType
	fileCache    string
	createdAt    *time.Time
}

func NewHTTPFile(path, version string, checksum []byte, size int64, architecture string, createdAt *time.Time) (*HTTPFile, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	return &HTTPFile{
		u:            u,
		version:      version,
		checksum:     checksum,
		size:         size,
		architecture: architecture,
		createdAt:    createdAt,
	}, nil
}

func (f *HTTPFile) Version() string {
	return f.version
}

func (f *HTTPFile) File() (io.ReadCloser, error) {
	// проверяем локальный кэш
	f.mutex.RLock()
	fileCache := f.fileCache
	f.mutex.RUnlock()

	if fileCache != "" {
		if fd, err := os.Open(fileCache); err == nil {
			return fd, nil
		}
	}

	response, err := http.Get(f.u.String())
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("remote server response not 200 OK")
	}

	// сохраняем в локальный кэш
	fd, err := ioutil.TempFile(os.TempDir(), "release-download-")
	if err == nil {
		_, err = io.Copy(fd, response.Body)
		if err == nil {
			f.mutex.Lock()
			f.fileCache = fd.Name()
			f.mutex.Unlock()

			return os.Open(fd.Name())
		}
	}

	return nil, err
}

func (f *HTTPFile) FileBinary() (io.ReadCloser, error) {
	return f.File()
}

func (f *HTTPFile) Path() string {
	return f.u.String()
}

func (f *HTTPFile) Checksum() []byte {
	return f.checksum
}

func (f *HTTPFile) Size() int64 {
	return f.size
}

func (f *HTTPFile) Architecture() string {
	return f.architecture
}

func (f *HTTPFile) Type() ota.FileType {
	f.mutex.RLock()
	ft := f.fileType
	f.mutex.RUnlock()

	if ft == nil {
		found := f.getFileType()
		ft = &found

		f.mutex.Lock()
		f.fileType = &found
		f.mutex.Unlock()
	}

	return *ft
}

func (f *HTTPFile) CreatedAt() *time.Time {
	return f.createdAt
}

func (f *HTTPFile) getFileType() ota.FileType {
	// попытка вычитать HEAD
	response, err := http.Head(f.u.String())
	if err == nil {
		if response.StatusCode == http.StatusOK {
			return ota.FileTypeFromMIME(response.Header.Get("Content-Type"))
		}
	}

	// если не удалось, то делает GET
	response, err = http.Get(f.u.String())
	if err == nil {
		if response.StatusCode == http.StatusOK {
			return ota.FileTypeFromMIME(response.Header.Get("Content-Type"))
		}
	}

	return ota.FileTypeUnknown
}
