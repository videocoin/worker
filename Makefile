.NOTPARALLEL:
.DEFAULT_GOAL := push
DOCKER_REGISTRY = gcr.io
PROJECT_ID= videocoin-network
SERVICE_NAME = transcoder
BRANCH=$$(git branch | grep \* | cut -d ' ' -f2)

VERSION=$$(git describe --abbrev=0)-$$(git rev-parse --abbrev-ref HEAD)-$$(git rev-parse --short HEAD)
IMAGE_TAG=$(DOCKER_REGISTRY)/$(PROJECT_ID)/$(SERVICE_NAME):$(VERSION)
LATEST_IMAGE_TAG=$(DOCKER_REGISTRY)/$(PROJECT_ID)/$(SERVICE_NAME):latest

ENV?=dev

version:
	@echo ${VERSION}

image-tag:
	@echo $(IMAGE_TAG)

deps:
	GO111MODULE=on go mod vendor
	cp -r $(GOPATH)/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1 \
	vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/
	cp -r $(GOPATH)/src/github.com/shirou/gopsutil/host/include \
	vendor/github.com/shirou/gopsutil/host/include
	cp -r $(GOPATH)/src/github.com/karalabe/usb/hidapi \
	vendor/github.com/karalabe/usb/hidapi
	cp -r $(GOPATH)/src/github.com/karalabe/usb/libusb \
	vendor/github.com/karalabe/usb/libusb

lint: docker-lint

docker-lint:
	docker run --rm \
		-v `PWD`:/go/src/github.com/videocoin/transcode \
		-w /go/src/github.com/videocoin/transcode \
		golangci/golangci-lint:v1.23.6 \
		golangci-lint run -v

build:
	@echo "==> Building..."
	go build -mod vendor -a -installsuffix cgo -ldflags="-w -s -X main.Version=${VERSION}" -o bin/$(SERVICE_NAME) cmd/main.go

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

docker-latest:
	docker tag $(IMAGE_TAG) $(LATEST_IMAGE_TAG)
	docker push $(LATEST_IMAGE_TAG)

clean:
	rm -rf bin/*

deploy:
	ENV=${ENV} deploy/deploy.sh

release: docker docker-push docker-latest

.PHONY : build deps test push clean docker deploy release
