DB_URL=postgres://postgres:postgres@localhost:5432/tablelink?sslmode=disable&TimeZone=Asia/Jakarta
MIGRATIONS_DIR = db/migrations

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up
