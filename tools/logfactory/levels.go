package logfactory

import (
	"DIEM-API/tools/filefactory"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"os"
	"path"
)

var zlogE zerolog.Logger
var zlogA zerolog.Logger

type access struct{}
type error struct{}
type stdErr struct{}

var Access = access{}
var Error = error{}
var StdErr = stdErr{}

func init() {
	ginMode := gin.Mode()

	zlogStderr := zerolog.ConsoleWriter{Out: os.Stderr}

	if ginMode == "debug" {
		zlogE = zerolog.New(zlogStderr).With().Timestamp().Logger()
		zlogA = zerolog.New(zlogStderr).With().Timestamp().Logger()
		return
	}

	zlogError := newFactory("error")
	zlogAccess := newFactory("access")

	go zlogError.rotate()
	go zlogAccess.rotate()

	zlogE = zerolog.New(zlogError.Writer).With().Timestamp().Logger()
	zlogA = zerolog.New(zlogAccess.Writer).With().Timestamp().Logger()

}

func newFactory(level string) *factory {
	aLogger := new(factory)
	aLogger.level = level
	aLogger.fullName = path.Join(logPath, level, level+".log")
	aLogger.Writer = filefactory.NewFile(aLogger.fullName)

	return aLogger
}

func (e *error) Debug() *zerolog.Event {
	return zlogE.Debug()
}

func (e *error) Error() *zerolog.Event {
	return zlogE.Error()
}

func (l *access) Debug() *zerolog.Event {
	return zlogA.Debug()
}

func (l *access) Error() *zerolog.Event {
	return zlogA.Error()
}
