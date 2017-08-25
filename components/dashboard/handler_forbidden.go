package dashboard

import (
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/session"
)

type ForbiddenHandler struct {
	Handler
}

func (h *ForbiddenHandler) IsAuth() bool {
	return false
}

func (h *ForbiddenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if username, err := session.GetString(r, "username"); err == nil && username != "" {
		h.Redirect("/", http.StatusFound, w, r)
		return
	}

	config := ConfigFromContext(r.Context())
	checkUsername := config.GetString(ConfigAuthUser)
	checkPassword := config.GetString(ConfigAuthPassword)

	if checkUsername == "" && checkPassword == "" {
		session.PutString(r, "username", "anonymous")
		h.Redirect("/", http.StatusFound, w, r)
		return
	}

	handleUrl := "/dashboard/login"

	if r.URL.Path != handleUrl {
		h.Redirect(handleUrl, http.StatusFound, w, r)
		return
	}

	var err error
	request := RequestFromContext(r.Context())

	if request.IsPost() {
		if err = r.ParseForm(); err == nil {
			username := r.PostForm.Get("username")
			password := r.PostForm.Get("password")

			if checkUsername == username && checkPassword == password {
				if err = session.RegenerateToken(r); err == nil {
					if err = session.PutString(r, "username", username); err == nil {
						h.Redirect("/", http.StatusFound, w, r)
						return
					}
				}
			} else {
				err = fmt.Errorf("Invalid username and/or password")
			}
		}
	}

	w.WriteHeader(http.StatusForbidden)
	h.RenderLayout(r.Context(), ComponentName, "login", "blank", map[string]interface{}{
		"login_url": handleUrl,
		"error":     err,
	})
}
