package ota

import (
	"bytes"
	"crypto/md5"
	"debug/elf"
	"debug/macho"
	"encoding/binary"
	"encoding/hex"
	"io"
	"path/filepath"
	"strings"
)

const (
	ArchitectureUnknown = "unknown"
)

type Release interface {
	Version() string
	File() (io.ReadCloser, error)
	Path() string
	Type() FileType
	Checksum() []byte
	Size() int64
	Architecture() string
}

type Repository interface {
	Releases(arch string) ([]Release, error)
}

type RepositoryRemover interface {
	Remove(Release) error
}

type goArchReader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

func GoArch(reader goArchReader) string {
	data := make([]byte, 16)
	if _, err := io.ReadFull(reader, data); err != nil {
		return ArchitectureUnknown
	}
	reader.Seek(0, 0)

	if bytes.HasPrefix(data, []byte("\x7FELF")) {
		if _elf, err := elf.NewFile(reader); err == nil {
			switch _elf.Machine {
			case elf.EM_386:
				return "386"
			case elf.EM_X86_64:
				return "amd64"
			case elf.EM_ARM:
				return "arm"
			case elf.EM_AARCH64:
				return "arm64"
			case elf.EM_PPC64:
				if _elf.ByteOrder == binary.LittleEndian {
					return "ppc64le"
				}
				return "ppc64"
			case elf.EM_S390:
				return "s390x"
			}
		}
	}

	if bytes.HasPrefix(data, []byte("\xFE\xED\xFA")) || bytes.HasPrefix(data[1:], []byte("\xFA\xED\xFE")) {
		if _macho, err := macho.NewFile(reader); err == nil {
			switch _macho.Cpu {
			case macho.Cpu386:
				return "386"
			case macho.CpuAmd64:
				return "amd64"
			case macho.CpuArm:
				return "arm"
			case macho.CpuArm64:
				return "arm64"
			case macho.CpuPpc:
				return "ppc"
			case macho.CpuPpc64:
				return "ppc64"
			}
		}
	}

	return ArchitectureUnknown
}

func GenerateReleaseID(rl Release) string {
	hasher := md5.New()

	hasher.Write(rl.Checksum())
	hasher.Write(separator)
	hasher.Write([]byte(rl.Path()))

	return hex.EncodeToString(hasher.Sum(nil))
}

func GenerateFileName(rl Release) string {
	basePath := filepath.Base(rl.Path())
	ext := rl.Type().Ext()

	if ext != "" && strings.HasSuffix(basePath, ext) {
		return basePath
	}

	return strings.ReplaceAll(basePath, " ", "_") +
		"." + strings.ReplaceAll(rl.Version(), " ", ".") +
		"." + rl.Architecture() + ext
}
