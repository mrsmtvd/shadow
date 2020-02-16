package internal

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) initAssetFS() {
	for _, component := range c.components {
		if cmp, ok := component.(dashboard.HasAssetFS); ok {
			c.RegisterAssetFS(component.Name(), cmp.AssetFS())
		}
	}
}

func (c *Component) RegisterAssetFS(name string, fs *assetfs.AssetFS) {
	if fs.Prefix != "" {
		fs.Prefix = ""
	}

	c.registryAssetFS.Store(name, fs)
}
