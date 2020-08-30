package main

import (
	"log"
	"net/http"
)

func main() {
	getConfig()
	initRMQ()
	router := NewRouter()
	log.Println("Starting Service on port :8080 ")
	log.Fatal(http.ListenAndServe(":8080", router))
}
