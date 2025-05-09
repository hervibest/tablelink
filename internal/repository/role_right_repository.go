package repository

import (
	"context"
	"tablelink/internal/domain"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRightRepository interface {
	CheckPermission(ctx context.Context, roleID int, section, route string) (*domain.RoleRight, error)
}

type roleRightRepository struct {
	pool *pgxpool.Pool
}

func NewRoleRightRepository(pool *pgxpool.Pool) RoleRightRepository {
	return &roleRightRepository{
		pool: pool,
	}
}

func (r *roleRightRepository) CheckPermission(ctx context.Context, roleID int, section, route string) (*domain.RoleRight, error) {
	rr := new(domain.RoleRight)
	query := `
	SELECT 
		id, role_id, section, route, r_create, r_read, r_update, r_delete 
	FROM 
		role_rights
	WHERE  
		role_id = $1 AND section= $2 AND route = $3`
	if err := pgxscan.Get(ctx, r.pool, rr, query, roleID, section, route); err != nil {
		return nil, err
	}

	return rr, nil
}
