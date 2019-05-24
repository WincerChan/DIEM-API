package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configPath := flag.String("config", "config.json", "MySQL config file.")
	if len(os.Args) <= 1 {
		fmt.Println("No enough arguments.")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "prod":
		initLogFile()
	case "dev":
	}
	initHitokotoDB(*configPath)
	MakeReturnMap()

	http.HandleFunc("/hitokoto/v2/", Hitokoto)
	http.HandleFunc("/hitokoto/get", Redirect301)
	// http.HandleFunc("/thinking/v1/", HandleThinkReq)
	log.Println("listening in " + config.Port + "port.")
	err := http.ListenAndServe(config.Port, nil)
	checkErr(err)
}
