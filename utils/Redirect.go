package utils

import (
	"net/http"

	"forum/config"
)

func Redirect(w http.ResponseWriter, r *http.Request, newPath string) {
	http.Redirect(w, r, newPath, http.StatusFound)
}

func RedirectIsAuth(w http.ResponseWriter, r *http.Request) {
	sessionId, err := r.Cookie("session")
	if err == nil && config.IsAuth(sessionId.Value) != nil {
		Redirect(w, r, "/")
	}
}
