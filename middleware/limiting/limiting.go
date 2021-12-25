package limiting

import (
	I "DIEM-API/rpcserver"
	T "DIEM-API/tools"

	"github.com/gin-gonic/gin"
)

var RalPool *I.Pool

func InitRateLimit(type_ string, addr string, poolSize int) {
	if type_ == "uds" {
		RalPool = I.NewPool(poolSize, addr, I.DialUDS)
	} else {
		RalPool = I.NewPool(poolSize, addr, I.DialTCP)
	}
}

// request redis's throttle module for limit-rating info.
func check(xff string) []interface{} {
	return I.Choke(xff, 10, 0.1, RalPool)
}

// check if current request is valid
func Limiting(c *gin.Context) {
	xff := c.GetHeader("X-Forwarded-For")
	ret := check(xff)
	c.Header("X-RateLimit-Limit", T.Str(ret[1]))
	c.Header("X-RateLimit-Remaining", T.Str(ret[2]))
	c.Header("X-RateLimit-Next", T.Str(ret[3]))
	// `0`: current request check passed.
	if T.Str(ret[0]) != "1" {
		c.String(200, "Sorry,  Your IP requests is too frequently.")
		c.Abort()
	}
}
