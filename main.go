package main

import (
	"log"
	"net/http"
)

func main() {
	router := NewRouter()

    log.Println("listening in 520 port.")

	log.Fatal(http.ListenAndServe(":520", router))
}
