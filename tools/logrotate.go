package tools

import (
	"io"
	"path"
	"time"
)

const (
	dayHour   = 23
	dayMinute = 59
	daySecond = 60

	logpath = "_log"
)

type Logger struct {
	Writer   io.Writer
	level    string
	fullName string
}

func NewLogger(level string) *Logger {
	newlogger := new(Logger)
	newlogger.level = level
	newlogger.fullName = path.Join(logpath, level, level+".log")
	newlogger.Writer = FileCreator.New(newlogger.fullName)

	return newlogger
}

func (l *Logger) doRollover(now time.Time) {
	CopyFile(l.fullName, l.fullName+now.Format("2006-01-02"))
	l.Writer = FileCreator.New(l.fullName)
}

func (l *Logger) rename(now time.Time) {
	restHour := time.Hour * time.Duration(dayHour-now.Hour())
	restMinute := time.Minute * time.Duration(dayMinute-now.Minute())
	restSecond := time.Second * time.Duration(daySecond-now.Second())

	t := time.NewTimer(restHour + restMinute + restSecond)
	<-t.C
	l.doRollover(now)
}

func (l *Logger) Rotate() {
	for {
		now := time.Now()
		l.rename(now)
	}
}
