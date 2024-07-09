postgres:
	docker run --name simplebank -p 5432:5432  -e POSTGRES_USER=postgres  -e POSTGRES_PASSWORD=postgres -d postgres:latest

createdb:
	docker exec -it simplebank createdb --username=postgres --owner=postgres simplebankdata

migratecreate:
	migrate create -ext sql -dir repository/migrations/ -seq init_schema

migrateup:
	 migrate -path repository/migrations/ -database "postgresql://postgres:postgres@localhost:5432/react_graphql_data?sslmode=disable" -verbose up

dropdb:
	docker exec -it schedule dropdb schedules

migratedown:
	migrate -path db/migration -database "postgresql://postgres:postgres@$5432:5432/react_graphql_data?sslmode=disable" -verbose down

.PHONY: postgres createdb createtestdb dropdb migrateup migratedown migratecreate