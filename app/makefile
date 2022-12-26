## 注意變數不要跟path重疊
POSTGRESQL_URL="postgres://daniel:mypassword@localhost:5432/test_db?sslmode=disable"
MIGRATION_PATH="db/migrations"

migrate-up:
		migrate -path ${MIGRATION_PATH} -database ${POSTGRESQL_URL} up
migrate-down:
		migrate -path ${MIGRATION_PATH} -database ${POSTGRESQL_URL} down
generateSql:
		sqlc generate
## v verbose log  ./... all package
test:
		go test -v -cover ./...