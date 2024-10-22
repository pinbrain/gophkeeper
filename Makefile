build_client:
	@go build -o cmd/client/client cmd/client/main.go

build_server:
	@go build -o cmd/server/server cmd/server/main.go

CERT_DNS = localhost
CERT_IP = 0.0.0.0
certs:
	@rm -rf cert
	@mkdir cert
	@echo "subjectAltName=DNS:$(CERT_DNS),IP:$(CERT_IP)" > cert/server-ext.cnf
	@openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout cert/ca-key.pem -out cert/ca-cert.pem -subj "/C=RU/ST=Moscow/O=Learn/OU=Education/CN=*"
	@openssl x509 -in cert/ca-cert.pem -noout -text
	@openssl req -newkey rsa:4096 -nodes -keyout cert/server-key.pem -out cert/server-req.pem -subj "/C=RU/ST=Moscow/L=Moscow/O=Learn/OU=Education/CN=*"
	@openssl x509 -req -in cert/server-req.pem -days 60 -CA cert/ca-cert.pem -CAkey cert/ca-key.pem -CAcreateserial -out cert/server-cert.pem -extfile cert/server-ext.cnf
	@openssl x509 -in cert/server-cert.pem -noout -text

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