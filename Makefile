# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: godf android ios godf-cross evm all test clean
.PHONY: godf-linux godf-linux-386 godf-linux-amd64 godf-linux-mips64 godf-linux-mips64le
.PHONY: godf-linux-arm godf-linux-arm-5 godf-linux-arm-6 godf-linux-arm-7 godf-linux-arm64
.PHONY: godf-darwin godf-darwin-386 godf-darwin-amd64
.PHONY: godf-windows godf-windows-386 godf-windows-amd64

GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run

godf:
	$(GORUN) build/ci.go install ./cmd/godf
	@echo "Done building."
	@echo "Run \"$(GOBIN)/godf\" to launch godf."

all:
	$(GORUN) build/ci.go install

android:
	$(GORUN) build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/godf.aar\" to use the library."
	@echo "Import \"$(GOBIN)/godf-sources.jar\" to add javadocs"
	@echo "For more info see https://stackoverflow.com/questions/20994336/android-studio-how-to-attach-javadoc"
	
ios:
	$(GORUN) build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Godf.framework\" to use the library."

test: all
	$(GORUN) build/ci.go test

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean:
	env GO111MODULE=on go clean -cache
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

godf-cross: godf-linux godf-darwin godf-windows godf-android godf-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/godf-*

godf-linux: godf-linux-386 godf-linux-amd64 godf-linux-arm godf-linux-mips64 godf-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/godf-linux-*

godf-linux-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/godf
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/godf-linux-* | grep 386

godf-linux-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/godf
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/godf-linux-* | grep amd64

godf-linux-arm: godf-linux-arm-5 godf-linux-arm-6 godf-linux-arm-7 godf-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/godf-linux-* | grep arm

godf-linux-arm-5:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/godf
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/godf-linux-* | grep arm-5

godf-linux-arm-6:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/godf
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/godf-linux-* | grep arm-6

godf-linux-arm-7:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/godf
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/godf-linux-* | grep arm-7

godf-linux-arm64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/godf
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/godf-linux-* | grep arm64

godf-linux-mips:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/godf
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/godf-linux-* | grep mips

godf-linux-mipsle:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/godf
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/godf-linux-* | grep mipsle

godf-linux-mips64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/godf
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/godf-linux-* | grep mips64

godf-linux-mips64le:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/godf
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/godf-linux-* | grep mips64le

godf-darwin: godf-darwin-386 godf-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/godf-darwin-*

godf-darwin-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/godf
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/godf-darwin-* | grep 386

godf-darwin-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/godf
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/godf-darwin-* | grep amd64

godf-windows: godf-windows-386 godf-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/godf-windows-*

godf-windows-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/godf
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/godf-windows-* | grep 386

godf-windows-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/godf
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/godf-windows-* | grep amd64
