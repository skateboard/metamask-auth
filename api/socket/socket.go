package socket

import (
	"github.com/gofiber/websocket/v2"
	"log"
	"sync"
)

var Socket = Server{}

type Server struct {
	connections map[*Connection]bool
}

type Connection struct {
	Conn    *websocket.Conn
	Mux     sync.Mutex
	SessionToken string
}

func Initialize() {
	Socket.connections = make(map[*Connection]bool)
}

func (s *Server) Handler(c *websocket.Conn) {
	sessionToken := c.Params("sessionToken")
	for existingConn := range s.connections { // prevent multiple connections with same session token
		if existingConn.SessionToken == sessionToken {
			log.Println("multiple connection attempt, closing previous socket, session token: " + sessionToken)
			err := existingConn.Conn.Close()
			delete(s.connections, existingConn)

			if err != nil {
				log.Println("failed to close existing conn,", err)
			}
		}
	}

	conn := Connection{
		Conn: c,
		SessionToken: sessionToken,
	}
	s.connections[&conn] = true
	log.Println("new socket connection, session token: " + sessionToken)

	//to stop cpu usage
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func (c *Connection) Send(data []byte)  {
	c.Mux.Lock()

	err := c.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("failed to write message to client,", err)
	}
	c.Mux.Unlock()
}

func (s *Server) GetConnection(sessionToken string) *Connection {
	for conn := range s.connections {
		if conn.SessionToken == sessionToken {
			return conn
		}
	}

	return nil
}

func (s *Server) Emit(data []byte) {
	for conn := range s.connections {
		conn := conn
		go func() {
			conn.Mux.Lock()

			err := conn.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println("failed to write message to client,", err)
				delete(s.connections, conn)
			}
			conn.Mux.Unlock()
		}()
	}
}