package dashboard

import (
	"context"
	"io"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

type Renderer interface {
	IsRegisterNamespace(ns string) bool
	RegisterNamespace(ns string, fs *assetfs.AssetFS) error
	Render(wr io.Writer, ctx context.Context, ns, view string, data map[string]interface{}) error
	RenderAndReturn(ctx context.Context, ns, view string, data map[string]interface{}) (string, error)
	RenderLayout(wr io.Writer, ctx context.Context, ns, view, layout string, data map[string]interface{}) error
	RenderLayoutAndReturn(ctx context.Context, ns, view, layout string, data map[string]interface{}) (string, error)
}

type HasTemplates interface {
	DashboardTemplates() *assetfs.AssetFS
}

type HasTemplateFunctions interface {
	DashboardTemplateFunctions() map[string]interface{}
}

type HasToolbar interface {
	DashboardToolbar(ctx context.Context) string
}
