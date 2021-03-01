package tools

import (
	Logf "DIEM-API/tools/logfactory"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
)

// check exception, just log. can't crash the service
func CheckException(err error, message string) {
	if err != nil {
		Logf.Error.Error().Msg(message)
	}
}

// check error, if not mute, crash the service
func CheckFatalError(err error, mute bool) {
	if err != nil && !mute {
		panic(err)
	}
}

// cast some type to string
func Str(arg interface{}) (ret string) {
	switch arg.(type) {
	case int64:
		ret = strconv.FormatInt(arg.(int64), 10)
	case int:
		ret = strconv.Itoa(arg.(int))
	case float64:
		ret = fmt.Sprintf("%.1f", arg.(float64))
	case uint32:
		ret = strconv.Itoa(int(arg.(uint32)))
	}
	return
}

// cast some type to int
func Int(arg interface{}) (ret int) {
	switch arg.(type) {
	case string:
		ret, _ = strconv.Atoi(arg.(string))

	}
	return
}

// load json file
func LoadJSON(JSONPath string) []byte {
	jsonFile, err := os.Open(JSONPath)
	CheckFatalError(err, false)

	byteValue, err := ioutil.ReadAll(jsonFile)
	CheckFatalError(err, false)

	return byteValue
}

func Int32ToBytes(num int) []byte {
	key := make([]byte, 4)
	binary.BigEndian.PutUint32(key, uint32(num))
	return key
}

func Min(num1, num2 int) int {
	if num1 < num2 {
		return num1
	}
	return num2
}

func Max(num1, num2 int) int {
	if num1 < num2 {
		return num2
	}
	return num1
}

func ConfigAbsPath(path string) string {
	return filepath.Join(viper.GetString("work_dir"), path)
}
