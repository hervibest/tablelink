DB_URL=postgres://postgres:postgres@localhost:5432/tablelink?sslmode=disable&TimeZone=Asia/Jakarta
MIGRATIONS_DIR = db/migrations

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-to-zero:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down-to 0

generate-proto-user:
	cd proto && protoc --go_out=. --go-grpc_out=. user.proto

generate-proto-auth:
	cd proto && protoc --go_out=. --go-grpc_out=. auth.proto

mockgen-role-right-repo:
	cd internal && mockgen -source=./repository/role_right_repository.go -destination=./repository/mock/mock_role_right_repository.go -package=mockrepo

mockgen-role-repo:
	cd internal && mockgen -source=./repository/role_repository.go -destination=./repository/mock/mock_role_repository.go -package=mockrepo

mockgen-user-repo:
	cd internal && mockgen -source=./repository/user_repository.go -destination=./repository/mock/mock_user_repository.go -package=mockrepo

mockgen-cache-adapter:
	cd internal && mockgen -source=./adapter/cache_adapter.go -destination=./adapter/mock/mock_cache_adapter.go -package=mockadapter

mockgen-auth-usecase:
	cd internal && mockgen -source=./usecase/auth_usecase.go -destination=./usecase/mock/mock_auth_usecase.go -package=mockusecase

mockgen-right-usecase:
	cd internal && mockgen -source=./usecase/right_usecase.go -destination=./usecase/mock/mock_right_usecase.go -package=mockusecase

mockgen-user-usecase:
	cd internal && mockgen -source=./usecase/user_usecase.go -destination=./usecase/mock/mock_user_usecase.go -package=mockusecase

mockgen-logger-helper:
	cd internal && mockgen -source=./helper/logger_helper.go -destination=./helper/mock/mock_logger_helper.go -package=mockhelper