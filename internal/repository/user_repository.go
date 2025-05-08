package repository

import (
	"context"
	"tablelink/internal/domain"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
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
func (u *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	query := `
	INSERT INTO users (id, name, email, password, role_id, last_access)
	VALUES ($1, $2,$ 3, $4, $5, $6)`
	_, err = tx.Exec(ctx, query, user.ID, user.Name, user.Email, user.Password, user.RoleID, user.LastAccess)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()
	query := `
	UPDATE users SET name = $1, email = $2, password = $3, role_id = $4, last_access= $5 
	WHERE id = $6`
	_, err = tx.Exec(ctx, query, user.Name, user.Email, user.Password, user.RoleID, user.LastAccess)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return user, nil

}

func (u *userRepository) Delete(ctx context.Context, id int) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()
	query := `
	DELETE FROM users WHERE id=$1`
	_, err = tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil

}

func (u *userRepository) ListAll(ctx context.Context) ([]*domain.User, error) {
	users := make([]*domain.User, 0)
	query := "SELECT id, name, email, role_id, last_access FROM users"
	if err := pgxscan.Select(ctx, u.pool, users, query); err != nil {
		return nil, err
	}
	return users, nil
}
