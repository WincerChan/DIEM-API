package logger

import (
	Logf "DIEM-API/tools/logfactory"

	"github.com/gin-gonic/gin"
)

// log every request
func Log(c *gin.Context) {
	c.Next()
	Logf.Access.Debug().
		Str("Error", c.Errors.String()).
		Int("Status", c.Writer.Status()).
		Str("Path", c.Request.URL.String()).
		Str("XFF", c.ClientIP()).
		Msg("")
}
