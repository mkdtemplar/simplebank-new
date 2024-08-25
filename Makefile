postgres:
	docker run --name simplebank-new  -p 5432:5432  -e POSTGRES_USER=postgres  -e POSTGRES_PASSWORD=postgres -d postgres:latest

createdb:
	docker exec -it simplebank-new createdb --username=postgres --owner=postgres simplebankdata

migratecreate:
	migrate create -ext sql -dir db/migrations/ -seq init_schema

migrateup:
	 migrate -path db/migrations/ -database "postgresql://postgres:postgres@localhost:5432/simplebankdata?sslmode=disable" -verbose up

dropdb:
	docker exec -it schedule dropdb schedules

migratedown:
	migrate -path db/migration -database "postgresql://postgres:postgres@$5432:5432/simplebankdata?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

mockdb:
	mockgen -package mockdb -destination db/mock/store.go github.com/mkdtemplar/simplebank-new/db/sqlc Store

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc


evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 redis:latest

.PHONY: postgres createdb createtestdb dropdb migrateup migratedown migratecreate test mockdb proto evans redis