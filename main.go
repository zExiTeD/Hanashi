package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"fmt"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type hellos struct {
	Hello string	`json:"text"`
}

func hello (w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.templ"))
	data := hellos{
		Hello : "hii",
	} 
	tmpl.Execute(w,data)
}


func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	fmt.Println(conn.LocalAddr())
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	log.Println("Client connected")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
	
		var MessAge hellos
		err = json.Unmarshal(message,&MessAge)

		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Received message:", MessAge.Hello, messageType)
		msg:=fmt.Sprintf(`<div id="chat" hx-swap-oob="beforeend"><p>message recived %s</p></div>`,MessAge.Hello)
		err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println(err)
			return
		}
	}

}

func main() {
	http.HandleFunc("/",hello)
	http.HandleFunc("/ws", WebSocketHandler)

	log.Print("Server starting at Port :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("error server couldnt start")
	}
}


