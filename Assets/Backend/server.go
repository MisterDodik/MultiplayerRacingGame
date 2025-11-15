package main

import (
	"log"
	"net/http"
)

func setupAPI() {
	manager := NewManager()

	//http.Handle("/", http.FileServer(http.Dir("./frontend")))

	http.HandleFunc("/ws", manager.serveWS)
	http.HandleFunc("/login", manager.login)
}
func main() {
	setupAPI()
	log.Println(http.ListenAndServe(":8080", nil))
}
