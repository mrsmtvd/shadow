package dashboard

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	Renderer() Renderer
	RegisterAssetFS(name string, fs *assetfs.AssetFS)
}
