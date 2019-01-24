.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.DEFAULT_GOAL := main

GOOS = linux
GOARCH = amd64

SERVICE_NAME = transcoder
DOCKER_REGISTRY = us.gcr.io
PROJECT_ID?=
RELEASE_BUCKET?=

IMAGE_TAG=$(DOCKER_REGISTRY)/${PROJECT_ID}/$(SERVICE_NAME):$(VERSION)
LATEST=$(DOCKER_REGISTRY)/${PROJECT_ID}/$(SERVICE_NAME):latest

VERSION=$$(git rev-parse --short HEAD)


main: package publish


test:
	@echo "==> Running tests..."
	go test -v -short ./...

deps:
	@echo "==> Running go dep..."
	go mod verify && go mod tidy

proto-update:
	env GO111MODULE=on go get github.com/videocoin/common@latest
	env GO111MODULE=on go mod vendor
	env GO111MODULE=on go mod tidy

build:
	@echo "==> Building..."
	go build -o bin/$(SERVICE_NAME) cmd/main.go

build-alpine:
	@echo "==> Building for alpine..."
	go build -o bin/$(SERVICE_NAME) --ldflags '-w -linkmode external -extldflags "-static"' cmd/main.go


docker:
	@echo "==> Docker building..."
	docker build -t ${IMAGE_TAG} -t $(LATEST) . --squash

package:
	@echo "==> Building package..."
	export GOOS=linux
	export GOARCH=amd64
	export CGO_ENABLED=0
	go build -a -installsuffix cgo -ldflags="-w -s" -o release/$(SERVICE_NAME) cmd/main.go
	tar -C release -cvjf release/$(VERSION)_transcoder_linux_amd64.tar.bz2 transcoder

publish:
	@echo "==> Pushing to storage..."
	gsutil -m cp release/$(VERSION)_transcoder_linux_amd64.tar.bz2 gs://${RELEASE_BUCKET}/transcoder/
	gsutil -m cp gs://${RELEASE_BUCKET}/transcoder/$(VERSION)_transcoder_linux_amd64.tar.bz2 gs://${RELEASE_BUCKET}/transcoder/latest_transcoder_linux_amd64.tar.bz2