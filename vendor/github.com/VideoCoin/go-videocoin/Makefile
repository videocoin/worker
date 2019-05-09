.PHONY: vid android ios vid-cross evm all test clean
.PHONY: vid-linux vid-linux-386 vid-linux-amd64 vid-linux-mips64 vid-linux-mips64le
.PHONY: vid-linux-arm vid-linux-arm-5 vid-linux-arm-6 vid-linux-arm-7 vid-linux-arm64
.PHONY: vid-darwin vid-darwin-386 vid-darwin-amd64
.PHONY: vid-windows vid-windows-386 vid-windows-amd64
.PHONY: vid-compress
.PHONY: faucet bootnode proxy

GOBIN = $(shell pwd)/build/bin
GO ?= latest

vid:
	build/env.sh go run build/ci.go install ./cmd/vid
	@echo "Done building."
	@echo "Run \"$(GOBIN)/vid\" to launch vid."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/vid.aar\" to use the library."

ios:
	build/env.sh go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/vid.framework\" to use the library."

test: all
	build/env.sh go run build/ci.go test

lint: ## Run linters.
	build/env.sh go run build/ci.go lint

clean:
	./build/clean_go_build_cache.sh
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

#vid-cross: vid-linux vid-darwin vid-windows vid-android vid-ios
#	@echo "Full cross compilation done:"
#	@ls -ld $(GOBIN)/vid-*

vid-cross:
	build/env.sh go run build/ci.go xgo -- --go=latest --targets=linux/amd64,linux/arm64,darwin/amd64,windows/amd64 -v ./cmd/vid
	mv build/bin/vid-darwin-10.6-amd64 build/bin/vid-darwin-amd64
	mv build/bin/vid-windows-4.0-amd64.exe build/bin/vid-windows-amd64.exe

vid-linux: vid-linux-386 vid-linux-amd64 vid-linux-arm vid-linux-mips64 vid-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/vid-linux-*

vid-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/vid
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/vid-linux-* | grep 386

vid-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/vid
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/vid-linux-* | grep amd64

vid-linux-arm: vid-linux-arm-5 vid-linux-arm-6 vid-linux-arm-7 vid-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/vid-linux-* | grep arm

vid-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/vid
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/vid-linux-* | grep arm-5

vid-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/vid
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/vid-linux-* | grep arm-6

vid-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/vid
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/vid-linux-* | grep arm-7

vid-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/vid
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/vid-linux-* | grep arm64

vid-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/vid
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/vid-linux-* | grep mips

vid-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/vid
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/vid-linux-* | grep mipsle

vid-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/vid
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/vid-linux-* | grep mips64

vid-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/vid
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/vid-linux-* | grep mips64le

vid-darwin: vid-darwin-386 vid-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/vid-darwin-*

vid-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/vid
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/vid-darwin-* | grep 386

vid-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/vid
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/vid-darwin-* | grep amd64

vid-windows: vid-windows-386 vid-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/vid-windows-*

vid-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/vid
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/vid-windows-* | grep 386

vid-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/vid
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/vid-windows-* | grep amd64

vid-compress:
	cd build/bin; for f in vid*; do zip $$f.zip $$f; rm $$f; done; cd -

faucet:
	build/env.sh go run build/ci.go install ./cmd/faucet
	@echo "Done building."
	@echo "Run \"$(GOBIN)/faucet\" to launch faucet."

proxy:
	build/env.sh go run build/ci.go install ./cmd/proxy
	@echo "Done building."
	@echo "Run \"$(GOBIN)/proxy\" to launch proxy."

bootnode:
	build/env.sh go run build/ci.go install ./cmd/bootnode
	@echo "Done building."
	@echo "Run \"$(GOBIN)/bootnode\" to launch bootnode."