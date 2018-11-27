package tools

import (
	"os"
	"path"

	"github.com/rs/zerolog/log"
)

type logFile struct {
}

var FileCreator logFile

func makeDir(ltype string) string {
	filepath := path.Join(logpath, ltype)
	os.MkdirAll(filepath, os.ModePerm)
	return path.Join(filepath, ltype+".log")
}

func (f *logFile) New(ltype string) *os.File {
	filename := makeDir(ltype)
	newFile, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Warn().Timestamp().Msg("Could not create log fiel.")
	}
	return newFile
}
