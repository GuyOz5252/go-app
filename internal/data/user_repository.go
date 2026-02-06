package data

import (
	"context"
	"database/sql"
	"errors"

	"github.com/GuyOz5252/go-app/internal/core"
)

type SqlUserRepository struct {
	db      *sql.DB
	queries *map[string]string
}

func NewSqlUserRepository(db *sql.DB, queries *map[string]string) core.UserRepository {
	return &SqlUserRepository{
		db:      db,
		queries: queries,
	}
}

func (r *SqlUserRepository) GetById(ctx context.Context, id int) (*core.User, error) {
	user := &core.User{}
	query, ok := (*r.queries)["get_by_id"]
	if !ok {
		return nil, core.ErrQueryNotConfigured
	}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *SqlUserRepository) Create(ctx context.Context, user *core.User) (int, error) {
	var userId int
	query, ok := (*r.queries)["create"]
	if !ok {
		return -1, core.ErrQueryNotConfigured
	}
	if err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.PasswordHash).Scan(&userId); err != nil {
		return -1, err
	}
	user.Id = userId
	return userId, nil
}

func (r *SqlUserRepository) Update(ctx context.Context, user *core.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	var userId int
	query, ok := (*r.queries)["update"]
	if !ok {
		return core.ErrQueryNotConfigured
	}
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Id).Scan(&userId)
	if err == sql.ErrNoRows {
		return core.ErrNotFound
	}

	return err
}

func (r *SqlUserRepository) Delete(ctx context.Context, id int) error {
	var userId int
	query, ok := (*r.queries)["delete"]
	if !ok {
		return core.ErrQueryNotConfigured
	}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&userId)
	if err == sql.ErrNoRows {
		return core.ErrNotFound
	}

	return err
}

func (r *SqlUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	query, ok := (*r.queries)["exists_by_username"]
	if !ok {
		return false, core.ErrQueryNotConfigured
	}
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *SqlUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query, ok := (*r.queries)["exists_by_email"]
	if !ok {
		return false, core.ErrQueryNotConfigured
	}
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *SqlUserRepository) GetByEmail(ctx context.Context, email string) (*core.User, error) {
	user := &core.User{}
	query, ok := (*r.queries)["get_by_email"]
	if !ok {
		return nil, core.ErrQueryNotConfigured
	}
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}
