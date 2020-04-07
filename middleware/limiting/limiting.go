package limiting

import (
	C "DIEM-API/config"
	T "DIEM-API/tools"
	"github.com/gin-gonic/gin"
)

// request redis's throttle module for limit-rating info.
func check(xff string) []interface{} {
	ret, err := C.RedisCli.Do("CL.THROTTLE", xff, "35", "36", "360").Result()
	T.CheckFatalError(err, false)
	return ret.([]interface{})
}

// check if current request is valid
func Limiting(c *gin.Context) {
	xff := c.GetHeader("X-Forwarded-For")
	ret := check(xff)
	c.Header("X-RateLimit-Limit", T.Str(ret[1]))
	c.Header("X-RateLimit-Remaining", T.Str(ret[2]))
	c.Header("X-RateLimit-Reset", T.Str(ret[4]))
	// `0`: current request check passed.
	if T.Str(ret[0]) != "0" {
		c.String(200, "Sorry,  Your IP requests is too frequently.")
		c.Abort()
	}
}
