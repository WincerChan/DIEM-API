package logfactory

import (
	"DIEM-API/tools/filefactory"
	"os"
	"path"

	"github.com/rs/zerolog"
)

var zlogE zerolog.Logger
var zlogA zerolog.Logger

type access struct{}
type error struct{}
type stdErr struct{}

var Access = access{}
var Error = error{}
var StdErr = stdErr{}

func InitLog() {

	zlogStderr := zerolog.ConsoleWriter{Out: os.Stderr}

	zlogE = zerolog.New(zlogStderr).With().Timestamp().Logger()
	zlogA = zerolog.New(zlogStderr).With().Timestamp().Logger() // comment is to disable access log
}

func newFactory(level, logPath string) *factory {
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
