package data

import (
	"database/sql"
	"errors"

	"github.com/GuyOz5252/go-app/internal/core"
)

type SqlUserRepository struct {
	db *sql.DB
	queries *map[string]string
}

func NewSqlUserRepository(db *sql.DB, queries *map[string]string) core.UserRepository {
	return &SqlUserRepository{
		db: db,
		queries: queries,
	}
}

func (r *SqlUserRepository) GetById(id int) (*core.User, error) {
	user := &core.User{}
	query, ok := (*r.queries)["get_by_id"]
	if !ok {
		return nil, errors.New("user query 'get_by_id' not configured")
	}
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
	query, ok := (*r.queries)["create"]
	if !ok {
		return -1, errors.New("user query 'create' not configured")
	}
	if err := r.db.QueryRow(query, user.Username, user.Email).Scan(&userId); err != nil {
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
	query, ok := (*r.queries)["update"]
	if !ok {
		return errors.New("user query 'update' not configured")
	}
	err := r.db.QueryRow(query, user.Username, user.Email, user.Id).Scan(&userId)
	if err == sql.ErrNoRows {
		return core.ErrNotFound
	}

	return err
}

func (r *SqlUserRepository) Delete(id int) error {
	var userId int
	query, ok := (*r.queries)["delete"]
	if !ok {
		return errors.New("user query 'delete' not configured")
	}
	err := r.db.QueryRow(query, id).Scan(&userId)
	if err == sql.ErrNoRows {
		return core.ErrNotFound
	}

	return err
}

func (r *SqlUserRepository) ExistsByUsername(username string) (bool, error) {
	var exists bool
	query, ok := (*r.queries)["exists_by_username"]
	if !ok {
		return false, errors.New("user query 'exists_by_username' not configured")
	}
	err := r.db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *SqlUserRepository) ExistsByEmail(email string) (bool, error) {
	var exists bool
	query, ok := (*r.queries)["exists_by_email"]
	if !ok {
		return false, errors.New("user query 'exists_by_email' not configured")
	}
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
