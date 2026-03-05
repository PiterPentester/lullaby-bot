.PHONY: build test docker-build clean

APP_NAME = lullaby-bot
DOCKER_IMAGE = $(APP_NAME):latest

build:
	go build -o bin/lullaby-bot cmd/bot/main.go

test:
	go test ./...

run: build
	./bin/lullaby-bot

docker-build:
	docker build --platform linux/arm64 -t $(DOCKER_IMAGE) .

clean:
	rm -rf bin/
