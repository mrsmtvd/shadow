package dashboard

import (
	"context"

	"github.com/elazarl/go-bindata-assetfs"
)

type Renderer interface {
	Render(ctx context.Context, c, v string, d map[string]interface{}) error
	RenderLayout(ctx context.Context, c, v, l string, d map[string]interface{}) error
}

type HasTemplates interface {
	DashboardTemplates() *assetfs.AssetFS
}
