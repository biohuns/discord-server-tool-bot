APP_NAME := discord_servertool

default: build

.PHONY: build
build:
ifeq ($(OS),Windows_NT)
	go build -o $(APP_NAME).exe
else
	go build -o $(APP_NAME)
endif

.PHONY: run
run: build
ifeq ($(OS),Windows_NT)
	./$(APP_NAME).exe
else
	./$(APP_NAME)
endif

.PHONY: clean
clean:
	go clean -x
