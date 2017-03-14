package main

import (
	"log"
	"net/http"
	"github.com/joaquinicolas/Elca/libs"
	"io/ioutil"
	"os"
)

func main() {
	libs.NewLog(ioutil.Discard,os.Stdout,os.Stdout,os.Stderr)
	router := CreateRouter()
	log.Fatal(http.ListenAndServe(":8080", router))

}

