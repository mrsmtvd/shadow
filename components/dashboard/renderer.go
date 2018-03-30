package dashboard

import (
	"context"

	"github.com/elazarl/go-bindata-assetfs"
)

type Renderer interface {
	Render(ctx context.Context, component, view string, data map[string]interface{}) error
	RenderLayout(ctx context.Context, component, view, layout string, data map[string]interface{}) error
}

type HasTemplates interface {
	DashboardTemplates() *assetfs.AssetFS
}

type HasTemplateFunctions interface {
	DashboardTemplateFunctions() map[string]interface{}
}
