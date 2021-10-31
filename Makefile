GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install


install:
	$(GOINSTALL) ./...

run-client:
	$(GOBUILD) -o ./cmd/chat-client/chat-client ./cmd/chat-client/
	./cmd/chat-client/chat-client

run-http-server:
	$(GOBUILD) -o ./cmd/http-server/http-server ./cmd/http-server/
	./cmd/http-server/http-server