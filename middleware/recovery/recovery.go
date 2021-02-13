package recovery

import (
	Logf "DIEM-API/tools/logfactory"
	"errors"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// log error message
func storeError() {
	Logf.Error.Error().Msg(string(debug.Stack()))
}

// recover from error, and save stack message to context.
func Recover(c *gin.Context) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		storeError()

		e := r.(error)
		c.Error(errors.New(e.Error()))

	}()
	c.Next()
}
