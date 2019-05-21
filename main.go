package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configPath := flag.String("config", "config.json", "MySQL config file.")
	initHitokotoDB(*configPath)
	switch os.Args[1] {
	case "prod":
		initHitokotoDB(*configPath)
		initLogFile()
	case "dev":
		initHitokotoDB(*configPath)
	}
	// initHitokotoDB()
	// initLogFile()
	MakeReturnMap()

	http.HandleFunc("/hitokoto/v2/", Hitokoto)
	http.HandleFunc("/hitokoto/get", Redirect301)
	// http.HandleFunc("/thinking/v1/", HandleThinkReq)
	http.HandleFunc("/", Home)
	log.Println("listening in 520 port.")
	err := http.ListenAndServe(":5200", nil)
	checkErr(err)
}
