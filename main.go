package main

import (
	L "DIEM-API/middleware/logger"
	C "DIEM-API/middleware/rate-limiting"
	R "DIEM-API/middleware/recovery"
	S "DIEM-API/services"

	"github.com/gin-gonic/gin"
)

func pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func pang(c *gin.Context) {
	k := make([]int, 0)
	_ = k[0]
	c.JSON(200, gin.H{
		"message": "pang",
	})
}

func main() {
	r := gin.New()
	r.Use(L.Log)
	r.Use(R.Recover)
	r.Use(C.Limiting)
	r.GET("/pong", pong)
	r.GET("/pang", pang)
	r.GET("/hitokoto/v2/", S.Hitokoto)
	r.Run()
}
