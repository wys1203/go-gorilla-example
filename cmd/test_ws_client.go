//go:build test_ws_client
// +build test_ws_client

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	WebSocketClientTest()
}

func WebSocketClientTest() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer c.Close()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Failed to read message:", err)
				return
			}

			var msg struct {
				Type string `json:"type"`
				Acct string `json:"acct"`
			}

			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.Println("Failed to unmarshal message:", err)
				continue
			}

			switch msg.Type {
			case "failed_sign_in":
				log.Printf("User with acct %s failed to sign in", msg.Acct)
			default:
				log.Printf("Unknown message type: %s", msg.Type)
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	http.Post("http://localhost:8080/signin", "application/json", strings.NewReader(`{"acct": "wys", "password": "test"}`))

	for {
		select {
		case <-ticker.C:
			log.Println("No notifications received for 10 seconds")
			return
		case <-interrupt:
			log.Println("Interrupt received, closing connection")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Failed to write close message:", err)
				return
			}
			return
		}
	}

}
