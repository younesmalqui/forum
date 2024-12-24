package database

import (
	"forum/config"
)

func execQuery(query string) error {
	db := config.DB
	_, err := db.Exec(query)
	if err != nil {
		return config.NewInternalError(err)
	}
	return nil
}

var Tables = []func() error{
	createPostTable,
	createUserTable,
	createSessionTable,
	createPostLikeTable,
	createCommentTable,
	createCommentLikeTable,
	createTagsTable,
	createPostTagsTable,
}

func createUserTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	`
	return execQuery(query)
}

func createSessionTable() error {
	query := `CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		username TEXT,
		userId TEXT,
		expiresAt DATETIME
	)`

	return execQuery(query)
}

func createPostTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
        userId INTEGER NOT NULL,
        content TEXT NOT NULL,
        createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	return execQuery(query)
}

func createPostLikeTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS post_reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userId TEXT NOT NULL,
    postId INTEGER NOT NULL,
    isLike INT NOT NULL,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(userId, postId)
	);`

	return execQuery(query)
}

func createCommentTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    postId INTEGER,
    userId INTEGER,
    comment TEXT NOT NULL,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (postId) REFERENCES posts(id) ON DELETE CASCADE
)`
	return execQuery(query)
}

func createCommentLikeTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS comment_reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userId INTEGER NOT NULL,
    commentId INTEGER NOT NULL,
    isLike INT NOT NULL,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(userId, commentId)
	);`
	return execQuery(query)
}

func createTagsTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE 
);`
	return execQuery(query)
}

func createPostTagsTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS post_tags (
    postId INTEGER,
    tagId INTEGER,
    PRIMARY KEY (postId, tagId),
    FOREIGN KEY (postId) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (tagId) REFERENCES tags(id) ON DELETE CASCADE
);`
	return execQuery(query)
}
