# parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=restarter

all: build

test:
	$(GOTEST)

build:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -ldflags "-s -w" -o $(BINARY_NAME) -v

get:
	$(GOGET) -u ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
