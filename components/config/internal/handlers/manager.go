package handlers

import (
	"net/http"

	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/logging"
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

	component config.Component
}

func NewManagerHandler(component config.Component) *ManagerHandler {
	return &ManagerHandler{
		component: component,
	}
}

func (h *ManagerHandler) ServeHTTP(w http.ResponseWriter, r *dashboard.Request) {
	var err error

	if r.IsPost() {
		err = r.Original().ParseForm()
		if err == nil {
			for key, values := range r.Original().PostForm {
				if !h.component.Has(key) || !h.component.IsEditable(key) || len(values) == 0 {
					continue
				}

				err = h.component.Set(key, values[0])
				if err != nil {
					break
				}

				user := r.User()
				if user != nil {
					logging.Log(r.Context()).Info("User change config "+key,
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

	for _, v := range h.component.Variables() {
		source := v.(hasSource).Source()

		cmp, ok := variables[source]
		if !ok {
			variables[source] = make([]variableView, 0)
			cmp = variables[source]
		}

		cmp = append(cmp, variableView{
			Variable: v,
			Watchers: h.component.Watchers(v.Key()),
		})
		variables[source] = cmp
	}

	if err != nil {
		r.Session().FlashBag().Error(err.Error())
	}

	h.Render(r.Context(), "manager", map[string]interface{}{
		"variables": variables,
	})
}
