CLIENT_PATH=./cmd/client/main.go
CLIENT_OUTPUT_PATH=./bin/Dblockchain

SERVER_PATH=./cmd/server/main.go
	SERVER_OUTPUT_PATH=./bin/Dblockchain_server

clean:
	@rm -rf bin

build:
	@go build -o $(CLIENT_OUTPUT_PATH) $(CLIENT_PATH)
	@chmod +x $(CLIENT_OUTPUT_PATH)
	@echo "build done"

test:
	go test ./...

run: build
	@$(CLIENT_OUTPUT_PATH)

.PHONY: clean build run
