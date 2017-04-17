package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joaquinicolas/Elca/router"
)

func main() {

	r := router.CreateRouter()
	fmt.Println("Listening on port 8080. http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
