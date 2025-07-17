CLIENT_PATH=./cmd/client/main.go
CLIENT_OUTPUT_PATH=./bin/Dblockchain

SERVER_PATH=./cmd/server/main.go
SERVER_OUTPUT_PATH=./bin/Dblockchain_server

clean:
	@rm -rf bin

build_client: build-deps tailwind
	@mv ./internals/frontend/static/favicon.svg  ./internals/frontend/assets/favicon.svg
	@go build -o $(CLIENT_OUTPUT_PATH) $(CLIENT_PATH)
	@chmod +x $(CLIENT_OUTPUT_PATH)

build_server:
	@go build -o $(SERVER_OUTPUT_PATH) $(SERVER_PATH)
	@chmod +x $(SERVER_OUTPUT_PATH)

build: build_client build_server
	@echo "build done"

test:
	go test ./...

run: build
	@$(CLIENT_OUTPUT_PATH)

templ-dev:
	go tool templ generate --watch --proxy="http://localhost:3000" --open-browser=false

templ:
	go tool templ generate 

dev-server:
	go tool air \
		--build.cmd "go build -o tmp/bin/main ./cmd/client/main.go" \
		--build.bin "tmp/bin/main" \
		--build.delay "100" \
		--build.exclude_dir "node_modules" \
    --build.include_ext "go,css,js" \
		--build.stop_on_error "false" \
		--misc.clean_on_exit true  

tailwind-dev:
	bun run tailwind:dev

tailwind:
	bun run tailwind:build

build-deps:
	@bun run build-deps

dev:
	make -j3 tailwind-dev templ-dev dev-server

prod:
	make -j3 tailwind templ build


.PHONY: templ server tailwind all
