package main

import (
	"log"
	"net/http"
)

func main() {

	router := ServerRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
