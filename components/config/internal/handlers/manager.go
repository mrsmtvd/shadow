package handlers

import (
	"net/http"
	"strings"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

const (
	defaultComponentName = "main"
)

type variableView struct {
	Variable config.Variable
	Watchers []config.Watcher
}

func (v variableView) HasView(n string) bool {
	if len(v.Variable.View()) == 0 {
		return false
	}

	for _, nv := range v.Variable.View() {
		if nv == n {
			return true
		}
	}

	return false
}

func (v variableView) GetViewOption(o string) interface{} {
	if len(v.Variable.ViewOptions()) > 0 {
		if opt, ok := v.Variable.ViewOptions()[o]; ok {
			return opt
		}
	}

	return nil
}

type ManagerHandler struct {
	dashboard.Handler

	Application shadow.Application
	Component   config.Component
}

func (h *ManagerHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	var err error

	vars := h.Component.Variables()

	if r.IsPost() {
		err = r.Original().ParseForm()
		if err == nil {
			for key, values := range r.Original().PostForm {
				if !r.Config().Has(key) || !r.Config().IsEditable(key) || len(values) == 0 {
					continue
				}

				err = r.Config().Set(key, values[0])
				if err != nil {
					break
				}

				user := r.User()
				if user != nil {
					r.Logger().Infof("User change config %s", key, map[string]interface{}{
						"user.id":   user.UserID,
						"user.name": user.Name,
					})
				}
			}

			if err == nil {
				h.Redirect(r.URL().String(), http.StatusFound, w, r)
				return
			}
		}
	}

	variables := map[string]map[string]variableView{}
	for k, v := range vars {
		parts := strings.SplitN(k, ".", 2)

		cmpName := parts[0]
		if !h.Application.HasComponent(cmpName) {
			cmpName = defaultComponentName
		}

		cmp, ok := variables[cmpName]
		if !ok {
			variables[cmpName] = map[string]variableView{}
			cmp = variables[cmpName]
		}

		cmp[k] = variableView{
			Variable: v,
			Watchers: h.Component.Watchers(k),
		}
		variables[cmpName] = cmp
	}

	h.Render(r.Context(), h.Component.GetName(), "manager", map[string]interface{}{
		"variables": variables,
		"error":     err,
	})
}
