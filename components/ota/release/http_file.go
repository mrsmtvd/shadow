package release

import (
	"io"
	"net/http"
	"net/url"
	"sync"

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
}

func NewHTTPFile(path, version string, checksum []byte, size int64, architecture string) (*HTTPFile, error) {
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
	}, nil
}

func (f *HTTPFile) Version() string {
	return f.version
}

func (f *HTTPFile) File() (io.ReadCloser, error) {
	response, err := http.Get(f.u.String())
	if err != nil {
		return nil, err
	}

	return response.Body, err
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
