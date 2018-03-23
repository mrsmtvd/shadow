package i18n

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/elazarl/go-bindata-assetfs"
)

const (
	TranslateFileExt = ".mo"
	MessagesDirName  = "LC_MESSAGES"
)

func FromAssetFS(fs *assetfs.AssetFS) map[string]io.ReadSeeker {
	root, err := fs.Open("")
	if err != nil {
		return nil
	}

	dirs, err := root.Readdir(0)
	if err != nil {
		return nil
	}

	locales := make(map[string]io.ReadSeeker, len(dirs))

	for _, d := range dirs {
		localeDir, err := fs.Open(filepath.Join(d.Name(), MessagesDirName))
		if err != nil {
			continue
		}

		localeFiles, err := localeDir.Readdir(0)
		if err != nil {
			continue
		}

		for _, f := range localeFiles {
			if f.IsDir() || filepath.Ext(f.Name()) != TranslateFileExt {
				continue
			}

			content, err := fs.Asset(filepath.Join(fs.Prefix, d.Name(), MessagesDirName, f.Name()))
			if err != nil {
				continue
			}

			locales[d.Name()] = bytes.NewReader(content)
			break
		}
	}

	return locales
}
