package main

import (
	M "DIEM-API/middleware"
	V "DIEM-API/views"

	F "DIEM-API/config"

	"github.com/gin-gonic/gin"
)

func main() {
	F.InitConfig()
	r := gin.New()
	// register for middlewares
	M.Register(r)
	// register for views
	V.Register(r)
	r.Run()
}
