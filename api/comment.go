package api

import (
	"log"
	"net/http"
	"strings"

	"forum/config"
	"forum/models"
	"forum/utils"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	sessionId := utils.GetSessionCookie(r)
	session := config.IsAuth(sessionId)
	if session == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusUnauthorized, "Invalid request method", nil)
		return
	}

	var comment models.Comment
	err := utils.ReadJSON(r, &comment)
	comment.UserID = session.UserId
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if strings.TrimSpace(comment.Comment) == "" {
		utils.WriteJSON(w, http.StatusBadRequest, "comment is required", nil)
		return
	}
	postRepo := models.NewPostRepository()
	commentRepo := models.NewCommentRepository()
	isExist, err := postRepo.IsPostExist(comment.PostID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isExist {
		utils.WriteJSON(w, http.StatusBadRequest, "Post does not exist", nil)
		return
	}
	err = commentRepo.Create(&comment)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	utils.WriteJSON(w, 200, "comment added succefully", comment)
}

// userId, postId, isLike
func HandleLikeComment(w http.ResponseWriter, r *http.Request) {
	sessionId := utils.GetSessionCookie(r)
	session := config.IsAuth(sessionId)
	if session == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var like models.CommentLike
	err := utils.ReadJSON(r, &like)
	like.UserID = session.UserId
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if like.IsLike != 1 && like.IsLike != -1 {
		utils.WriteJSON(w, http.StatusBadRequest, "invalid request", nil)
		return
	}
	comntRepo := models.NewCommentRepository()
	isExist, err := comntRepo.IsCommentExist(like.CommentId)
	if err != nil {
	}
	if !isExist {
		utils.WriteJSON(w, http.StatusBadRequest, "Comment does not exits", nil)
		return
	}
	isReactionExist, err := comntRepo.IsReactionExist(like.UserID, like.CommentId, like.IsLike)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if isReactionExist {
		err = comntRepo.DeleteReaction(like.UserID, like.CommentId)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	} else {
		err = comntRepo.ReactComment(like)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}
	commentReaction, err := comntRepo.GetCommentReaction(like.CommentId)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.WriteJSON(w, 200, "comment like updated succefully", commentReaction)
}
