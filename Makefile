build:
	@go build -o bin/gobank
run: build
	@bin/gobank
test:
	@go test -v ./...
dbup:
	@docker compose up -d --remove-orphans
dbdown:
	@docker compose down
createdb:
	@docker exec -it pg createdb --username=postgres --owner=postgres dev
dropdb:
	@docker exec -it pg dropdb --username=postgres dev