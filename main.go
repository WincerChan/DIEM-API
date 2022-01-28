package main

import (
	M "DIEM-API/middleware"
	"flag"
	"fmt"
	"os"

	H "DIEM-API/models/hitokoto"

	F "DIEM-API/config"

	"github.com/gin-gonic/gin"
)

func getServerFromArgs() string {
	isMigrate := false
	service := ""
	config := ""
	flag.BoolVar(&isMigrate, "migrate", false, "shuld migrate?")
	flag.BoolVar(&isMigrate, "m", false, "shuld migrate?")
	flag.StringVar(&service, "service", "", "running service")
	flag.StringVar(&service, "s", "", "running service")
	flag.StringVar(&config, "config", "", "config file")
	flag.StringVar(&config, "c", "", "config file")
	flag.Parse()
	// config not provided
	if config == "" {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	F.InitConfig(config)
	if isMigrate {
		H.MigrateBolt()
		os.Exit(0)
	}
	return service
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
