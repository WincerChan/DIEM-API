package limiting

import (
	C "DIEM-API/config"
	"strconv"

	"github.com/gin-gonic/gin"
)

func check(xff string) (limitInfo []string) {
	ret, _ := C.RedisCli.Do("CL.THROTTLE", xff, "35", "36", "360").Result()
	limitInfo = make([]string, len(ret.([]interface{})))
	for i, arg := range ret.([]interface{}) {
		limitInfo[i] = strconv.FormatInt(arg.(int64), 10)
	}
	return
}

func Limiting(c *gin.Context) {
	xff := c.GetHeader("X-Forwarded-For")
	ret := check(xff)
	c.Header("X-RateLimit-Limit", ret[1])
	c.Header("X-RateLimit-Remaining", ret[2])
	c.Header("X-RateLimit-Reset", ret[4])
	// 0 代表通过检测
	if ret[0] != "0" {
		c.String(200, "Sorry,  Your IP requests is too frequently.")
		c.Abort()
	}
}
