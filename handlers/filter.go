package handlers

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"forum/config"
	c "forum/config"
	"forum/models"
	"forum/utils"
)

const (
	ALL = iota
	MY_POST
	LIKED_POST
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		PostFilter(w, r)
	default:
	}
}

func pagination(r *http.Request) (int, int) {
	pageStr := r.URL.Query().Get("page")
	currPage, err := strconv.Atoi(pageStr)
	if err != nil || currPage < 1 {
		currPage = 1
	}
	return currPage, c.LIMIT_PER_PAGE
}

func PostFilter(w http.ResponseWriter, r *http.Request) {
	currPage, limit := pagination(r)
	sessionID := utils.GetSessionCookie(r)
	session, err := c.SESSION.GetSession(sessionID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var userId int64 = -1
	if c.IsAuth(sessionID) != nil {
		userId = session.UserId
	}
	r.ParseForm()
	postType := 0
	query := r.FormValue("query")
	if r.FormValue("options") != "" {
		postType, err = strconv.Atoi(r.FormValue("options"))
	}
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	postType, err = selectPostType(postType, userId != -1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	postRep := models.NewPostRepository()

	posts, err := postRep.GetPostsBy(query, postType, userId, currPage, limit)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	posts, err = getPostsFilter(posts)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	count := len(posts)
	if (currPage-1)*limit > count {
		currPage = max(int(math.Ceil(float64(count)/config.LIMIT_PER_PAGE)), 1)
	}
	sliceOfPosts := posts[(currPage-1)*limit : min(count, (currPage-1)*limit+limit)]
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page := NewPageStruct("forum", sessionID, nil)
	page.Data = IndexStruct{
		Posts:       sliceOfPosts,
		TotalPages:  int(math.Ceil(float64(count) / config.LIMIT_PER_PAGE)),
		CurrentPage: currPage,
		Query:       query,
		Option:      postType,
	}
	config.TMPL.Render(w, "filter.html", page)
}

func getPostsFilter(posts []*models.Post) ([]*models.Post, error) {
	tagsRepo := models.NewTagRepository()

	for _, post := range posts {
		tags, err := tagsRepo.GetTagsForPost(post.ID)
		if err != nil {
			return nil, err
		}
		post.Tags = tags
	}
	return posts, nil
}

func selectPostType(value int, isAuth bool) (int, error) {
	if value < 0 || value > 2 {
		return 0, errors.New("invalid option")
	}
	if isAuth {
		return value, nil
	}
	if value != 0 {
		return 0, errors.New("bad request, you can't filter with that option")
	}
	return 0, nil
}
