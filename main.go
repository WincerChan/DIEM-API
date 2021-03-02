package main

import (
	L "DIEM-API/middleware/logger"
	R "DIEM-API/middleware/recovery"

	V "DIEM-API/views"

	F "DIEM-API/config"

	"github.com/gin-gonic/gin"
)

func main() {
	F.InitConfig()
	r := gin.New()
	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("daterange", B.Daterange)
	// }
	r.Use(L.Log)
	r.Use(R.Recover)
	// r.Use(C.Limiting)
	r.GET("/hitokoto/v2/", V.Hitokoto)
	r.GET("/gaviews/v1/", V.GAViews)
	r.GET("/blog-search/v1/", V.BlogSearchViews)
	r.Run()
}
