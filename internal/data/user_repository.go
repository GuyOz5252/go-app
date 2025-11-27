package data

import (
	"database/sql"

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
			return nil, core.NotFoundError
		}
		return nil, err
	}
	return user, nil
}

func (userRepository *PostgresUserRepository) Create(user *core.User) error {
	return nil
}

func (userRepository *PostgresUserRepository) Update(user *core.User) error {
	return nil
}

func (userRepository *PostgresUserRepository) Delete(id int) error {
	return nil
}
