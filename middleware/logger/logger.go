package logger

import (
	Logf "DIEM-API/tools/logfactory"
	"github.com/gin-gonic/gin"
)

// log every request
func Log(c *gin.Context) {
	c.Next()
	Logf.Access.Debug().Int("| Status", c.Writer.Status()).
		Str("| Path", c.Request.URL.Path).
		Str("| Error", c.Errors.String()).
		Msg("SUCCESS")
}
