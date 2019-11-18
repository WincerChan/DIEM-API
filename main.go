package main

import (
	"flag"
	"log"
	"net/http"
)

func init() {
	configPath := flag.String(
		"config",
		"config.yaml",
		"Database Config file.")
	flag.Parse()
	initConfig(*configPath)
	initHitokotoDB()
	initRedis()
}

func main() {
	log.Println("listening in " + config.ListenPort + " port.")
	http.HandleFunc("/hitokoto/v2/", Hitokoto)

	err = http.ListenAndServe(config.ListenPort, nil)
	checkErr(err)
}
