.NOTPARALLEL:
.DEFAULT_GOAL := push
DOCKER_REGISTRY = gcr.io
PROJECT_ID= videocoin-network
SERVICE_NAME = transcoder
BRANCH=$$(git branch | grep \* | cut -d ' ' -f2)

VERSION=$$(git describe --abbrev=0)-$$(git rev-parse --abbrev-ref HEAD)-$$(git rev-parse --short HEAD)
IMAGE_TAG=$(DOCKER_REGISTRY)/$(PROJECT_ID)/$(SERVICE_NAME):$(VERSION)
TRANSINIT_IMAGE_TAG=$(DOCKER_REGISTRY)/$(PROJECT_ID)/transinit:$(VERSION)

ENV?=snb

version:
	@echo ${VERSION}

image-tag:
	@echo $(IMAGE_TAG)

deps:
	GO111MODULE=on go mod vendor
	cp -r $(GOPATH)/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1 \
	vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/

build:
	@echo "==> Building..."
	go build -a -installsuffix cgo -ldflags="-w -s -X main.Version=${VERSION}" -o bin/$(SERVICE_NAME) cmd/main.go

build-transinit:
	@echo "==> Building transinit..."
	go build -ldflags="-w -s" -o bin/transinit cmd/transinit.go

test:
	@echo "==> Running tests..."
	go test -v ./...

test-coverage:
	@echo "==> Running tests..."
	go test -cover ./...

docker:
	@echo "==> Docker building..."
	docker build -t ${IMAGE_TAG} .

docker-push:
	docker push $(IMAGE_TAG)

clean:
	rm -rf bin/*

deploy:
	ENV=${ENV} deploy/deploy.sh

release: docker docker-push

.PHONY : build deps test push clean docker deploy release
