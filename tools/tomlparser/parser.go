package tomlparser

import (
	T "DIEM-API/tools"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

var Config *toml.Tree

func LoadTOML(path string) {
	var err error
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
	file := GetString(key)
	if filepath.IsAbs(file) {
		return file
	}
	return filepath.Join(base, file)
}
