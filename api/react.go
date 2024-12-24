package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"forum/config"
	"forum/models"
	"forum/utils"
)

func ReactToPostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleReactPost(w, r)
	case http.MethodGet:
		handleReactGet(w, r)
	default:
		utils.WriteJSON(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}

func handleReactPost(w http.ResponseWriter, r *http.Request) {
	session := config.IsAuth(utils.GetSessionCookie(r))
	if session == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var like models.Like
	err := utils.ReadJSON(r, &like)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	like.UserID = session.UserId
	postRepo := models.NewPostRepository()
	isExistPost, _ := postRepo.IsPostExist(like.PostID)
	if !isExistPost {
		utils.WriteJSON(w, http.StatusBadRequest, "Invalid Request", nil)
		return
	}
	if like.IsLike != -1 && like.IsLike != 1 {
		utils.WriteJSON(w, http.StatusBadRequest, "Invalid Request", nil)
		return
	}
	likeRepo := models.NewLikeRepository()
	ok, err := likeRepo.IsUserReactToPost(like.UserID, like.PostID, like.IsLike)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if ok {
		if err := likeRepo.DeleteLike(like.UserID, like.PostID); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	} else {
		if err := likeRepo.AddReaction(&like); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reaction added successfully"})
}

func handleReactGet(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.ParseInt(r.URL.Query().Get("postId"), 10, 64)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, "Bad Request", nil)
		return
	}

	likeRepo := models.NewLikeRepository()

	count, err := likeRepo.GetPostLikes(postId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusNotFound, err.Error(), nil)
		} else {
			utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, "", count)
}
