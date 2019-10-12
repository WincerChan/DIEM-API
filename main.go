package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Config file.")
	initConfig(*configPath)
	initHitokotoDB()
	initRedis()
	MakeReturnMap()

	http.HandleFunc("/hitokoto/v2/", Hitokoto)
	log.Println("listening in " + config.ListenPort + "port.")
	err := http.ListenAndServe(config.ListenPort, nil)
	checkErr(err)
}
