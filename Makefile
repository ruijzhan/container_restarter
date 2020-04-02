# parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=restarter

all: build

build:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -o $(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

docker_build:
	docker run --rm -i -t -v $(PWD):/v -w /v golang make get_build

docker_image:
	docker build . -t ruijzhan/container_restarter:$(TAG)
