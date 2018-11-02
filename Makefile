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

rr:
	export GOOS=linux
	export GOARCH=amd64
	export CGO_ENABLED=0

	go build -a -installsuffix cgo -ldflags="-w -s" -o release/$(SERVICE_NAME) cmd/main.go
	tar -C release -cvjf release/$(VERSION)_transcoder_linux_amd64.tar.bz2 transcoder

publish:
	gsutil -m cp release/$(VERSION)_transcoder_linux_amd64.tar.bz2 gs://vc-releases/transcoder/
	gsutil -m cp gs://vc-releases/transcoder/$(VERSION)_transcoder_linux_amd64.tar.bz2 gs://vc-releases/transcoder/latest_transcoder_linux_amd64.tar.bz2