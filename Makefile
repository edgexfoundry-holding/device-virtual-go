.PHONY: build test clean docker

GO = CGO_ENABLED=0 GO111MODULE=on go

MICROSERVICES=cmd/device-virtual

.PHONY: $(MICROSERVICES)

DOCKERS=docker_device_virtual_go
.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION)
GIT_SHA=$(shell git rev-parse HEAD)
GOFLAGS=-ldflags "-X github.com/edgexfoundry/device-virtual.Version=$(VERSION)"

build: $(MICROSERVICES)
	$(GO) build ./...

cmd/device-virtual:
	$(GO) build $(GOFLAGS) -o $@ ./cmd

test:
	$(GO) test ./... -cover

clean:
	rm -f $(MICROSERVICES)

docker: $(DOCKERS)

docker_device_virtual_go:
	docker build \
		--label "git_sha=$(GIT_SHA)" \
		-t edgexfoundry/docker-device-virtual-go:$(GIT_SHA) \
		-t edgexfoundry/docker-device-virtual-go:$(VERSION)-dev \
		.
