.NOTPARALLEL:
.DEFAULT_GOAL := push

GOOS = linux
GOARCH = amd64
GCP_PROJECT ?= videocoin-network
NAME = worker
BRANCH=$$(git branch | grep \* | cut -d ' ' -f2)
VERSION?=$$(git describe --abbrev=0)-$$(git rev-parse --abbrev-ref HEAD)-$$(git rev-parse --short HEAD)

IMAGE_TAG=gcr.io/$(GCP_PROJECT)/$(NAME):$(VERSION)
LATEST_IMAGE_TAG=gcr.io/$(GCP_PROJECT)/$(NAME):latest

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
	docker build -f Dockerfile.lint .

build:
	@echo "==> Building..."
	go build -mod vendor -a -installsuffix cgo -ldflags="-w -s -X main.Version=${VERSION}" -o bin/$(NAME) cmd/main.go

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


release: docker docker-push docker-latest

publish:
	@echo "==> Building..."
	GOOS=linux GOARCH=amd64 go build -mod vendor -a -installsuffix cgo -ldflags="-w -s -X main.Version=${VERSION}" -o bin/$(NAME)-linux-amd64 cmd/main.go
	gsutil cp bin/worker-linux-amd64 gs://videocoin-releases/worker/${VERSION}/worker-linux-amd64
	gsutil acl ch -u AllUsers:R gs://videocoin-releases/worker/${VERSION}/worker-linux-amd64
	gsutil cp capacity_test.mp4 gs://videocoin-releases/worker/${VERSION}/capacity_test.mp4
	gsutil acl ch -u AllUsers:R gs://videocoin-releases/worker/${VERSION}/capacity_test.mp4


.PHONY : build deps test push clean docker release
