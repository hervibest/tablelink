package repository

import (
	"context"
	"tablelink/internal/domain"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository interface {
	GetByID(ctx context.Context, id int) (*domain.Role, error)
}

type roleRepository struct {
	pool *pgxpool.Pool
}

func NewRoleRepository(pool *pgxpool.Pool) RoleRepository {
	return &roleRepository{
		pool: pool,
	}
}

func (r *roleRepository) GetByID(ctx context.Context, id int) (*domain.Role, error) {
	role := new(domain.Role)
	query := `SELECT id, name FROM roles WHERE id = $1`

	if err := pgxscan.Get(ctx, r.pool, role, query, id); err != nil {
		return nil, err
	}

	return role, nil

}
