build_client:
	@go build -o cmd/client/client cmd/client/main.go

build_server:
	@go build -o cmd/server/server cmd/server/main.go

proto:
	@rm -f internal/proto/*.go

	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/proto/user.proto

	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/proto/vault.proto

migration_down:
	@goose -dir internal/storage/postgres/migrations \
		postgres "host=192.168.0.27 port=5412 user=keeper password=keeper dbname=keeper sslmode=disable" \
		down

migration_create:
	@goose -s -dir internal/storage/postgres/migrations \
		create $(NAME) sql