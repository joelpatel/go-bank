postgres:
	docker run --name postgres-db-container -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:15-alpine

postgres-start:
	docker start postgres-db-container

postgres-stop:
	docker stop postgres-db-container

createdb:
	docker exec -it postgres-db-container createdb --username=root --owner=root bank

dropdb:
	docker exec -it postgres-db-container dropdb bank

migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test