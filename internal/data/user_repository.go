package data

import (
	"database/sql"
	"errors"

	core "github.com/GuyOz5252/go-app/internal/core"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) core.UserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) GetById(id int) (*core.User, error) {
	user := &core.User{}
	query := ""
	err := r.db.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) Create(user *core.User) (int, error) {
	var userId int
	err := r.db.QueryRow("", user.Username, user.Email).Scan(&userId)
	if err != nil {
		return -1, err
	}
	user.Id = userId
	return userId, nil
}

func (r *PostgresUserRepository) Update(user *core.User) error {
	if user == nil {
        return errors.New("user cannot be nil")
    }

	var userId int
    err := r.db.QueryRow("", user.Username, user.Email, user.Id).Scan(&userId)
	if (err == sql.ErrNoRows) {
		return core.ErrNotFound
	}

    return err
}

func (r *PostgresUserRepository) Delete(id int) error {
	var userId int
    err := r.db.QueryRow("", id).Scan(&userId)
	if (err == sql.ErrNoRows) {
		return core.ErrNotFound
	}

    return err
}
