package api

import (
	"net/http"

	"forum/config"
	"forum/models"
	"forum/services"
	"forum/utils"
)

func PostApi(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePost(w, r)
	default:
		utils.WriteJSON(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	sessionId := utils.GeTCookie("session", r)
	session, _ := config.SESSION.GetSession(sessionId)
	if session == nil {
		utils.WriteJSON(w, http.StatusBadRequest, "Unauthorized access", nil)
		return
	}
	var post models.Post
	err := utils.ReadJSON(r, &post)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if utils.IsEmpty(post.Content) || utils.IsEmpty(post.Title) {
		utils.WriteJSON(w, http.StatusBadRequest, "All fields required!", nil)
		return
	}
	post.UserID = session.UserId
	err = services.CreateNewPost(&post)
	if err != nil {
		if err.(*config.CustomError).IsInternal() {
			utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		utils.WriteJSON(w, http.StatusBadRequest, "err", nil)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, "Post created successfully", post)
}
