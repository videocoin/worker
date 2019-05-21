.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.DEFAULT_GOAL := docker

DOCKER_REGISTRY = us.gcr.io
CIRCLE_ARTIFACTS = ./bin
SERVICE_NAME = transcoder
RELEASE_BUCKET?=

PROJECT_ID=$$(gcloud config list --format 'value(core.project)' 2>/dev/null)
VERSION=$$(git rev-parse --short HEAD)
IMAGE_TAG=$(DOCKER_REGISTRY)/${PROJECT_ID}/$(SERVICE_NAME):$(VERSION)
LATEST=$(DOCKER_REGISTRY)/${PROJECT_ID}/$(SERVICE_NAME):latest

IMAGE_TARBALL_PATH=$(CIRCLE_ARTIFACTS)/$(SERVICE_NAME)-$(VERSION).tar

version:
	@echo $(VERSION)

tag:
	@echo $(IMAGE_TAG)

deps:
	env GO111MODULE=on go mod vendor
	cp -r $(GOPATH)/src/github.com/VideoCoin/go-videocoin/crypto/secp256k1/libsecp256k1 \
	vendor/github.com/VideoCoin/go-videocoin/crypto/secp256k1/

build:
	export GOOS=linux
	export GOARCH=amd64
	export CGO_ENABLED=0
	@echo "==> Building..."
	go build -a -installsuffix cgo -ldflags="-w -s" -o bin/$(SERVICE_NAME) cmd/main.go

test:
	@echo "==> Running tests..."
	go test -v ./...

coverage:
	@echo "==> Running tests..."
	go test -cover ./...

docker: deps
	@echo "==> Docker building..."
	docker build -t $(IMAGE_TAG) -t $(LATEST) . --squash

push:
	docker push $(IMAGE_TAG)
	docker push $(LATEST)

clean:
	rm -rf release/*


.PHONY : build deps test push clean docker 
