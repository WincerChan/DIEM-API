package main

import (
	L "DIEM-API/middleware/logger"
	R "DIEM-API/middleware/recovery"

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
	r.GET("/pong", pong)
	r.GET("/pang", pang)
	r.Run()
}
