package dashboard

import (
	"fmt"
	"net/http"
)

type ForbiddenHandler struct {
	Handler
}

func (h *ForbiddenHandler) ServeHTTP(w *Response, r *Request) {
	if username, err := r.Session().GetString(SessionUsername); err == nil && username != "" {
		h.Redirect("/", http.StatusFound, w, r)
		return
	}

	checkUsername := r.Config().GetString(ConfigAuthUser)
	checkPassword := r.Config().GetString(ConfigAuthPassword)

	if checkUsername == "" && checkPassword == "" {
		r.Session().PutString(SessionUsername, "anonymous")
		h.Redirect("/", http.StatusFound, w, r)
		return
	}

	handleUrl := "/dashboard/login"

	if r.URL().Path != handleUrl {
		h.Redirect(handleUrl, http.StatusFound, w, r)
		return
	}

	var err error

	if r.IsPost() {
		if err = r.Original().ParseForm(); err == nil {
			username := r.Original().PostForm.Get("username")
			password := r.Original().PostForm.Get("password")

			if checkUsername == username && checkPassword == password {
				if err = r.Session().RenewToken(); err == nil {
					if err = r.Session().PutString(SessionUsername, username); err == nil {
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
