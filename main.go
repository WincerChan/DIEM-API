package main

import (
	"flag"
	"log"
	"net/http"
)

var configPath = flag.String("config", "config.yaml", "Config file.")

func main() {
	flag.Parse()
	
	initConfig(*configPath)
	initHitokotoDB()
	initRedis()

	log.Println("listening in " + config.ListenPort + "port.")
	http.HandleFunc("/hitokoto/v2/", Hitokoto)

	err = http.ListenAndServe(config.ListenPort, nil)
	checkErr(err)
}
