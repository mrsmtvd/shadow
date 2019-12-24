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
	Architecture386     = "386"
	ArchitectureAMD64   = "amd64"
	ArchitectureARM     = "arm"
	ArchitectureARM64   = "arm64"
	ArchitecturePPC     = "ppc"
	ArchitecturePPC64   = "ppc64"
	ArchitecturePPC64LE = "ppc64le"
	ArchitectureS390X   = "s390x"
)

type Release interface {
	Version() string
	File() (io.ReadCloser, error)
	FileBinary() (io.ReadCloser, error)
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

func ArchitectureFromReader(reader goArchReader) string {
	data := make([]byte, 16)
	if _, err := io.ReadFull(reader, data); err != nil {
		return ArchitectureUnknown
	}
	reader.Seek(0, io.SeekStart)

	if bytes.HasPrefix(data, []byte("\x7FELF")) {
		if _elf, err := elf.NewFile(reader); err == nil {
			switch _elf.Machine {
			case elf.EM_386:
				return Architecture386
			case elf.EM_X86_64:
				return ArchitectureAMD64
			case elf.EM_ARM:
				return ArchitectureARM
			case elf.EM_AARCH64:
				return ArchitectureARM64
			case elf.EM_PPC64:
				if _elf.ByteOrder == binary.LittleEndian {
					return ArchitecturePPC64LE
				}
				return ArchitecturePPC64
			case elf.EM_S390:
				return ArchitectureS390X
			}
		}
	}

	if bytes.HasPrefix(data, []byte("\xFE\xED\xFA")) || bytes.HasPrefix(data[1:], []byte("\xFA\xED\xFE")) {
		if _macho, err := macho.NewFile(reader); err == nil {
			switch _macho.Cpu {
			case macho.Cpu386:
				return Architecture386
			case macho.CpuAmd64:
				return ArchitectureAMD64
			case macho.CpuArm:
				return ArchitectureARM
			case macho.CpuArm64:
				return ArchitectureARM64
			case macho.CpuPpc:
				return ArchitecturePPC
			case macho.CpuPpc64:
				return ArchitecturePPC64
			}
		}
	}

	return ArchitectureUnknown
}

func GenerateReleaseID(rl Release) string {
	h := md5.New()

	h.Write(rl.Checksum())
	h.Write(separator)
	h.Write([]byte(rl.Path()))

	return hex.EncodeToString(h.Sum(nil))
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
