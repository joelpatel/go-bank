include .env

DB_URL=postgres://$(DATABASE_USER):$(DATABASE_PASS)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=$(DATABASE_SSLMODE)
TEST_DB_URL=postgres://$(TEST_DATABASE_USER):$(TEST_DATABASE_PASS)@$(TEST_DATABASE_HOST):$(TEST_DATABASE_PORT)/$(TEST_DATABASE_NAME)?sslmode=$(TEST_DATABASE_SSLMODE)

postgres:
	docker pull postgres:16.1-alpine
	docker run --name postgres16.1 -p $(DATABASE_PORT):5432 -e POSTGRES_DB=$(DATABASE_NAME) -e POSTGRES_USER=$(DATABASE_USER) -e POSTGRES_PASSWORD=$(DATABASE_PASS) -d postgres:16.1-alpine

migrateup:
	migrate -path sql -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path sql -database "$(DB_URL)" -verbose down -all

test_migrateup:
	migrate -path sql -database "$(TEST_DB_URL)" -verbose up

test_migratedown:
	migrate -path sql -database "$(TEST_DB_URL)" -verbose down -all

new_migration:
	migrate create -ext sql -dir sql -seq $(name)

pre_test:
	docker pull postgres:16.1-alpine
	docker run --name TestDB -p $(TEST_DATABASE_PORT):5432 -e POSTGRES_DB=$(TEST_DATABASE_NAME) -e POSTGRES_USER=$(TEST_DATABASE_USER) -e POSTGRES_PASSWORD=$(TEST_DATABASE_PASS) -d postgres:16.1-alpine
	sleep 1
	migrate -path sql -database "$(TEST_DB_URL)" up

test:
	go test -v -cover ./...

post_test:
	migrate -path sql -database "$(TEST_DB_URL)" down -all
	docker rm -f TestDB

execute_tests:
	make pre_test && make test && make post_test

.PHONY: postgres migrateup migratedown test_migrateup test_migratedown new_migration pre_test test post_test execute_tests