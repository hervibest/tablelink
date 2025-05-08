package domain

import "time"

type User struct {
	ID         int        `db:"id"`
	Name       string     `db:"name"`
	Email      string     `db:"email"`
	Password   string     `db:"password"`
	RoleID     string     `db:"role_id"`
	LastAccess *time.Time `db:"last_access"`
}
