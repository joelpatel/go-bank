test:
	docker pull postgres:16.1-alpine
	docker run --name TestDB -p 54321:5432 -e POSTGRES_DB=test-bank -e POSTGRES_USER=test-user -e POSTGRES_PASSWORD=test-password -d postgres:16.1-alpine
	go test -v -cover ./...
	docker rm -f TestDB

.PHONY: test