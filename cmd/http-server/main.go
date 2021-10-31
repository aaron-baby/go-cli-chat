package main

import (
	"errors"
	"github.com/Luqqk/go-cli-chat/pkg/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"os"
	"os/signal"
)

//
// REST
// ====
// $ curl -X POST -d '{"msg":"awesomeness"}' -H 'Content-Type: application/json' http://localhost:3000/messages
// {"Msg":"awesomeness"}
//
// $ curl localhost:3000/messages
// [{"Msg":"a [21:59:24] hi"},{"Msg":"awesomeness"},{"Msg":"a [21:59:49] pp"}]

var ReceivedMessages []*nats.Msg

func main() {
	nc := client.GetConn()

	nc.Subscribe("msg.test", func(m *nats.Msg) {
		ReceivedMessages = append(ReceivedMessages, m)
	})
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/messages", func(r chi.Router) {
		r.Get("/", ListReceivedMessages)
		r.Post("/", CreateMessage)
	})
	http.ListenAndServe(":3000", r)
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
}
func ListReceivedMessages(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, NewMessageListResponse(ReceivedMessages)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	data := &MsgRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	msg := &nats.Msg{Subject: "msg.test", Data: []byte(data.Msg.Msg)}
	conn := client.GetConn()
	err := conn.PublishMsg(msg)
	if err != nil {
		log.Fatal(err)
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewMsgResponse(msg))
}

type MsgRequest struct {
	*Msg
}

func (a *MsgRequest) Bind(r *http.Request) error {
	// a.Article is nil if no Article fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if a.Msg == nil {
		return errors.New("missing required Msg fields.")
	}

	return nil
}

type MsgResponse struct {
	Msg string
}

func NewMsgResponse(msg *nats.Msg) *MsgResponse {
	resp := &MsgResponse{Msg: string(msg.Data)}

	return resp
}

func (rd *MsgResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewMessageListResponse(msgs []*nats.Msg) []render.Renderer {
	list := []render.Renderer{}
	for _, msg := range msgs {
		list = append(list, NewMsgResponse(msg))
	}
	return list
}

type Msg struct {
	Msg string `json:"msg"`
}
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}
