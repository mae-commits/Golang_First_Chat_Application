package main

import (
	"chatapp/internal/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := routes()
	log.Println("Starting channel listener")
	go handler.ListenToWsChannel()
	fmt.Println("Server Start Up ...... localhost:8080/login")
	_ = http.ListenAndServe(":8080", mux)
}
