package handlers

import (
	"net/http"
	"time"

	"forum/config"
	"forum/utils"
)

func deleteCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func LogutHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := utils.GetSessionCookie(r)
	config.SESSION.DeleteSession(sessionId)
	deleteCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
