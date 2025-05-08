DB_URL=postgres://postgres:postgres@localhost:5432/tablelink?sslmode=disable&TimeZone=Asia/Jakarta
MIGRATIONS_DIR = db/migrations

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

generate-proto-user:
	cd proto && protoc --go_out=. --go-grpc_out=. user.proto

generate-proto-auth:
	cd proto && protoc --go_out=. --go-grpc_out=. auth.proto