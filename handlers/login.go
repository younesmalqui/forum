package handlers

import (
	"net/http"

	"forum/config"
	"forum/services"
	"forum/utils"
)

type pageData struct {
	Error  string
	Method string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	utils.RedirectIsAuth(w, r)
	switch r.Method {
	case http.MethodGet:
		getLogin(w)
	case http.MethodPost:
		postLogin(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func getLogin(w http.ResponseWriter) {
	page := pageData{
		Method: "GET",
	}
	config.TMPL.Render(w, "login.html", page)
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	page := pageData{
		Method: "POST",
	}
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := services.LoginUser(username, password)
	if err != nil {
		// TODO make page
		page.Error = err.Error()
		config.TMPL.Render(w, "login.html", page)
		return
	}

	session, err := config.SESSION.CreateSession(user.Username, user.ID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	cookies := http.Cookie{
		Name:    "session",
		Value:   session.ID,
		Expires: session.ExpiresAt,
		Path:    "/",
	}
	http.SetCookie(w, &cookies)
	config.TMPL.Render(w, "login.html", page)
}
