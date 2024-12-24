package models

import (
	"database/sql"
	"errors"

	"forum/config"
)

var (
	ErrDB           = errors.New("database error")
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: config.DB}
}

func (r *UserRepository) CreateUser(user *User) error {
	query := "INSERT INTO users (email, username, password) VALUES (?,?,?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return config.NewInternalError(err)
	}
	defer stmt.Close()
	if err != nil {
		return config.NewInternalError(err)
	}
	result, err := stmt.Exec(user.Email, user.Username, user.Password)
	if err != nil {
		return config.NewInternalError(err)
	}
	user.ID, _ = result.LastInsertId()
	return nil
}

func (r *UserRepository) GetUserByID(id string) (*User, error) {
	query := "SELECT id, email, username, password FROM users WHERE id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, ErrDB
	}
	row := stmt.QueryRow(id)
	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	query := "SELECT id, email, username, password FROM users WHERE email = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, config.NewInternalError(err)
	}
	row := stmt.QueryRow(email)

	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, config.NewError(ErrUserNotFound)
		}
		return nil, config.NewInternalError(err)
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*User, error) {
	query := "SELECT id, email, username, password FROM users WHERE username = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, config.NewInternalError(err)
	}
	row := stmt.QueryRow(username)

	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, config.NewError(ErrUserNotFound)
		}
		return nil, config.NewInternalError(err)
	}
	return &user, nil
}

func (r *UserRepository) UserExists(username, email string) (bool, error) {
	var count int
	query := `
    SELECT COUNT(*) FROM users 
    WHERE username = ? OR email = ?
    `
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return false, config.NewInternalError(err)
	}
	err = stmt.QueryRow(username, email).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, config.NewError(ErrUserNotFound)
		}
		return false, config.NewInternalError(err)
	}
	return count > 0, nil
}
