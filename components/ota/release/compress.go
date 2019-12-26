package release

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/kihamo/shadow/components/ota"
)

// подходит первый попавшийся файл
var defaultFnSearchFile = func(s string) bool {
	return true
}

type Compress struct {
	original     ota.Release
	fnSearchFile func(string) bool
	architecture string
}

func NewCompress(release ota.Release) *Compress {
	return NewCompressWithFn(release, defaultFnSearchFile)
}

func NewCompressWithFn(release ota.Release, fnSearchFile func(string) bool) *Compress {
	rl := &Compress{
		original:     release,
		fnSearchFile: fnSearchFile,
		architecture: release.Architecture(),
	}

	// если архитектура не определена, то пытаемся найти бинарник в архиве и определить через него
	if rl.architecture == ota.ArchitectureUnknown {
		if reader, err := rl.FileBinary(); err == nil {
			fd, _, err := createTempFileFromReader(reader)
			if err == nil {
				rl.architecture = ota.ArchitectureFromReader(fd)
				fd.Close()
			}

			reader.Close()
		}
	}

	return rl
}

func (f *Compress) Version() string {
	return f.original.Version()
}

func (f *Compress) File() (io.ReadCloser, error) {
	return f.original.File()
}

func (f *Compress) FileBinary() (io.ReadCloser, error) {
	reader, err := f.original.FileBinary()
	if err != nil {
		return nil, err
	}

	if f.Type() != ota.FileTypeZip {
		return reader, nil
	}

	defer reader.Close()

	// в остальных случая нам надо:
	// 1. создать временный файл, записать туда содержимое
	fd, size, err := createTempFileFromReader(reader)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	// 3. открыть его через архиватор
	readerArchive, err := zip.NewReader(fd, size)
	if err != nil {
		return nil, err
	}

	// 4. найти в архиве нужный файл
	for _, fileInArchive := range readerArchive.File {
		if f.fnSearchFile(fileInArchive.Name) {
			// копируем бинарник во временный файл
			readerBinary, err := fileInArchive.Open()
			if err != nil {
				return nil, err
			}

			// TODO: непонятно как удалять этот файл

			fdBinary, _, err := createTempFileFromReader(readerBinary)
			readerBinary.Close()

			return fdBinary, err
		}
	}

	return nil, errors.New("binary file not found in archive")
}

func (f *Compress) Path() string {
	return f.original.Path()
}

func (f *Compress) Checksum() []byte {
	return f.original.Checksum()
}

func (f *Compress) Size() int64 {
	return f.original.Size()
}

func (f *Compress) Architecture() string {
	return f.architecture
}

func (f *Compress) Type() ota.FileType {
	return f.original.Type()
}

func (f *Compress) CreatedAt() *time.Time {
	return f.original.CreatedAt()
}

func (f *Compress) Validate() error {
	return f.original.Validate()
}

func createTempFileFromReader(reader io.Reader) (_ *fileTemp, size int64, err error) {
	// если пришел дескриптор то ничего создавать не нужно, файл уже есть, переиспользуем его
	if exist, ok := reader.(*os.File); ok {
		if info, err := exist.Stat(); err == nil {
			return &fileTemp{exist, false}, info.Size(), err
		}
	}

	// в противном случае создаем временный файл и копируем туда содержимое ридера
	var fd *os.File

	fd, err = ioutil.TempFile(os.TempDir(), "release-temp-")
	if err == nil {
		size, err = io.Copy(fd, reader)
		if err == nil {
			_, err = fd.Seek(0, io.SeekStart)
		}
	}

	return &fileTemp{fd, true}, size, err
}

type fileTemp struct {
	*os.File
	autoRemove bool
}

func (f *fileTemp) Close() (err error) {
	err = f.File.Close()
	if err == nil && f.autoRemove {
		err = os.Remove(f.Name())
	}

	return err
}
