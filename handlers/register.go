package handlers

import (
	"net/http"

	"forum/config"
	"forum/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	utils.RedirectIsAuth(w, r)

	config.TMPL.Render(w, "register.html", nil)
}
