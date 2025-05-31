clean:
	@rm -rf bin

build:
	@go build -o ./bin/DBlockchain ./cmd/api/main.go
	@chmod +x ./bin/DBlockchain
	@echo "build done"

run: build
	@./bin/DBlockchain

.PHONY: clean build run
