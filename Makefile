
# Image URL to use all building/pushing image targets
BIN := kubectl-mutate-ingress2httpproxy
REGISTRY ?= projects.registry.vmware.com/tanzu_migrator
IMAGE    ?= $(REGISTRY)/$(BIN)
VERSION  ?= latest

GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Build kubectl-mutate binary
build: fmt vet
	go build -o bin/$(BIN) *.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run ./main.go

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Build the docker image
docker-build:
	docker build . -t ${IMAGE}:${VERSION}

# Push the docker image
docker-push:
	docker push ${IMAGE}:${VERSION}

install: build
	cp bin/$(BIN) /usr/local/bin
	chmod +x /usr/local/bin/$(BIN)
