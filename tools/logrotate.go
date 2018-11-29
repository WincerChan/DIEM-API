package tools

import (
	"io"
	"os"
	"strings"
	"time"
)

const (
	dayHour   = 23
	dayMinute = 59
	daySecond = 60

	logpath = "_log"
)

type Logger struct {
	Writer io.Writer
	Ltype  string
}

func CopyFile(srcFile, destFile string) error {
	file, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer file.Close()
	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, file)
	return err
}

func NewLogger(ltype string) *Logger {
	llog := new(Logger)
	llog.Writer = FileCreator.New(ltype, false)
	llog.Ltype = ltype
	return llog
}

func (l *Logger) rename(now time.Time) {
	restHour := time.Hour * time.Duration(dayHour-now.Hour())
	restMinute := time.Minute * time.Duration(dayMinute-now.Minute())
	restSecond := time.Second * time.Duration(daySecond-now.Second())

	t := time.NewTimer(restHour + restMinute + restSecond)
	<-t.C
	CopyFile(strings.Join([]string{logpath, l.Ltype, l.Ltype + ".log"}, "/"),
		strings.Join([]string{logpath, l.Ltype, l.Ltype + now.Format("2006-01-02")}, "/"))
	l.Writer = FileCreator.New(l.Ltype, true)
	return
}

func (l *Logger) Rotate() {
	for {
		now := time.Now()
		l.rename(now)
	}
}
