.PHONY: deploy
.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.DEFAULT_GOAL := docker

DOCKER_REGISTRY = us.gcr.io
CIRCLE_ARTIFACTS = ./bin
SERVICE_NAME = transcoder
RELEASE_BUCKET?=

PROJECT_ID=$$(gcloud config list --format 'value(core.project)' 2>/dev/null)
VERSION=$$(git describe --abbrev=0)-$$(git rev-parse --short HEAD)
IMAGE_TAG=$(DOCKER_REGISTRY)/${PROJECT_ID}/$(SERVICE_NAME):$(VERSION)
LATEST=$(DOCKER_REGISTRY)/${PROJECT_ID}/$(SERVICE_NAME):latest


main: docker


test:
	@echo "==> Running tests..."
	go test -v -short ./...


deps:
	env GO111MODULE=on go mod vendor
	cp -r $(GOPATH)/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1 \
	vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/

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
	docker build -t $(IMAGE_TAG) -t $(LATEST) .

package:
	cd cmd && xgo --targets=linux/amd64 -dest ../release -out $(SERVICE_NAME) .
	cp keys/$(SERVICE_NAME).key release
	tar -C release -cvjf release/$(VERSION)_$(SERVICE_NAME)_linux_amd64.tar.bz2 $(SERVICE_NAME)-linux-amd64 $(SERVICE_NAME).key

store:
	gsutil -m cp release/$(VERSION)_$(SERVICE_NAME)_linux_amd64.tar.bz2 gs://${RELEASE_BUCKET}/$(SERVICE_NAME)/
	gsutil -m cp gs://${RELEASE_BUCKET}/$(SERVICE_NAME)/$(VERSION)_$(SERVICE_NAME)_linux_amd64.tar.bz2 gs://${RELEASE_BUCKET}/$(SERVICE_NAME)/latest_$(SERVICE_NAME)_linux_amd64.tar.bz2
	
clean:
	rm -rf release/*

deploy:
	@cd ./deploy && ./deploy.sh

	
push:
	docker push $(IMAGE_TAG)
	docker push $(LATEST)
