package repository

import (
	"context"
	"tablelink/internal/domain"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	CountByEmail(ctx context.Context, email string) (int, error)
	Create(ctx context.Context, user *domain.User) error
	UpdateName(ctx context.Context, user *domain.User) error
	UpdateLastAccess(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
	ListAll(ctx context.Context) ([]*domain.User, error)
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{
		pool: pool,
	}
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := new(domain.User)
	query := `
	SELECT id, name, email, password, role_id, last_access
	FROM users WHERE email = $1
	`
	if err := pgxscan.Get(ctx, u.pool, user, query, email); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepository) CountByEmail(ctx context.Context, email string) (int, error) {
	var total int
	query := `
	SELECT COUNT(*) FROM users WHERE email = $1
	`
	if err := pgxscan.Get(ctx, u.pool, &total, query, email); err != nil {
		return 0, err
	}

	return total, nil
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) error {
	err := BeginTx(ctx, u.pool, func(tx pgx.Tx) error {
		query := `INSERT INTO users (name, email, password, role_id) VALUES ($1, $2, $3, $4)`
		_, err := tx.Exec(ctx, query, user.Name, user.Email, user.Password, user.RoleID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepository) UpdateName(ctx context.Context, user *domain.User) error {
	err := BeginTx(ctx, u.pool, func(tx pgx.Tx) error {
		query := `UPDATE users SET name = $1 WHERE id = $2`
		_, err := tx.Exec(ctx, query, user.Name, user.ID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepository) UpdateLastAccess(ctx context.Context, user *domain.User) error {
	err := BeginTx(ctx, u.pool, func(tx pgx.Tx) error {
		query := `UPDATE users SET last_access= $1 WHERE id = $2`
		_, err := tx.Exec(ctx, query, user.LastAccess, user.ID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepository) Delete(ctx context.Context, id int) error {
	err := BeginTx(ctx, u.pool, func(tx pgx.Tx) error {
		query := `DELETE FROM users WHERE id=$1`
		_, err := tx.Exec(ctx, query, id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepository) ListAll(ctx context.Context) ([]*domain.User, error) {
	users := make([]*domain.User, 0)
	query := "SELECT id, name, email, role_id, last_access FROM users"
	if err := pgxscan.Select(ctx, u.pool, &users, query); err != nil {
		return nil, err
	}
	return users, nil
}
