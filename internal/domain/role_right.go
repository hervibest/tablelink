package domain

type RoleRight struct {
	Id      int    `db:"id"`
	RoleId  int    `db:"role_id"`
	Section string `db:"section"`
	route   string `db:"route"`
	RCreate bool   `db:"r_create"`
	RRead   bool   `db:"r_read"`
	RUpdate bool   `db:"r_update"`
	RDelete bool   `db:"r_delete"`
}
