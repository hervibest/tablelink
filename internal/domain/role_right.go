package domain

type RoleRight struct {
	Id      int    `db:"id"`
	RoleId  int    `db:"role_id"`
	Section string `db:"section"`
	Route   string `db:"route"`
	RCreate int    `db:"r_create"`
	RRead   int    `db:"r_read"`
	RUpdate int    `db:"r_update"`
	RDelete int    `db:"r_delete"`
}
