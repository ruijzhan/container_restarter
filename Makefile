# parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=restarter

all: build

get_build: get build

build:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -o $(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

get:
	$(GOGET) github.com/docker/docker/api/types
	$(GOGET) github.com/docker/docker/client

docker_build:
	docker run --rm -i -t -v $(PWD):/v -w /v golang make get_build

docker_image:
	docker build . -t ruijzhan/container_restarter:$(TAG)
