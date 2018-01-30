package main

import (
	"log"
	"net/http"
)

func main() {
	router := NewRouter()

    log.Println("listening in 8080 port.")

	log.Fatal(http.ListenAndServe(":8080", router))
}
