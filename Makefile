.PHONY: deploy
.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.DEFAULT_GOAL := docker

DOCKER_REGISTRY = us.gcr.io
CIRCLE_ARTIFACTS = ./bin
SERVICE_NAME = transcoder
RELEASE_BUCKET?=

PROJECT_ID= videocoin-network
VERSION=$$(git describe --abbrev=0)-$$(git rev-parse --short HEAD)
IMAGE_TAG=$(DOCKER_REGISTRY)/$(PROJECT_ID)/$(SERVICE_NAME):$(VERSION)
LATEST=$(DOCKER_REGISTRY)/$(PROJECT_ID)/$(SERVICE_NAME):latest


main: docker

test:
	@echo "==> Running tests..."
	go test -v -short ./...

deps:
	#env GO111MODULE=on go mod vendor
	cp -r $(GOPATH)/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1 \
	vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/

build:
	@echo "==> Building..."
	go build -o bin/$(SERVICE_NAME) cmd/main.go

build-alpine:
	@echo "==> Building for alpine..."
	go build -o bin/$(SERVICE_NAME) --ldflags '-w -linkmode external -extldflags "-static"' cmd/main.go

docker:
	@echo "==> Docker building..."
	docker build -t $(IMAGE_TAG) -t $(LATEST) .
	docker push $(IMAGE_TAG)
	docker push $(LATEST)

clean:
	rm -rf release bin

deploy:
	@cd ./deploy && ./deploy.sh

push:
	docker push $(IMAGE_TAG)
	docker push $(LATEST)
