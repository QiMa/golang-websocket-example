package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Connection struct {
	ws       *websocket.Conn
	messages chan []byte
	cancel   chan bool
}

func broadcastMessage(messageType int, message []byte) {
	switch messageType {
	case websocket.BinaryMessage:
	case websocket.TextMessage:
		fmt.Println("Broadcasting message", string(message))
		hub.Broadcast(message)
	default:
		fmt.Println("Unknown message", messageType, message)
	}
}

func handleConnectionError(connection *Connection, err error) {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
		log.Printf("Unexpected error from connection: %q", err)
	}

	hub.Unregister(connection);
}

func (connection *Connection) WriteMessages(hub *Hub) {
	for {
		select {
		case <-connection.cancel:
			return
		default:
			message := <-connection.messages
			err := connection.ws.WriteMessage(websocket.TextMessage, message)

			if err != nil {
				handleConnectionError(connection, err)
				return
			}
		}
	}
}

func (connection *Connection) ReadMessages(hub *Hub) {
	for {
		select {
		case <-connection.cancel:
			break
		default:
			messageType, message, err := connection.ws.ReadMessage()

			if err != nil {
				handleConnectionError(connection, err)
				return
			} else {
				broadcastMessage(messageType, message);
			}
		}
	}
}

func (connection *Connection) Close() {
	connection.ws.Close()
	connection.cancel <- true
}

func NewConnection(ws *websocket.Conn) *Connection {
	return &Connection{
		ws: ws,
		messages: make(chan []byte, 256),
		cancel: make(chan bool, 1),
	}
}
