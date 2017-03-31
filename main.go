package main

import (
	"log"
	"net/http"
	
	"io/ioutil"
	"os"
)

func main() {
	
	router := CreateRouter()
	log.Fatal(http.ListenAndServe(":8080", router))

}

