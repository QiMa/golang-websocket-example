package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Connection struct {
	ws       *websocket.Conn
	messages chan []byte
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
		message := <-connection.messages

		if message == nil {
			return
		}

		err := connection.ws.WriteMessage(websocket.TextMessage, message)

		if err != nil {
			handleConnectionError(connection, err)
		}
	}
}

func (connection *Connection) ReadMessages(hub *Hub) {
	for {
		fmt.Println("Reading ... ")
		messageType, message, err := connection.ws.ReadMessage()

		if err != nil {
			handleConnectionError(connection, err)
			return
		} else {
			fmt.Println("Read." + string(message))
			broadcastMessage(messageType, message)
		}
	}
}

func (connection *Connection) Close() {
	close(connection.messages)
	connection.ws.Close()
}

func NewConnection(ws *websocket.Conn) *Connection {
	return &Connection{
		ws: ws,
		messages: make(chan []byte, 256),
	}
}
