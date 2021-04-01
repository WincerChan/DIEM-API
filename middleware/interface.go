package middleware

import (
	C "DIEM-API/middleware/limiting"
	L "DIEM-API/middleware/logger"
	R "DIEM-API/middleware/recovery"

	T "DIEM-API/tools/tomlparser"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	if T.GetBool("rate-limit.enable") {
		r.Use(C.Limiting)
	}
	r.Use(L.Log)
	r.Use(R.Recover)
}
