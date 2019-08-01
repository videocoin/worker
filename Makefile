.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.DEFAULT_GOAL := push
DOCKER_REGISTRY = us.gcr.io
PROJECT_ID= videocoin-network
SERVICE_NAME = transcoder

VERSION=$$(git describe --abbrev=0)-$$(git rev-parse --short HEAD)
IMAGE_TAG=$(DOCKER_REGISTRY)/$(PROJECT_ID)/$(SERVICE_NAME):$(VERSION)
LATEST=$(DOCKER_REGISTRY)/${PROJECT_ID}/$(SERVICE_NAME):latest


version:
	@echo $(VERSION)

image-tag:
	@echo $(IMAGE_TAG)

deps:
	env GO111MODULE=on go mod vendor
	cp -r $(GOPATH)/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1 \
	vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/

build:
	export GOOS=linux
	export GOARCH=amd64
	export CGO_ENABLED=0
	@echo "==> Building..."
	go build -a -installsuffix cgo -ldflags="-w -s" -o bin/$(SERVICE_NAME) cmd/main.go

test:
	@echo "==> Running tests..."
	go test -v ./...

test-coverage:
	@echo "==> Running tests..."
	go test -cover ./...

docker:
	@echo "==> Docker building..."
	docker build -t ${IMAGE_TAG} -t $(LATEST) .

docker-push:
	docker push $(IMAGE_TAG)
	docker push $(LATEST)

clean:
	rm -rf bin/*

deploy:
	@cd deploy && ./deploy.sh

push: docker docker-push

.PHONY : build deps test push clean docker deploy release
