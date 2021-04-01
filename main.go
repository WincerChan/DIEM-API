package main

import (
	M "DIEM-API/middleware"
	V "DIEM-API/views"
	"flag"
	"os"

	F "DIEM-API/config"

	"github.com/gin-gonic/gin"
)

func getServerFromArgs() string {
	isMigrate := flag.Bool("migrate", false, "should migrate?")
	service := flag.String("view", "", "running service")
	flag.Parse()
	if *isMigrate {
		F.MigrateBolt()
		os.Exit(0)
	}
	return *service
}

func main() {
	service := getServerFromArgs()
	F.InitConfig(service)
	r := gin.New()
	// register for middlewares
	M.Register(r)
	// register for views
	V.Register(r, service)
	r.Run()
}
