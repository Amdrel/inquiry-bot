GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install

DOCKERIMAGE_NAME=stickmanventures.com/inquirybot

all: build

.PHONY: build docker install clean

build:
	$(GOBUILD)

docker:
	CGO_ENABLED=0 $(GOBUILD) -a -installsuffix cgo
	docker build --no-cache -t $(DOCKERIMAGE_NAME) .

install:
	$(GOINSTALL)

clean:
	$(GOCLEAN)
