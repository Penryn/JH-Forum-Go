.PHONY: all build build-web run test clean fmt pre-commit help

PROJECT = JH-Forum-Go
TARGET = JH-Forum
ifeq ($(OS),Windows_NT)
TARGET := $(TARGET).exe
endif
TARGET_BIN = $(basename $(TARGET))

ifeq (n$(CGO_ENABLED),n)
CGO_ENABLED := 1
endif

RELEASE_ROOT = release
RELEASE_FILES = LICENSE README.md  config.yaml.example docker-compose.yaml scripts docs
RELEASE_LINUX_AMD64 = $(RELEASE_ROOT)/linux-amd64/$(PROJECT)
RELEASE_DARWIN_AMD64 = $(RELEASE_ROOT)/darwin-amd64/$(PROJECT)
RELEASE_DARWIN_ARM64 = $(RELEASE_ROOT)/darwin-arm64/$(PROJECT)
RELEASE_WINDOWS_AMD64 = $(RELEASE_ROOT)/windows-amd64/$(PROJECT)

BUILD_VERSION := $(shell git describe --tags --always)
BUILD_DATE := $(shell date +'%Y-%m-%d %H:%M:%S %Z')
SHA_SHORT := $(shell git rev-parse --short HEAD)

MOD_NAME = JH-Forum
LDFLAGS = -X "${MOD_NAME}/pkg/version.version=${BUILD_VERSION}" \
          -X "${MOD_NAME}/pkg/version.buildDate=${BUILD_DATE}" \
          -X "${MOD_NAME}/pkg/version.commitID=${SHA_SHORT}" \
          -X "${MOD_NAME}/pkg/version.buildTags=${TAGS}" \
		  -w -s

all: fmt build

build:
	@go mod download
	@echo Build JH-Forum
	@go build -trimpath -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(RELEASE_ROOT)/$(TARGET)

run:
	@go run -trimpath -gcflags "all=-N -l" -tags '$(TAGS)' -ldflags '$(LDFLAGS)' . serve

.PHONY: release
release: linux-amd64 darwin-amd64 darwin-arm64 windows-x64
	@echo Package JH-Forum
	@cp -rf $(RELEASE_FILES) $(RELEASE_LINUX_AMD64)
	@cp -rf $(RELEASE_FILES) $(RELEASE_DARWIN_AMD64)
	@cp -rf $(RELEASE_FILES) $(RELEASE_DARWIN_ARM64)
	@cp -rf $(RELEASE_FILES) $(RELEASE_WINDOWS_AMD64)
	@cd $(RELEASE_LINUX_AMD64)/.. && rm -f *.zip && zip -r $(PROJECT)-linux_amd64.zip $(PROJECT) && cd -
	@cd $(RELEASE_DARWIN_AMD64)/.. && rm -f *.zip && zip -r $(PROJECT)-darwin_amd64.zip $(PROJECT) && cd -
	@cd $(RELEASE_DARWIN_ARM64)/.. && rm -f *.zip && zip -r $(PROJECT)-darwin_arm64.zip $(PROJECT) && cd -
	@cd $(RELEASE_WINDOWS_AMD64)/.. && rm -f *.zip && zip -r $(PROJECT)-windows_amd64.zip $(PROJECT) && cd -

.PHONY: linux-amd64
linux-amd64:
	@echo Build JH-Forum [linux-amd64] CGO_ENABLED=$(CGO_ENABLED) TAGS="'$(TAGS)'"
	@set "CGO_ENABLED=$(CGO_ENABLED)" && set "GOOS=linux" && set "GOARCH=amd64" && go build -trimpath -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(RELEASE_LINUX_AMD64)/$(TARGET_BIN)

.PHONY: darwin-amd64
darwin-amd64:
	@echo Build JH-Forum [darwin-amd64] CGO_ENABLED=$(CGO_ENABLED) TAGS="'$(TAGS)'"
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=amd64 go build -trimpath  -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(RELEASE_DARWIN_AMD64)/$(TARGET_BIN)

.PHONY: darwin-arm64
darwin-arm64:
	@echo Build JH-Forum [darwin-arm64] CGO_ENABLED=$(CGO_ENABLED) TAGS="'$(TAGS)'"
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=arm64 go build -trimpath -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(RELEASE_DARWIN_ARM64)/$(TARGET_BIN)

.PHONY: windows-x64
windows-x64:
	@echo Build JH-Forum [windows-x64] CGO_ENABLED=$(CGO_ENABLED) TAGS="'$(TAGS)'"
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=amd64 go build -trimpath  -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(RELEASE_WINDOWS_AMD64)/$(TARGET_BIN).exe

.PHONY: generate
generate: gen-mir

.PHONY: gen-mir
gen-mir:
	@go generate mirc/gen.go
	@go fmt ./mirc/auto/api/...


clean:
	@go clean
	@find ./release -type f -exec rm -r {} +

fmt:
	@echo Formatting...
	@go fmt ./internal/...
	@go fmt ./pkg/...
	@go vet -composites=false ./internal/...
	@go vet -composites=false ./pkg/...

test:
	@go test ./...

pre-commit: fmt
	@go mod tidy

.PHONY: install-protobuf-plugins
install-protobuf-plugins:
	@go install github.com/bufbuild/buf/cmd/buf@v1.25.0
	@go install github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@v1.25.0
	@go install github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@v1.25.0
	@go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest
	@go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

help:
	@echo "make: make"
	@echo "make run: start api server"
	@echo "make build: build executable"
	@echo "make release: build release executables"
	@echo "make generate: generate code"
	@echo "make clean: clean up"
	@echo "make fmt: format code"
	@echo "make test: run tests"
	@echo "make pre-commit: run pre-commit checks"
	@echo "make install-protobuf-plugins: install protobuf plugins"
	
