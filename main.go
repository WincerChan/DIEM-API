package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	initHitokotoDB()
	// initLogFile()
	MakeReturnMap()

	http.HandleFunc("/hitokoto/v2/", Hitokoto)
	http.HandleFunc("/hitokoto/get", Redirect301)
	http.HandleFunc("/thinking/v1/", HandleThinkReq)
	http.HandleFunc("/", Home)
	log.Println("listening in 520 port.")
	err := http.ListenAndServe(":5200", nil)
	checkErr(err)
}
