package main

import (
	M "DIEM-API/middleware"
	"flag"
	"os"

	H "DIEM-API/models/hitokoto"

	F "DIEM-API/config"

	"github.com/gin-gonic/gin"
)

func getServerFromArgs() string {
	isMigrate := flag.Bool("migrate", false, "should migrate?")
	service := flag.String("view", "", "running service")
	config := flag.String("config", "", "config file")
	flag.Parse()
	F.InitConfig(*config)
	if *isMigrate {
		H.MigrateBolt()
		os.Exit(0)
	}
	return *service
}

func main() {
	service := getServerFromArgs()
	r := gin.New()
	// register for middlewares
	M.Register(r)
	// register for views
	F.InitService(r, service)
	r.Run()
}
