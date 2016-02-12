package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Connection struct {
	ws       *websocket.Conn
	messages chan []byte
}

func (connection *Connection) WriteMessages(hub *Hub) {
	for {
		message := <-connection.messages
		connection.ws.WriteMessage(websocket.TextMessage, message)
	}
}

func (connection *Connection) ReadMessages(hub *Hub) {
	for {
		messageType, message, err := connection.ws.ReadMessage()

		if err != nil {
			break
		}

		switch messageType {
		case websocket.BinaryMessage:
		case websocket.TextMessage:
			fmt.Println("Broadcasting message", string(message))
			hub.Broadcast(message)
		default:
			fmt.Println("Unknown message", messageType, message)
			break
		}
	}
	hub.unregister <- connection
}

func (connection *Connection) Close() {
	connection.ws.Close()
}

func NewConnection(ws *websocket.Conn) *Connection {
	return &Connection{
		ws: ws,
		messages: make(chan []byte, 256),
	}
}
