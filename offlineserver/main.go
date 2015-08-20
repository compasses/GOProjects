package main

import (
	"log"
	"net/http"
)

func main() {

	router := ServerRouter()

	log.Fatal(http.ListenAndServe(":10080", router))
}
