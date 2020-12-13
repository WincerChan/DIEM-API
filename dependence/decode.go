package main

import (
	"bytes"
	"encoding/binary"
	"math"
)

type RPCDecode struct {
	data   []byte
	buff   bytes.Buffer
	cursor int
}

func (d *RPCDecode) getLength() (size int) {
	sizePart := d.data[d.cursor : d.cursor+4]
	buf := bytes.NewBuffer(sizePart)
	binary.Read(buf, binary.BigEndian, size)
	d.cursor += 4
	return
}

func (d *RPCDecode) decodeFloat() (value float64) {
	d.cursor += 1 + 4
	valuePart := d.data[d.cursor : d.cursor+8]
	bits := binary.BigEndian.Uint64(valuePart)
	value = math.Float64frombits(bits)
	d.cursor += 8
	return
}

func (d *RPCDecode) decodeInteger() (value uint32) {
	d.cursor += 1 + 4
	valuePart := d.data[d.cursor : d.cursor+4]
	value = binary.BigEndian.Uint32(valuePart)
	d.cursor += 4
	return
}

func (d *RPCDecode) extract() []interface{} {
	results := make([]interface{}, 0)
outside:
	for d.cursor < len(d.data) {
		switch d.data[d.cursor] {
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
