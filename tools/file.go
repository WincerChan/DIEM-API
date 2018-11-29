package tools

import (
	"os"
	"path"

	"github.com/rs/zerolog/log"
)

type logFile struct {
}

var FileCreator logFile

func makeDir(ltype string, isRotate bool) string {
	filepath := path.Join(logpath, ltype)
	if !isRotate {
		os.MkdirAll(filepath, os.ModePerm)
	}
	return path.Join(filepath, ltype+".log")
}

func (f *logFile) New(ltype string, isRotate bool) *os.File {
	var newFile *os.File
	var err error

	filename := makeDir(ltype, isRotate)
	if isRotate {
		newFile, err = os.OpenFile(filename,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	} else {
		newFile, err = os.OpenFile(filename,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	if err != nil {
		log.Warn().Timestamp().Msg("Could not create log fiel.")
	}
	return newFile
}
