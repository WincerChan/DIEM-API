package main

import (
	C "DIEM-API/middleware/limiting"
	L "DIEM-API/middleware/logger"
	R "DIEM-API/middleware/recovery"

	V "DIEM-API/views"

	F "DIEM-API/config"

	"github.com/gin-gonic/gin"
)

func main() {
	F.InitConfig()
	r := gin.New()
	r.Use(L.Log)
	r.Use(R.Recover)
	r.Use(C.Limiting)
	r.GET("/hitokoto/v2/", V.Hitokoto)
	r.GET("/gaviews/v1/", V.GAViews)
	r.Run()
}
