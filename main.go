package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"fmt"
	"github.com/gorilla/websocket"
)

type client struct {
	connection	*websocket.Conn
	name				string 
}
var clients_list = make(map[client]bool)

var broadcast =make(chan string)
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
	log.Println(conn.RemoteAddr())
	checkws()
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	client_id := len(clients_list)+1
	var name_conn = fmt.Sprintf("connection no. -> %d",client_id)
	connection := client{conn,name_conn}
	clients_list[connection]=true

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
		msg:=fmt.Sprintf(`<div id="chat" hx-swap-oob="beforeend"><p> msg from %d : %s</p></div>`,client_id,MessAge.Hello)
		broadcast <- msg
//		err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
//		if err != nil {
//			log.Println(err)
//			return
//		}
	}

}

func handlemsg(){
	for {
		msg := <-broadcast
		for client := range clients_list {
			err:=client.connection.WriteMessage(websocket.TextMessage,[]byte(msg))
			if err!= nil {
				client.connection.Close()
				fmt.Println(err)
				log.Println("failed to write message")
			}
		}
	}
}

func checkws(){
		for cleint := range clients_list {
			test := true
			for ppl := range clients_list{
				if test{
					if cleint == ppl {
						test = false
					}
				} else {
					log.Println(ppl.connection.LocalAddr())
					ppl.connection.Close()
					delete(clients_list,ppl)
					log.Print("no of clients ", len(clients_list))
				}
			}
		}
	
}

func main() {
	http.HandleFunc("/",hello)
	http.HandleFunc("/ws", WebSocketHandler)
	
	go handlemsg()
	log.Print("Server starting at Port :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("error server couldnt start")
	}
}


