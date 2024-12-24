package utils

import "net/http"

func GeTCookie(name string, r *http.Request) string {
	session, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return session.Value
}

func GetSessionCookie(r *http.Request) string {
	return GeTCookie("session", r)
}
