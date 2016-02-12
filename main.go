package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
	"time"
	"os/exec"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var hub = NewHub()

func socketHandler(response http.ResponseWriter, request *http.Request) {
	websocketConnection, err := upgrader.Upgrade(response, request, nil)

	if err == nil {
		hub.Register(NewConnection(websocketConnection))
	}
}

func openBrowser()  {
	time.Sleep(100 * time.Millisecond)
	out, err := exec.Command("open", "http://localhost:4000").Output()
	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
}

func main() {
	go hub.Run()
	http.HandleFunc("/ws", socketHandler)
	staticFiles := http.FileServer(http.Dir("static"))
	http.Handle("/", staticFiles)
	fmt.Println("Starting application: http://localhost:4000")
	go openBrowser()
	http.ListenAndServe(":4000", nil)
	fmt.Println("Listening...")
}
