package tomlparser

import (
	T "DIEM-API/tools"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

const filename = "diem.toml"

var Config *toml.Tree

func checkConfigExists(name string) string {
	path := os.Getenv("CONFIG")
	if path == "" {
		println("Please set environment variable `CONFIG`")
		os.Exit(1)
	}
	return filepath.Join(path, name)
}

func LoadTOML() {
	var err error
	path := checkConfigExists(filename)
	Config, err = toml.LoadFile(path)
	T.CheckFatalError(err, false)
}

func GetString(key string) string {
	value := Config.Get(key)
	return T.Str(value)
}

func GetInt(key string) int {
	value := Config.Get(key)
	return T.Int(value)
}

func GetBool(key string) bool {
	value := Config.Get(key)
	return value.(bool)
}

func ConfigAbsPath(key string) string {
	base := GetString("config_dir")
	return filepath.Join(base, GetString(key))
}
