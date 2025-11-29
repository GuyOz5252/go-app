package data

import (
	"database/sql"
	"errors"

	"github.com/GuyOz5252/go-app/internal/core"
)

type SqlUserRepository struct {
	db *sql.DB
}

func NewSqlUserRepository(db *sql.DB) core.UserRepository {
	return &SqlUserRepository{
		db: db,
	}
}

func (r *SqlUserRepository) GetById(id int) (*core.User, error) {
	user := &core.User{}
	query := "SELECT id, username, email FROM users WHERE id = $1"
	err := r.db.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *SqlUserRepository) Create(user *core.User) (int, error) {
	var userId int
	err := r.db.QueryRow("INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id", user.Username, user.Email).Scan(&userId)
	if err != nil {
		return -1, err
	}
	user.Id = userId
	return userId, nil
}

func (r *SqlUserRepository) Update(user *core.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	var userId int
	err := r.db.QueryRow("", user.Username, user.Email, user.Id).Scan(&userId)
	if err == sql.ErrNoRows {
		return core.ErrNotFound
	}

	return err
}

func (r *SqlUserRepository) Delete(id int) error {
	var userId int
	err := r.db.QueryRow("", id).Scan(&userId)
	if err == sql.ErrNoRows {
		return core.ErrNotFound
	}

	return err
}
