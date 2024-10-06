build_client:
	@go build -o cmd/client/client cmd/client/main.go

build_server:
	@go build -o cmd/server/server cmd/server/main.go