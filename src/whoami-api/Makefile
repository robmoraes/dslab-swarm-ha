IMAGE ?= robmoraes/dslab-whoami
VERSION ?= 0.1.0
TAG ?= $(VERSION)
PLATFORM ?= linux/amd64
APP_BIN ?= bin/dslab-whoami
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo local)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)

.PHONY: help tidy test build run docker-build docker-run push release clean

help:
	@echo "Targets:"
	@echo "  make test          Run Go tests"
	@echo "  make build         Build local binary"
	@echo "  make run           Run API on :8080"
	@echo "  make docker-build  Build Docker image $(IMAGE):$(TAG)"
	@echo "  make docker-run    Run Docker image on localhost:8080"
	@echo "  make push          Push $(IMAGE):$(TAG) and latest"
	@echo "  make release       Build and push image"

tidy:
	go mod tidy

test:
	go test ./...

build:
	mkdir -p $(dir $(APP_BIN))
	CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o $(APP_BIN) ./cmd/api

run:
	go run -ldflags="$(LDFLAGS)" ./cmd/api

docker-build:
	docker build \
		--platform $(PLATFORM) \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(IMAGE):$(TAG) \
		-t $(IMAGE):latest \
		.

docker-run:
	docker run --rm -p 8080:8080 $(IMAGE):$(TAG)

push:
	docker push $(IMAGE):$(TAG)
	docker push $(IMAGE):latest

release: docker-build push

clean:
	rm -rf bin
