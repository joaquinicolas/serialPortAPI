package main

import (
	"log"
	"net/http"
)

func main() {
	router := CreateRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
