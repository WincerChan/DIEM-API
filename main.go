package main

import (
	M "DIEM-API/middleware"
	V "DIEM-API/views"
	"os"

	F "DIEM-API/config"

	"github.com/gin-gonic/gin"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		F.MigrateBolt()
		os.Exit(0)
	}
	F.InitConfig()
	r := gin.New()
	// register for middlewares
	M.Register(r)
	// register for views
	V.Register(r)
	r.Run()
}
