postgress:
	docker run --name go-postgress-k8 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it go-postgress-k8 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it go-postgress-k8 dropdb simple_bank

migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	docker run --rm -v C:\GoProjects\src\github.com\arun6783\go-postgress-k8s\db:/src -w /src kjconroy/sqlc generate
.PHONY: postgress createdb	dropdb migrateup migratedown sqlc