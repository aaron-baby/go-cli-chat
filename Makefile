GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

build: 
	$(GOBUILD) -o ./cmd/chat-client/chat-client ./cmd/chat-client/

install:
	$(GOINSTALL) ./...

run-client:
	$(GOBUILD) -o ./cmd/chat-client/chat-client ./cmd/chat-client/
	./cmd/chat-client/chat-client
