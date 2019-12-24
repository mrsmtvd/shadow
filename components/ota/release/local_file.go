package release

import (
	"crypto/md5"
	"io"
	"io/ioutil"
	"os"

	"github.com/kihamo/shadow/components/ota"
)

type LocalFile struct {
	path         string
	version      string
	checksum     []byte
	architecture string
	fileInfo     os.FileInfo
	fileType     ota.FileType
}

func NewLocalFile(path, version string) (*LocalFile, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	return NewLocalFileFromFD(fd, version)
}

func NewLocalFileFromStream(stream io.Reader, version, dir string) (*LocalFile, error) {
	if dir == "" {
		dir = os.TempDir()
	}

	fd, err := ioutil.TempFile(dir, "release-")
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	_, err = io.Copy(fd, stream)
	if err != nil {
		return nil, err
	}

	return NewLocalFileFromFD(fd, version)
}

func NewLocalFileFromFD(fd *os.File, version string) (*LocalFile, error) {
	fd.Seek(0, 0)

	stat, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	h := md5.New()
	if _, err := io.Copy(h, fd); err != nil {
		return nil, err
	}

	fd.Seek(0, 0)

	fileType := ota.FileTypeFromData(fd)

	fd.Seek(0, 0)

	if version == "" {
		version = stat.Name()
	}

	return &LocalFile{
		path:         fd.Name(),
		version:      version,
		checksum:     h.Sum(nil),
		architecture: ota.ArchitectureFromReader(fd),
		fileInfo:     stat,
		fileType:     fileType,
	}, nil
}

func (f *LocalFile) Version() string {
	return f.version
}

func (f *LocalFile) File() (io.ReadCloser, error) {
	return os.Open(f.path)
}

func (f *LocalFile) Path() string {
	return f.path
}

func (f *LocalFile) Checksum() []byte {
	return f.checksum
}

func (f *LocalFile) Size() int64 {
	return f.fileInfo.Size()
}

func (f *LocalFile) Architecture() string {
	return f.architecture
}

func (f *LocalFile) FileInfo() os.FileInfo {
	return f.fileInfo
}

func (f *LocalFile) Type() ota.FileType {
	return f.fileType
}
