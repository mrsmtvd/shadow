package dashboard

import (
	"github.com/elazarl/go-bindata-assetfs"
)

type HasAssetFS interface {
	AssetFS() *assetfs.AssetFS
}
