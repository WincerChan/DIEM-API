package tools

import (
	Logf "DIEM-API/tools/logfactory"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func CheckException(err error, message string) {
	if err != nil {
		Logf.Error.Error().Msg(message)
	}
}

func CheckFatalError(err error, mute bool) {
	if err != nil && !mute {
		panic(err)
	}
}

func Str(arg interface{}) (ret string) {
	switch arg.(type) {
	case int64:
		ret = strconv.FormatInt(arg.(int64), 10)
	case int:
		ret = strconv.Itoa(arg.(int))
	case float64:
		ret = fmt.Sprintf("%f", arg.(float64))
	}
	return
}

func Int(arg interface{}) (ret int) {
	switch arg.(type) {
	case string:
		ret, _ = strconv.Atoi(arg.(string))

	}
	return
}

func LoadJSON(JSONPath string) []byte {
	jsonFile, err := os.Open(JSONPath)

	CheckFatalError(err, false)

	byteValue, err := ioutil.ReadAll(jsonFile)

	CheckFatalError(err, false)

	return byteValue
}
