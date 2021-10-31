
## 💬 go-cli-chat

Chat server and client written in Go (simple prototype). The application heavily utilizes goroutines and channels. Go makes the concurrency easy to use and I had a lot of fun during development of this simple app.

![chat-client](assets/chat.png)

### Usage

```bash
$ go get github.com/aaron-baby/go-cli-chat/...
```

Now you can run client:


```bash
$ $GOPATH/bin/chat-client
```

You can also use `make` commands:


Build and run `chat-client`:

```bash
$ make run-client
```

Build `chat-client` and put binaries into corresponding `cmd/*` dir:

```bash
$ make build
```

Install `chat-client` and put binaries into `$GOPATH/bin/` dir:

```bash
$ make install
```
