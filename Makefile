.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.DEFAULT_GOAL := main


SERVICE_NAME = transcoder
GOOGLE_PROJECT = videocoin
DOCKER_REGISTRY = us.gcr.io

IMAGE_TAG=$(DOCKER_REGISTRY)/$(GOOGLE_PROJECT)/$(SERVICE_NAME):$(VERSION)

VERSION=$$(git rev-parse --short HEAD)


main: deps build


deps:
	go mod verify && go mod tidy

build:
	go build -o bin/$(SERVICE_NAME) cmd/main.go

docker:
	@echo "==> Docker building..."
	docker build -t ${IMAGE_TAG} . --squash
