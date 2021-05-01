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
	String  = 0
	Atom    = 1
	Integer = 2
	Float   = 3

	BasicType    = 0
	CompoundType = 1
)

var once sync.Once

var rpcConn *RPCConn
var wg sync.WaitGroup

type RPCConn struct {
	conn   net.Conn
	signal chan error
	addr   *net.TCPAddr
	tmp    bytes.Buffer
	reader *bufio.Reader
}

func encodeStringList(bf *bytes.Buffer, values []string) {
	bf.WriteByte(16)
	subBf := new(bytes.Buffer)
	for _, value := range values {
		encodeString(subBf, value)
	}
	subBytes := subBf.Bytes()
	bf.Write(integerToBytes(len(subBytes), 4))
	bf.Write(subBytes)
}

func encodeIntegerList(bf *bytes.Buffer, values []int) {
	bf.WriteByte(16)
	subBf := new(bytes.Buffer)
	for _, value := range values {
		encodeInteger(subBf, value)
	}
	subBytes := subBf.Bytes()
	bf.Write(integerToBytes(len(subBytes), 4))
	bf.Write(subBytes)
}

func encodeString(bf *bytes.Buffer, value string) {
	bf.WriteByte(String)
	bf.Write(integerToBytes(len(value), 4))
	bf.Write([]byte(value))
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

func encodeAtom(bf *bytes.Buffer, value string) {
	bf.WriteByte(Atom)
	bf.Write(integerToBytes(len(value), 4))
	bf.Write([]byte(value))
}

func encodeInteger(bf *bytes.Buffer, value int) {
	bf.WriteByte(Integer)
	bf.Write(integerToBytes(8, 4))
	bf.Write(integerToBytes(value, 8))
}

func encodeFloat(bf *bytes.Buffer, value float64) {
	bf.WriteByte(Float)
	bits := math.Float64bits(value)
	bf.Write(integerToBytes(8, 4))
	bf.Write(integerToBytes(int(bits), 8))
}

func execute(bf *bytes.Buffer, conn *Conn) []interface{} {
	conn.WriteOnce(bf.Bytes())
	body := conn.ReadOnce()
	d := &RPCDecode{data: body}
	return d.extract()
}

// func main() {
// 	times, _ := strconv.Atoi(os.Args[1])
// 	p := NewPool(10, "127.0.0.1:4004", DialTCP)
// 	wg.Add(times)
// 	start := time.Now()
// 	for i := 0; i < times; i++ {
// 		k := strconv.Itoa(i % 3000)
// 		go Choke(k, 8, 0.1, p)
// 	}
// 	wg.Wait()
// 	log.Println(time.Since(start))
// }
