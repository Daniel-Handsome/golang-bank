## 注意變數不要跟path重疊
## localhost是因為github ci 有postgres的container
POSTGRESQL_URL="postgres://daniel:mypassword@localhost:5432/test_db?sslmode=disable"
MIGRATION_PATH="db/migrations"

AWS_POSTGRESQL_URL="postgres://postgres:X569B210@simple-bank.c3pxclhxzdp0.us-east-1.rds.amazonaws.com/simple_bank"

migrate-up:
		migrate -path ${MIGRATION_PATH} -database ${POSTGRESQL_URL} up
migrate-down:
		migrate -path ${MIGRATION_PATH} -database ${POSTGRESQL_URL} down
generateSql:
		sqlc generate
## v verbose log  ./... all package
test:
		go test -v -cover ./...
run :
		go run main.go
mockStore:
		mockgen --build_flags=--mod=mod -package mockdb  -destination db/mock/store.go github.com/daniel/master-golang/db/sqlc Store
### 安裝3.14 用apt 最不夠版本 裝最新的也會有問題 裝3.14.0剛好
proto:
		rm -f ./pb/*.go
		rm -f ./docs/swagger/*.json
		protoc --proto_path=proto --go_out=./pb --go_opt=paths=source_relative \
		--go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=./pb --grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_opt logtostderr=true \
		--openapiv2_out ./docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
		proto/*.proto

.PHONY: run test generateSql migrate-up migrate-down mockStore proto