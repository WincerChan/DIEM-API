package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := NewRouter()

	logPath := "hitokoto.log"
	httpPort := 520

	openLogFile(logPath)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	http.HandleFunc("/", Hitokoto)

	fmt.Printf("listening on %v\n", httpPort)
	fmt.Printf("logging to %v", logPath)

	err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("check to make sure it works")

	log.Fatal(http.ListenAndServe(":8080", router))
}
