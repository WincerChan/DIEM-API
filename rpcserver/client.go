package rpcserver

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"math"
	"net"
	"sync"
)

const (
	// Basic Type
	String  = 0
	Atom    = 1
	Integer = 2
	Float   = 3

	// Compound Type
	List = 16
)

var once sync.Once

var rpcConn *RPCConn
var wg sync.WaitGroup

type RPCConn struct {
	conn   net.Conn
	signal chan error
	addr   *net.TCPAddr
	reader *bufio.Reader
}

func integerToBytes(integer, bytes int) []byte {
	data := make([]byte, bytes, bytes)
	switch bytes {
	case 4:
		binary.BigEndian.PutUint32(data, uint32(integer))
	case 8:
		binary.BigEndian.PutUint64(data, uint64(integer))
	}
	return data
}

func encodeStringList(bf *bytes.Buffer, values []string) {
	bf.WriteByte(List)
	subBf := new(bytes.Buffer)
	for _, value := range values {
		encodeString(subBf, value)
	}
	subBytes := subBf.Bytes()
	bf.Write(integerToBytes(len(subBytes), SizeBytes))
	bf.Write(subBytes)
}

func encodeIntegerList(bf *bytes.Buffer, values []int) {
	bf.WriteByte(List)
	subBf := new(bytes.Buffer)
	for _, value := range values {
		encodeInteger(subBf, value)
	}
	subBytes := subBf.Bytes()
	bf.Write(integerToBytes(len(subBytes), SizeBytes))
	bf.Write(subBytes)
}

func encodeString(bf *bytes.Buffer, value string) {
	bf.WriteByte(String)
	bf.Write(integerToBytes(len(value), SizeBytes))
	bf.Write([]byte(value))
}

func encodeAtom(bf *bytes.Buffer, value string) {
	bf.WriteByte(Atom)
	bf.Write(integerToBytes(len(value), SizeBytes))
	bf.Write([]byte(value))
}

func encodeInteger(bf *bytes.Buffer, value int) {
	bf.WriteByte(Integer)
	bf.Write(integerToBytes(IntegerBytes, SizeBytes))
	bf.Write(integerToBytes(value, IntegerBytes))
}

func encodeFloat(bf *bytes.Buffer, value float64) {
	bf.WriteByte(Float)
	bits := math.Float64bits(value)
	bf.Write(integerToBytes(IntegerBytes, SizeBytes))
	bf.Write(integerToBytes(int(bits), IntegerBytes))
}

func execute(bf *bytes.Buffer, conn *Conn) []interface{} {
	conn.WriteOnce(bf.Bytes())
	body := conn.ReadOnce()
	return extract(&body)
}
