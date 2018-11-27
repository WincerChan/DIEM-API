package config

import (
	"io"
	"os"

	T "DIEM-API/tools"

	"github.com/rs/zerolog"
)

var stderr io.Writer
var errLog *T.Log
var accessLog *T.Log

func init() {
	stderr = zerolog.ConsoleWriter{Out: os.Stderr}
	errLog = T.NewLogger("error")
	accessLog = T.NewLogger("access")
	go errLog.Rotate()
	go accessLog.Rotate()
}

func GetWriter(w string) io.Writer {
	switch w {
	case "std":
		return stderr
	case "error":
		return errLog.Writer
	case "access":
		return accessLog.Writer
	}
	return nil
}
