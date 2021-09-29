APP_REPO ?= quay.io/synpse/metrics-nats-example-app
TAG_NAME = $(shell git rev-parse --short=7 HEAD)$(shell [[ $$(git status --porcelain) = "" ]] || echo -dirty)
OUTPUT_BIN_APP ?= release
LDFLAGS		+= -s -w
GOARCH ?= arm64


APP_REPO ?= quay.io/synpse/aws-iot-core-example

.PHONY: image
image:
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t ${APP_REPO}:latest --push -f Dockerfile .

.PHONY: app-arm, app-arm64, app-amd64
app-arm:
	CGO_ENABLED=0 GOARCH=arm go build -ldflags "$(LDFLAGS)" -o ${OUTPUT_BIN_APP}/linux/arm/v7/app ./cmd/app

app-arm64:
	CGO_ENABLED=0 GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o ${OUTPUT_BIN_APP}/linux/arm64/app ./cmd/app

app-amd64:
	CGO_ENABLED=0 GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ${OUTPUT_BIN_APP}/linux/amd64/app ./cmd/app

app: app-arm app-arm64 app-amd64

run-dev:
	docker-compose up --build