package rpcserver

import (
	"encoding/binary"
	"math"
)

const (
	SizeBytes    = 4
	IntegerBytes = 8
	TypeByte     = 1
)

func decodeFloat(data *[]byte) (value float64) {
	*data = (*data)[SizeBytes:]
	valuePart := (*data)[:IntegerBytes]
	bits := binary.BigEndian.Uint64(valuePart)
	value = math.Float64frombits(bits)
	*data = (*data)[IntegerBytes:]
	return
}

func bytesToInteger(value *[]byte, bytes int) (length int) {
	lengthPart := (*value)[:bytes]
	switch bytes {
	case 4:
		length = int(binary.BigEndian.Uint32(lengthPart))
	case 8:
		length = int(binary.BigEndian.Uint64(lengthPart))
	}
	*value = (*value)[bytes:]
	return
}

func decodeString(data *[]byte) (value string) {
	length := bytesToInteger(data, SizeBytes)
	valuePart := (*data)[:length]
	value = string(valuePart)
	*data = (*data)[length:]
	return value
}

func decodeInteger(data *[]byte) (value int) {
	*data = (*data)[SizeBytes:]
	value = bytesToInteger(data, IntegerBytes)
	return
}

func extract(data *[]byte) []interface{} {
	results := make([]interface{}, 0)
outside:
	for len(*data) != 0 {
		eleType := (*data)[0]
		*data = (*data)[1:]
		switch eleType {
		case 0:
			results = append(results, decodeString(data))
		case 2:
			results = append(results, decodeInteger(data))
		case 3:
			results = append(results, decodeFloat(data))
		default:
			break outside
		}
	}
	return results
}
