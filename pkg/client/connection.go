package client

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
)

var (
	nc       *nats.Conn
	sub      *nats.Subscription
	username string
)

var lock = &sync.Mutex{}

func GetConn() *nats.Conn {
	var err error
	lock.Lock()
	defer lock.Unlock()

	if nc == nil {
		var urls = nats.DefaultURL
		// Connect Options.
		opts := []nats.Option{nats.Name("NATS Sample Subscriber")}
		// Connect to NATS
		nc, err = nats.Connect(urls, opts...)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nc
}

// Disconnect from chat and close
func Disconnect(g *gocui.Gui, v *gocui.View) error {
	// Unsubscribe
	sub.Unsubscribe()
	GetConn().Close()
	return gocui.ErrQuit
}

// Send message
func Send(g *gocui.Gui, v *gocui.View) error {
	conn := GetConn()
	currentTime := time.Now().Format("15:04:05")
	data := strings.TrimSuffix(v.Buffer(), "\n")
	msg := fmt.Sprintf("%s [%s] %s", username, currentTime, data)
	err := conn.Publish("msg.test", []byte(msg))
	if err != nil {
		log.Fatal(err)
	}
	g.Update(func(g *gocui.Gui) error {
		v.Clear()
		v.SetCursor(0, 0)
		v.SetOrigin(0, 0)
		return nil
	})
	return nil
}

// Connect to the server, create new reader, writer and set client name
func Connect(g *gocui.Gui, v *gocui.View) error {
	username = strings.TrimSuffix(v.Buffer(), "\n")
	// Channel Subscriber
	ch := make(chan *nats.Msg, 64)
	subj := "msg.test"
	sub, _ = GetConn().ChanSubscribe(subj, ch)

	// Some UI changes
	g.SetViewOnTop("messages")
	g.SetViewOnTop("input")
	g.SetCurrentView("input")
	// Wait for server messages in new goroutine
	messagesView, _ := g.View("messages")
	go func() {
		i := 0
		for {
			msg := <-ch

			g.Update(func(g *gocui.Gui) error {
				i += 1
				fmt.Fprintln(messagesView, formatMsg(msg, i))
				return nil
			})

		}
	}()
	return nil
}

func formatMsg(m *nats.Msg, i int) string {
	return fmt.Sprintf("[#%d] Received on [%s]: '%s'", i, m.Subject, string(m.Data))
}
