GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

PKG_NAME := github.com/fatindeed/docker-compose-deployer

EXTENSION:=
ifeq ($(GOOS),windows)
  EXTENSION:=.exe
endif

STATIC_FLAGS=CGO_ENABLED=0

GIT_TAG?=$(shell git rev-parse --short HEAD)

LDFLAGS="-s -w -X $(PKG_NAME)/cmd.Version=${GIT_TAG}"
GO_BUILD=$(STATIC_FLAGS) go build -trimpath -ldflags=$(LDFLAGS)

COMPOSE_BINARY?=bin/docker-compose-deployer
COMPOSE_BINARY_WITH_EXTENSION=$(COMPOSE_BINARY)$(EXTENSION)

TAGS:=
ifdef BUILD_TAGS
  TAGS=-tags $(BUILD_TAGS)
  LINT_TAGS=--build-tags $(BUILD_TAGS)
endif

all: test build

.PHONY: test
test:
	go test $(TAGS) -cover ./...

.PHONY: lint
lint:
	golangci-lint run $(LINT_TAGS) --timeout 10m0s ./...

.PHONY: build
build:
	$(GO_BUILD) $(TAGS) -o $(COMPOSE_BINARY_WITH_EXTENSION) main.go

.PHONY: cross
cross:
	GOOS=linux   GOARCH=amd64 $(GO_BUILD) $(TAGS) -o $(COMPOSE_BINARY)-linux-x86_64 main.go
	GOOS=linux   GOARCH=arm64 $(GO_BUILD) $(TAGS) -o $(COMPOSE_BINARY)-linux-aarch64 main.go
	GOOS=linux   GOARM=6 GOARCH=arm $(GO_BUILD) $(TAGS) -o $(COMPOSE_BINARY)-linux-armv6 main.go
	GOOS=linux   GOARM=7 GOARCH=arm $(GO_BUILD) $(TAGS) -o $(COMPOSE_BINARY)-linux-armv7 main.go
	GOOS=linux   GOARCH=s390x $(GO_BUILD) $(TAGS) -o $(COMPOSE_BINARY)-linux-s390x main.go
	GOOS=darwin  GOARCH=amd64 $(GO_BUILD) $(TAGS) -o $(COMPOSE_BINARY)-darwin-x86_64 main.go
	GOOS=darwin  GOARCH=arm64 $(GO_BUILD) $(TAGS) -o $(COMPOSE_BINARY)-darwin-aarch64 main.go
	GOOS=windows GOARCH=amd64 $(GO_BUILD) $(TAGS) -o $(COMPOSE_BINARY)-windows-x86_64.exe main.go

docker-build:
	DOCKER_BUILDKIT=0 docker build -t fatindeed/docker-compose-deployer .