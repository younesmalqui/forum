package services

import (
	"strings"
	"time"

	"forum/config"
	"forum/models"
)

func CreateNewPost(post *models.Post) error {
	postRepo := models.NewPostRepository()
	TagsRepo := models.NewTagRepository()
	// check if input empty
	if strings.TrimSpace(post.Content) == "" || post.IsTagsEmpty() {
		return config.NewError(errFieldsEmpty)
	}
	post.CreatedAt = time.Now()
	err := postRepo.Create(post)
	if err != nil {
		return err
	}
	TagsRepo.LinkTagsToPost(post.ID, post.Tags)
	return err
}
