package logfactory

import (
	"DIEM-API/tools/filefactory"
	"io"
	"time"
)

const (
	dayHour   = 23
	dayMinute = 59
	daySecond = 60

	logPath = "_log"
)

type factory struct {
	Writer   io.Writer
	level    string
	fullName string
}

// rollover logfile everyday.
func (l *factory) doRollover(now time.Time) {
	filefactory.CopyFile(l.fullName, l.fullName+now.Format("2006-01-02"))
	l.Writer = filefactory.NewFile(l.fullName)
}

// run rotate at 00:00:00
func (l *factory) rotate() {
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
