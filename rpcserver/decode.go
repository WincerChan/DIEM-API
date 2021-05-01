package rpcserver

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

type RPCDecode struct {
	data   []byte
	buff   bytes.Buffer
	cursor int
}

const (
	SizeBytes    = 4
	IntegerBytes = 8
	TypeByte     = 2
)

func (d *RPCDecode) getLength() (size int) {
	sizePart := d.data[d.cursor : d.cursor+IntegerBytes]
	buf := bytes.NewBuffer(sizePart)
	binary.Read(buf, binary.BigEndian, size)
	d.cursor += IntegerBytes
	return
}

func (d *RPCDecode) decodeFloat() (value float64) {
	d.cursor += TypeByte + SizeBytes
	valuePart := d.data[d.cursor : d.cursor+IntegerBytes]
	bits := binary.BigEndian.Uint64(valuePart)
	value = math.Float64frombits(bits)
	d.cursor += IntegerBytes
	return
}

func (d *RPCDecode) decodeString() (value string) {
	d.cursor += TypeByte + SizeBytes
	lengthPart := d.data[TypeByte:d.cursor]
	length := binary.BigEndian.Uint32(lengthPart)
	valuePart := d.data[d.cursor : d.cursor+int(length)]
	value = string(valuePart)
	d.cursor += int(length)
	return
}

func (d *RPCDecode) decodeInteger() (value uint64) {
	d.cursor += TypeByte + SizeBytes
	valuePart := d.data[d.cursor : d.cursor+IntegerBytes]
	value = binary.BigEndian.Uint64(valuePart)
	d.cursor += IntegerBytes
	return
}

func (d *RPCDecode) extract() []interface{} {
	results := make([]interface{}, 0)
outside:
	for d.cursor < len(d.data) {
		switch d.data[d.cursor+1] {
		case 0:
			results = append(results, d.decodeString())
		case 2:
			results = append(results, d.decodeInteger())
		case 3:
			results = append(results, d.decodeFloat())
		default:
			break outside
		}
	}
	return results
}

func main() {
	st := time.Now()
	for i := 0; i < 10000000; i++ {
		a := strings.Split("choke$$0.1$10", "$")
		strconv.Atoi(a[3])
		strconv.ParseFloat(a[2], 64)
	}
	log.Println(time.Since(st))
}
