package internal

import (
	"io"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/i18n"
)

func (c *Component) I18n() map[string]io.ReadSeeker {
	return i18n.FromAssetFS(&assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "locales",
	})
}
