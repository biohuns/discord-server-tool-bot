APP_NAME := discord-servertool

default: build

.PHONY: prep
prep:
	go install github.com/google/wire/cmd/wire

.PHONY: build
build: prep
	go generate
	go build -o $(APP_NAME)

.PHONY: run
run: build
	./$(APP_NAME)

.PHONY: clean
clean:
	go clean -x

.PHONY: help
help:
	@cat Makefile
