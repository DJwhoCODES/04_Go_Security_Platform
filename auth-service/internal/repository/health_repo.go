package repository

import (
	"context"
)

type HealthRepository struct {
	*Repository
}

func NewHealthRepository(repo *Repository) *HealthRepository {
	return &HealthRepository{repo}
}

func (r *HealthRepository) CheckDB(ctx context.Context) error {
	return r.DB.Ping(ctx)
}
