postgres:
	docker run --name postgres12 -p 5432:5432 -v data:/var/lib/postgresql/data -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root quiz

dropdb:
	docker exec -it postgres12 dropdb quiz

migratestart:
	migrate create -ext sql -dir app/db/migration -seq init_schema

migrateup:
	migrate -path app/db/migration -database "postgresql://root:secret@localhost:5432/quiz?sslmode=disable" -verbose up

migratedown:
	migrate -path app/db/migration -database "postgresql://root:secret@localhost:5432/quiz?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

build:
	go build

server:
	go run main.go

.PHONY: network postgres createdb dropdb migrateup migratedown sqlc test build server