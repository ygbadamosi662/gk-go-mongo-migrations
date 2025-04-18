# === CONFIG ===
CMD_DIR=cmd/gk
BINARY_NAME=gk
MODULE_PATH=github.com/ygbadamosi662/gk-go-mongo-migrations

# === INSTALL & BUILD ===

install:
	go install $(MODULE_PATH)/$(CMD_DIR)@latest

build:
	go build -o $(BINARY_NAME) ./$(CMD_DIR)

# === LOCAL DEV RUN ===

run:
	go run ./$(CMD_DIR) $(args)

# === SHORTCUTS ===

init:
	go run ./$(CMD_DIR) init

generate:
	go run ./$(CMD_DIR) generate name=$(name)

# === VERSIONING ===

release:
	git tag $(version)
	git push origin $(version)

# === CLEAN ===

clean:
	rm -f $(BINARY_NAME)
