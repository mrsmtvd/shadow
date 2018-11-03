package handlers

import (
	"net/http"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

type variableView struct {
	Variable config.Variable
	Watchers []config.Watcher
}

type hasSource interface {
	Source() string
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
}

func (h *ManagerHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	var err error

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
					h.Logger().Info("User change config "+key,
						"user.id", user.UserID,
						"user.name", user.Name,
					)
				}
			}

			if err == nil {
				h.Redirect(r.URL().String(), http.StatusFound, w, r)
				return
			}
		}
	}

	variables := map[string][]variableView{}
	for _, v := range r.Config().Variables() {
		source := v.(hasSource).Source()

		cmp, ok := variables[source]
		if !ok {
			variables[source] = make([]variableView, 0)
			cmp = variables[source]
		}

		cmp = append(cmp, variableView{
			Variable: v,
			Watchers: r.Config().Watchers(v.Key()),
		})
		variables[source] = cmp
	}

	h.Render(r.Context(), "manager", map[string]interface{}{
		"variables": variables,
		"error":     err,
	})
}
