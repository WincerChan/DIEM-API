package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	initHitokotoDB()
	initLogFile()

	http.HandleFunc("/hitokoto/v2/", Hitokoto)
	http.HandleFunc("/hitokoto/get", Redirect301)
	http.Handle("/", http.FileServer(http.Dir("./template")))
	log.Println("listening in 520 port.")
	err := http.ListenAndServe(":520", nil)
	checkErr(err)
}
