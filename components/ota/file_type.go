package ota

import (
	"bytes"
	"io"
	"strings"
)

const (
	FileTypeUnknown FileType = "unknown"
	FileTypeBinary  FileType = "binary"
	FileTypeZip     FileType = "zip"
)

type FileType string

func (t FileType) String() string {
	switch t {
	case FileTypeBinary, FileTypeZip, FileTypeUnknown:
		return string(t)
	}

	return FileTypeUnknown.String()
}

func (t FileType) Ext() string {
	switch t {
	case FileTypeBinary:
		return ".bin"
	case FileTypeZip:
		return ".zip"
	}

	return ""
}

func (t FileType) MIME() string {
	switch t {
	case FileTypeBinary:
		return "application/x-binary"
	case FileTypeZip:
		return "application/zip"
	}

	return ""
}

var (
	separator     = []byte("|")
	fileTypeSigns = []struct {
		fileType   FileType
		magicBytes []byte
	}{
		{FileTypeBinary, []byte("\x7FELF")},
		{FileTypeBinary, []byte("\xFE\xED\xFA")},
		{FileTypeBinary, []byte("\xFA\xED\xFE")},
		{FileTypeBinary, []byte("\xCF\xFA\xED\xFE")},
		{FileTypeBinary, []byte("\xFE\xED\xFA\xCE")},
		{FileTypeBinary, []byte("\xFE\xED\xFA\xCF")},
		{FileTypeBinary, []byte("\xCE\xFA\xED\xFE")},
		{FileTypeZip, []byte("PK\x03\x04")},
	}
	fileTypeMIME = []struct {
		fileType FileType
		mime     string
	}{
		{FileTypeBinary, "application/x-binary"},
		{FileTypeBinary, "application/zip"},
	}
)

func FileTypeFromData(data io.Reader) FileType {
	buf := make([]byte, 8)

	if _, err := data.Read(buf); err != nil {
		return FileTypeUnknown
	}

	for _, sign := range fileTypeSigns {
		if bytes.Equal(sign.magicBytes, buf[:len(sign.magicBytes)]) {
			return sign.fileType
		}
	}

	return FileTypeUnknown
}

func FileTypeFromMIME(contentType string) FileType {
	contentType = strings.ToLower(contentType)

	for _, m := range fileTypeMIME {
		if strings.HasPrefix(m.mime, contentType) {
			return m.fileType
		}
	}

	return FileTypeUnknown
}
