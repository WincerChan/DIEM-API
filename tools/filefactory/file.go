package filefactory

import (
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// CopyFile xxx
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

func makeFileDir(filename string) error {
	filepath := filepath.Dir(filename)
	os.MkdirAll(filepath, os.ModePerm)
	return nil
}

// NewFile xxx
func NewFile(filename string) *os.File {
	_ = makeFileDir(filename)

	newFile, err := os.OpenFile(filename,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Warn().Timestamp().
			Msg("Could not create log file.")
	}
	return newFile
}
