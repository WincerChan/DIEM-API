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
	configPath := flag.String("config", "config.yaml", "Config file.")
	if len(os.Args) <= 1 {
		fmt.Println("No enough arguments.")
		fmt.Println("Use: ./server [prod|dev]")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "prod":
		initLogFile()
	case "dev":
	}
	initConfig(*configPath)
	initHitokotoDB()
	initRedis()
	MakeReturnMap()

	http.HandleFunc("/hitokoto/v2/", Hitokoto)
	http.HandleFunc("/hitokoto/get", Redirect301)
	log.Println("listening in " + config.ListenPort + "port.")
	err := http.ListenAndServe(config.ListenPort, nil)
	checkErr(err)
}
