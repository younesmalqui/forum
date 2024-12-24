package models

import (
	"database/sql"

	"forum/config"
)

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TagRepository struct {
	db *sql.DB
}

func NewTagRepository() *TagRepository {
	return &TagRepository{db: config.DB}
}

func (r *TagRepository) CreateTag(name string) (*Tag, error) {
	stmt, err := r.db.Prepare("INSERT INTO tags (name) VALUES (?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Tag{ID: id, Name: name}, nil
}

func (r *TagRepository) GetAllTags() ([]Tag, error) {
	stmt, err := r.db.Prepare("SELECT id, name FROM tags")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *TagRepository) IsTagExists(name string) (bool, error) {
	var exists bool

	stmt, err := r.db.Prepare("SELECT COUNT(*) > 0 FROM tags WHERE name = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(name).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists, nil
}

func (r *TagRepository) GetTagsForPost(postId int64) ([]string, error) {
	stmt, err := r.db.Prepare(`
        SELECT t.name
        FROM tags t
        JOIN post_tags pt ON t.id = pt.tagId
        WHERE pt.postId = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tagName string
		if err := rows.Scan(&tagName); err != nil {
			return nil, err
		}
		tags = append(tags, tagName)
	}

	return tags, nil
}

func (r *TagRepository) LinkTagsToPost(postId int64, tagNames []string) error {
	selectStmt, err := r.db.Prepare("SELECT id FROM tags WHERE name = ?")
	if err != nil {
		return err
	}
	defer selectStmt.Close()

	insertStmt, err := r.db.Prepare("INSERT INTO post_tags (postId, tagId) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer insertStmt.Close()

	for _, tagName := range tagNames {
		var tagId int64

		err := selectStmt.QueryRow(tagName).Scan(&tagId)
		if err != nil {
			if err == sql.ErrNoRows {
				newTag, err := r.CreateTag(tagName)
				if err != nil {
					return err
				}
				tagId = newTag.ID
			} else {
				return err
			}
		}

		_, err = insertStmt.Exec(postId, tagId)
		if err != nil {
			return err
		}
	}

	return nil
}
