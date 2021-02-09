package tools

import (
	Logf "DIEM-API/tools/logfactory"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
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
		ret = fmt.Sprintf("%f", arg.(float64))
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
