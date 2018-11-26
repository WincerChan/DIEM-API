package config

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type logWriter struct {
	stderr    io.Writer
	errlog    io.Writer
	accesslog io.Writer
}

var Log *logWriter

func (l *logWriter) GetWriter(w string) io.Writer {
	switch w {
	case "std":
		return l.stderr
	case "error":
		return l.errlog
	case "access":
		return l.accesslog
	}
	return nil
}

func init() {
	Log = new(logWriter)
	Log.stderr = zerolog.ConsoleWriter{Out: os.Stderr}
	Log.errlog, _ = os.Create("error.log")
	Log.accesslog, _ = os.Create("access.log")
}
