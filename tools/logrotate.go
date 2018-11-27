package tools

import (
	"io"
	"os"
	"time"
)

const (
	dayHour   = 24
	dayMinute = 60
	daySecond = 60

	logpath = "_log"
)

type Log struct {
	Writer io.Writer
	Ltype  string
}

func NewLogger(ltype string) *Log {
	llog := new(Log)
	llog.Writer = FileCreator.New(ltype)
	llog.Ltype = ltype
	return llog
}

func (l *Log) rename(now time.Time) {
	restHour := time.Hour * time.Duration(dayHour-now.Hour())
	restMinute := time.Minute * time.Duration(dayMinute-now.Minute())
	restSecond := time.Second * time.Duration(daySecond-now.Second())

	t := time.NewTimer(restHour + restMinute + restSecond)
	<-t.C
	os.Rename(logpath+l.Ltype+".log", logpath+l.Ltype+now.Format("2019-05-20"))
	l.Writer = FileCreator.New(l.Ltype)
}

func (l *Log) Rotate() {
	for {
		now := time.Now()
		l.rename(now)
	}
}
