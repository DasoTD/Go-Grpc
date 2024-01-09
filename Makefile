DB_URL=postgresql://root:secret@localhost:5432/cypherdb?sslmode=disable

postgres:
	sudo docker run -d --name cypher -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres
creatdb:
	sudo docker exec -it cypher createdb --username=root --owner=root cypherdb

dropdb:
	# sudo docker exec -it cypher dropdb --username=root --owner=root cypherdb
	sudo docker exec -it cypher dropdb cypherdb

connectDB:
	sudo docker exec -it cypher psql -U root

makemigrattion:
	migrate create -ext sql -dir db/migration -seq add_users

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate
server:
	go run main.go
proto:
	# rm -f pb/*.go
	# rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc
.PHONY: creatdb dropdb postgres migratedown migrateup migratedown1 migrateup1 test server mock docker sqlc


# protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
# 	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
# 	proto/*.proto



# protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative  --go-grpc_out=pb --go-grpc_opt=paths=source_relative  proto/*.proto

# go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# scoop install migrate