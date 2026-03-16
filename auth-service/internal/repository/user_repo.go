package repository

import (
	"context"

	"github.com/djwhocodes/auth-service/internal/model"
)

type UserRepository struct {
	*Repository
}

func NewUserRepository(repo *Repository) *UserRepository {
	return &UserRepository{repo}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {

	query := `
	INSERT INTO users(id,email,password_hash,is_active)
	VALUES($1,$2,$3,$4)
	`

	_, err := r.DB.Exec(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.IsActive,
	)

	return err
}

func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {

	query := `
	SELECT id,email,password_hash,is_active,created_at
	FROM users
	WHERE email=$1
	`

	row := r.DB.QueryRow(ctx, query, email)

	user := &model.User{}

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
