package dashboard

import (
	"fmt"
	"net/http"
)

type ForbiddenHandler struct {
	Handler
}

func (h *ForbiddenHandler) getRedirectURL(r *Request) string {
	redirectURL, err := r.Session().GetString(SessionLastURL)
	if err == nil && redirectURL != "" {
		return redirectURL
	}

	return "/"
}

func (h *ForbiddenHandler) ServeHTTP(w *Response, r *Request) {
	session := r.Session()

	if username, err := session.GetString(SessionUsername); err == nil && username != "" {
		h.Redirect(h.getRedirectURL(r), http.StatusFound, w, r)
		return
	}

	checkUsername := r.Config().GetString(ConfigAuthUser)
	checkPassword := r.Config().GetString(ConfigAuthPassword)

	if checkUsername == "" && checkPassword == "" {
		session.PutString(SessionUsername, "anonymous")
		h.Redirect(h.getRedirectURL(r), http.StatusFound, w, r)
		return
	}

	handleUrl := "/dashboard/login"

	if r.URL().Path != handleUrl {
		if !r.IsAjax() {
			session.PutString(SessionLastURL, r.URL().Path)
		}

		h.Redirect(handleUrl, http.StatusFound, w, r)
		return
	}

	var err error

	if r.IsPost() {
		if err = r.Original().ParseForm(); err == nil {
			username := r.Original().PostForm.Get("username")
			password := r.Original().PostForm.Get("password")

			if checkUsername == username && checkPassword == password {
				if err = session.RenewToken(); err == nil {
					if err = session.PutString(SessionUsername, username); err == nil {
						h.Redirect(h.getRedirectURL(r), http.StatusFound, w, r)
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
