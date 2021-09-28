APP_REPO ?= quay.io/synpse/metrics-nats-example-app
TAG_NAME = $(shell git rev-parse --short=7 HEAD)$(shell [[ $$(git status --porcelain) = "" ]] || echo -dirty)
OUTPUT_BIN_APP ?= release/app
LDFLAGS		+= -s -w
GOARCH ?= arm64

.PHONY: app
app:
	CGO_ENABLED=0 GOARCH=${GOARCH} go build -ldflags "$(LDFLAGS)" -o ${OUTPUT_BIN_APP}/app ./cmd/app

.PHONY: build
build: 
	docker build --build-arg GOARCH=${GOARCH} . -t ${APP_REPO}:${GOARCH} -f Dockerfile 

.PHONY: push
push:
	docker push ${APP_REPO}:${GOARCH}

.PHONY: build-arm64
build-arm64: 
	GOARCH=arm64 make build

.PHONY: build-arm
build-arm:
	GOARCH=arm make build

.PHONY: build-amd64
build-amd64:
	GOARCH=amd64 make build

.PHONY: push-arm64
push-arm64: build-arm64
	GOARCH=arm64 make push

.PHONY: push-arm
push-arm: build-arm
	GOARCH=arm make push

.PHONY: push-amd64
push-amd64:  build-amd64
	GOARCH=amd64 make push

.PHONY: push-all
push-all: push-arm64 push-arm push-amd64


run-dev:
	docker-compose up --build