package config

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type logWriter struct {
	Stderr  io.Writer
	ErrFile io.Writer
}

var LogWriter *logWriter

func init() {
	LogWriter = new(logWriter)
	LogWriter.Stderr = zerolog.ConsoleWriter{Out: os.Stderr}
	LogWriter.ErrFile, _ = os.Create("error.log")
}
