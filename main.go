package main

import (
	L "DIEM-API/middleware/logger"
	C "DIEM-API/middleware/rate-limiting"
	R "DIEM-API/middleware/recovery"
	S "DIEM-API/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(L.Log)
	r.Use(R.Recover)
	r.Use(C.Limiting)
	r.GET("/hitokoto/v2/", S.Hitokoto)
	r.Run()
}
