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
	Ltype    string
	fullName string
}

func (l *Logger) createFile() {
	l.fullName = path.Join(logpath, l.Ltype, l.Ltype+".log")
	l.Writer = FileCreator.New(l.fullName)
}

func NewLogger(logtype string) *Logger {
	newlogger := new(Logger)
	newlogger.Ltype = logtype
	newlogger.createFile()

	return newlogger
}

func (l *Logger) doRotate(now time.Time) {
	CopyFile(l.fullName, l.fullName+now.Format("2006-01-02"))
	l.Writer = FileCreator.New(l.fullName)
}

func (l *Logger) rename(now time.Time) {
	restHour := time.Hour * time.Duration(dayHour-now.Hour())
	restMinute := time.Minute * time.Duration(dayMinute-now.Minute())
	restSecond := time.Second * time.Duration(daySecond-now.Second())

	t := time.NewTimer(restHour + restMinute + restSecond)
	<-t.C
	l.doRotate(now)
}

func (l *Logger) Rotate() {
	for {
		now := time.Now()
		l.rename(now)
	}
}
