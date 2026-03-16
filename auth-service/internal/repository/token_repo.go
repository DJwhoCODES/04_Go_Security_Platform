package repository

import (
	"context"

	"github.com/djwhocodes/auth-service/internal/model"
)

type TokenRepository struct {
	*Repository
}

func NewTokenRepository(repo *Repository) *TokenRepository {
	return &TokenRepository{repo}
}

func (r *TokenRepository) Save(
	ctx context.Context,
	token *model.RefreshToken,
) error {

	query := `
	INSERT INTO refresh_tokens(id,user_id,token_hash,expires_at)
	VALUES($1,$2,$3,$4)
	`

	_, err := r.DB.Exec(
		ctx,
		query,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
	)

	return err
}

func (r *TokenRepository) Revoke(
	ctx context.Context,
	hash string,
) error {

	query := `
	UPDATE refresh_tokens
	SET revoked=true
	WHERE token_hash=$1
	`

	_, err := r.DB.Exec(ctx, query, hash)

	return err
}

func (r *TokenRepository) FindByHash(
	ctx context.Context,
	hash string,
) (*model.RefreshToken, error) {

	query := `
	SELECT id,user_id,token_hash,expires_at,revoked,created_at
	FROM refresh_tokens
	WHERE token_hash=$1
	`

	row := r.DB.QueryRow(ctx, query, hash)

	token := &model.RefreshToken{}

	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.Revoked,
		&token.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return token, nil
}
