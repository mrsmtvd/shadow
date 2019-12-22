package release

import (
	"io"
	"net/http"
	"net/url"
)

type HTTPFile struct {
	u            *url.URL
	version      string
	checksum     []byte
	size         int64
	architecture string
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

func (f *HTTPFile) BinFile() (io.ReadCloser, error) {
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
