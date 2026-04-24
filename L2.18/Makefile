.PHONY: build run test lint clean docker-build docker-run

BINARY_NAME=calendar-app
BIN_DIR=bin
CONFIG_PATH=config/local.yaml

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/server

run: build
	CONFIG_PATH=$(CONFIG_PATH) ./$(BIN_DIR)/$(BINARY_NAME)

test:
	go test -v -race ./internal/...

docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -p 8080:8080 -e CONFIG_PATH=config/local.yaml $(BINARY_NAME)

lint:
	golangci-lint run

clean:
	rm -rf $(BIN_DIR)
	go clean