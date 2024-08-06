package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "<H1>Hello There</H1>\n")
}

func main() {
	
	mux := http.NewServeMux()

	mux.HandleFunc("/",home)
	
	log.Println("\x1b[38;5;196m[LOG]\x1b[0m  Server Starting at Port :8080")
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		fmt.Println("\x1b[38;5;196m[ERROR]\x1b[0m Failed to Start Server |")
	}

}
