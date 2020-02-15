package main

import (
	C "DIEM-API/middleware/limiting"
	L "DIEM-API/middleware/logger"
	R "DIEM-API/middleware/recovery"

	S "DIEM-API/services"

	F "DIEM-API/config"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(L.Log)
	r.Use(R.Recover)
	if F.Enabled {
		r.Use(C.Limiting)
	}
	r.GET("/hitokoto/v2/", S.Hitokoto)
	r.Run()
}
