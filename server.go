package main

import (
	handler "Golang_First_Chat_Application/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/login/create", handler.CreateHandler)
	http.HandleFunc("/login/create", handler.DeleteHandler)
	http.HandleFunc("/login/newResistration", handler.NewResistrationHandler)
	http.HandleFunc("/login/newResistrationPost", handler.NewResistrationPostHandler)
	http.HandleFunc("/chat", handler.ChatHandler)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))
	fmt.Println("Server Start Up ...... localhost:8080/login")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
