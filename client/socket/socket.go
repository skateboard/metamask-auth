package socket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

type Socket struct {
	SessionToken 		string

	MessageHandler func(message Message)

	connection      	*websocket.Conn
	isPinging 			bool
}

func (socket *Socket) Connect() error {
	ws := url.URL{Scheme: "ws", Host: "127.0.0.1", Path: "/v1/ws/" + socket.SessionToken}
	connection, _, err := websocket.DefaultDialer.Dial(ws.String(), nil)
	if err != nil {
		log.Println("failed to connect to socket:", err)
		return err
	}
	socket.connection = connection

	go socket.handleMessages()
	if !socket.isPinging {
		go socket.startPingInterval()
	}

	return nil
}

func (socket *Socket) startPingInterval() {
	socket.isPinging = true
	for {
		err := socket.ping()
		if err != nil {
			log.Println("failed to ping socket:", err)
		}

		time.Sleep(10 * time.Second)
	}
}

func (socket *Socket) ping() error {
	err := socket.connection.WriteMessage(websocket.TextMessage, []byte("PING_SOCKET"))
	if err != nil {
		return err
	}

	return nil
}

func (socket *Socket) reconnect() {
	log.Println("reconnecting to socket")

	if err := socket.Connect(); err != nil {
		log.Println("failed to reconnect to socket:", err)
		time.Sleep(5 * time.Second)
		socket.reconnect()
	}
}

func (socket *Socket) parseMessage(data []byte) Message {
	var message Message
	err := json.Unmarshal(data, &message)
	if err != nil {
		log.Println("failed to parse message:", err)

		return Message{}
	}

	return message
}

func (socket *Socket) handleMessages() {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := socket.connection.ReadMessage()
			if err != nil {
				log.Println("failed to read message from socket:", err)
				socket.reconnect()
				return
			}

			parsedMessage := socket.parseMessage(message)
			if parsedMessage.Type == "" {
				continue
			}

			go socket.MessageHandler(parsedMessage)
		}
	}()
}