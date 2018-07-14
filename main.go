package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	initHitokotoDB()

	http.HandleFunc("/hitokoto/v2/", Hitokoto)
	log.Println("listening in 520 port.")
	err := http.ListenAndServe(":8020", nil)
	checkErr(err)
}
