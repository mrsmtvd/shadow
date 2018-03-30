package dashboard

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow"
)

const (
	AssetFSPrefixRoute     = "assets"
	AssetFSPrefixTemplates = "templates"
)

type HasAssetFS interface {
	AssetFS() *assetfs.AssetFS
}

func RouteFromAssetFS(component HasAssetFS) Route {
	fs := component.AssetFS()
	fs.Prefix = AssetFSPrefixRoute

	return NewRoute("/"+component.(shadow.Component).Name()+"/assets/*filepath", fs).
		WithMethods([]string{http.MethodGet})
}

func TemplatesFromAssetFS(component HasAssetFS) *assetfs.AssetFS {
	fs := component.AssetFS()
	fs.Prefix = AssetFSPrefixTemplates

	return fs
}
