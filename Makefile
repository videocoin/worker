.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.DEFAULT_GOAL := main

GOOS = linux
GOARCH = amd64

SERVICE_NAME = transcoder
GOOGLE_PROJECT = videocoin
DOCKER_REGISTRY = us.gcr.io

IMAGE_TAG=$(DOCKER_REGISTRY)/$(GOOGLE_PROJECT)/$(SERVICE_NAME):$(VERSION)
LATEST=$(DOCKER_REGISTRY)/$(GOOGLE_PROJECT)/$(SERVICE_NAME):latest

VERSION=$$(git rev-parse --short HEAD)


main: deps build


deps:
	go mod verify && go mod tidy

build:
	go build -o bin/$(SERVICE_NAME) cmd/main.go

build-alpine:
	go build -o bin/$(SERVICE_NAME) --ldflags '-w -linkmode external -extldflags "-static"' cmd/main.go


docker:
	@echo "==> Docker building..."
	docker build -t ${IMAGE_TAG} -t $(LATEST) . --squash
