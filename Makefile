GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install

DOCKERIMAGE_NAME=stickmanventures.com/inquirybot

all: inquiry-bot

inquiry-bot:
	$(GOBUILD)

.PHONY: docker install clean

alpine:
	CGO_ENABLED=0 $(GOBUILD) -a -installsuffix cgo

docker: alpine
	docker build -t $(DOCKERIMAGE_NAME) .

install: inquiry-bot
	$(GOINSTALL)

clean:
	$(GOCLEAN)
