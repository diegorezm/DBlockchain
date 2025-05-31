APP_API_MAIN_PATH=./cmd/api/main.go
APP_OUTPUT_PATH=./bin/Dblockchain

clean:
	@rm -rf bin

build:
	@go build -o $(APP_OUTPUT_PATH) $(APP_API_MAIN_PATH)
	@chmod +x $(APP_OUTPUT_PATH)
	@echo "build done"

test:
	go test ./...

run: build
	@$(APP_OUTPUT_PATH)

.PHONY: clean build run
