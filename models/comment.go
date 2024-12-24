package models

import (
	"database/sql"
	"time"

	"forum/config"
)

type CommentReaction struct {
	CommnetId int64 `json:"commentId"`
	Likes     int64 `json:"likes"`
	Dislikes  int64 `json:"dislikes"`
}

type CommentLike struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"userId"`
	CommentId int64 `json:"commentId"`
	IsLike    int   `json:"isLike"`
}

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"postId"`
	UserID    int64     `json:"userId"`
	Username  string    `json:"username"`
	Likes     int       `json:"likes"`
	DisLikes  int       `json:"disLikes"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{db: config.DB}
}

func (r *CommentRepository) Create(comment *Comment) error {
	query := `INSERT INTO comments (postId, userId, comment) VALUES (?, ?, ?)`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(comment.PostID, comment.UserID, comment.Comment)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	comment.ID = id
	return nil
}

func (r *CommentRepository) IsCommentExist(commentId int64) (bool, error) {
	query := "SELECT COUNT(id) FROM comments WHERE ID = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	var count int64
	err = stmt.QueryRow(commentId).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}

func (r *CommentRepository) GetPostComments(postID int64) ([]Comment, error) {
	query := `SELECT c.id ,c.postId ,c.userId, c.comment ,c.createdAt ,u.username , 
	(SELECT count(*) from comment_reactions WHERE isLike=1 AND commentId=c.id ) likes,
	(SELECT count(*) from comment_reactions WHERE isLike=-1 AND commentId=c.id ) dislike
	FROM comments c 
	LEFT JOIN comment_reactions l ON c.id = l.commentId 
	LEFT JOIN users u ON c.userId = u.id WHERE c.postId = ? GROUP BY c.id HAVING count(c.id) > 0 ORDER BY c.createdAt desc`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var comments []Comment
	rows, err := stmt.Query(postID)
	if err != nil {
		return comments, nil
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Comment, &comment.CreatedAt,
			&comment.Username, &comment.Likes, &comment.DisLikes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepository) DeleteReaction(userId int64, commentId int64) error {
	query := "DELETE FROM comment_reactions WHERE userId = ? AND commentId = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userId, commentId)
	if err != nil {
		return err
	}
	return nil
}

func (r *CommentRepository) IsReactionExist(userId int64, commmentId int64, isLike int) (bool, error) {
	query := "SELECT COUNT(*) FROM comment_reactions WHERE userId = ? AND commentId = ? AND isLike = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return false, err
	}
	var count int
	err = stmt.QueryRow(userId, commmentId, isLike).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}

func (r *CommentRepository) ReactComment(like CommentLike) error {
	stmt, err := r.db.Prepare(`
        INSERT INTO comment_reactions (userId, commentId, isLike)
        VALUES (?, ?, ?)
        ON CONFLICT(userId, commentId) DO UPDATE SET isLike = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(like.UserID, like.CommentId, like.IsLike, like.IsLike)
	return err
}

func (r *CommentRepository) GetCommentReaction(commentId int64) (*CommentReaction, error) {
	query := `SELECT
	SUM(CASE WHEN isLike = 1 THEN 1 ELSE 0 END) as likes,
	SUM(CASE WHEN islike = -1 THEN 1 ELSE 0 END) as dislike 
	FROM comment_reactions WHERE commentId = ? GROUP BY commentId`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	var commentReaction CommentReaction
	commentReaction.CommnetId = commentId
	err = stmt.QueryRow(&commentId).Scan(&commentReaction.Likes, &commentReaction.Dislikes)
	if err != nil {
		if err == sql.ErrNoRows {
			return &commentReaction, nil
		}
		return nil, err
	}
	return &commentReaction, nil
}
