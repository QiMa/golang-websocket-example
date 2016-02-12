package main

import (
	"fmt"
)

type Hub struct {
	connections map[*Connection]bool
	broadcast   chan []byte
	register    chan *Connection
	unregister  chan *Connection
}

func (h *Hub) Unregister(connection *Connection) {
	h.unregister <- connection
}

func (h *Hub) Register(connection *Connection) {
	h.register <- connection
}

func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

func (hub *Hub) Run() {
	for {
		select {
		case connection := <-hub.unregister:
			fmt.Println("unregistering connection")
			if _, ok := hub.connections[connection]; ok {
				connection.Close()
				delete(hub.connections, connection)
			}
		case connection := <-hub.register:
			fmt.Println("registering connection")
			hub.connections[connection] = true
			go connection.ReadMessages(hub)
			go connection.WriteMessages(hub)
		case message := <-hub.broadcast:
			for connection := range hub.connections {
				connection.messages <- message
			}
		}
	}
}

func NewHub() (*Hub) {
	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: make(map[*Connection]bool),
	}
}
