postgres:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine
createdb:
	docker exec -it postgres17 createdb --username=root --owner=root musli

dropdb:
	docker exec -it postgres17 dropdb musli
stoppostgres:
	sudo systemctl stop postgresql
migrateup:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/musli?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/musli?sslmode=disable" -verbose down
sqlc:
	sqlc generate
server:
	go run main.go
.PHONY: createdb dropdb postgres stoppostgres migrateup migratedown server