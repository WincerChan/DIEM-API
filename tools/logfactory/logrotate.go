package logfactory

import (
	"DIEM-API/tools/filefactory"
	"io"
	"path"
	"time"
)

const (
	dayHour   = 23
	dayMinute = 59
	daySecond = 60

	logPath = "_log"
)

var (
	Stderr    io.Writer
	ErrLog    *logger
	AccessLog *logger
)

type logger struct {
	Writer   io.Writer
	level    string
	fullName string
}

func NewLogger(level string) *logger {
	aLogger := new(logger)
	aLogger.level = level
	aLogger.fullName = path.Join(logPath, level, level+".log")
	aLogger.Writer = filefactory.NewFile(aLogger.fullName)

	return aLogger
}

// rollover logfile everyday.
func (l *logger) doRollover(now time.Time) {
	filefactory.CopyFile(l.fullName, l.fullName+now.Format("2006-01-02"))
	l.Writer = filefactory.NewFile(l.fullName)
}

// run rotate at 00:00:00
func (l *logger) Rotate() {
	for {
		now := time.Now()
		restHour := time.Hour * time.Duration(dayHour-now.Hour())
		restMinute := time.Minute * time.Duration(dayMinute-now.Minute())
		restSecond := time.Second * time.Duration(daySecond-now.Second())
		t := time.NewTimer(restHour + restMinute + restSecond)
		<-t.C
		l.doRollover(now)
	}
}

// GetWriter of each log level.
func GetWriter(w string) io.Writer {
	switch w {
	case "std":
		return Stderr
	case "error":
		return ErrLog.Writer
	case "access":
		return AccessLog.Writer
	}
	return nil
}
