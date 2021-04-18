package rpcserver

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	Delimiter = 18
	String    = 0
	Atom      = 1
	Integer   = 2
	Float     = 3

	BasicType    = 0
	CompoundType = 1
)

var once sync.Once

var rpcConn *RPCConn
var wg sync.WaitGroup

type RPCEncode struct {
	buffer bytes.Buffer
}

type RPCConn struct {
	conn   net.Conn
	signal chan error
	addr   *net.TCPAddr
	tmp    bytes.Buffer
	reader *bufio.Reader
}

func (r *RPCEncode) putInteger(value, len int) {
	data := make([]byte, len)
	if 4 == len {
		binary.BigEndian.PutUint32(data, uint32(value))
	} else {
		binary.BigEndian.PutUint64(data, uint64(value))
	}
	r.buffer.Write(data)
}

func (r *RPCEncode) getLength(value uint32) []byte {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, value)
	return data
}

func (r *RPCEncode) encodeStringList(values []string) {
	r.buffer.WriteByte(Delimiter)
	r.buffer.WriteByte(CompoundType)
	r.buffer.WriteByte(String)
	sub := new(RPCEncode)
	for _, value := range values {
		sub.encodeString(value)
	}
	subBytes := sub.buffer.Bytes()
	r.putInteger(len(subBytes), 4)
	r.buffer.Write(subBytes)
}

func (r *RPCEncode) encodeIntegerList(values []int) {
	r.buffer.WriteByte(Delimiter)
	r.buffer.WriteByte(CompoundType)
	r.buffer.WriteByte(Integer)
	sub := new(RPCEncode)
	for _, value := range values {
		sub.encodeInteger(value)
	}
	subBytes := sub.buffer.Bytes()
	r.putInteger(len(subBytes), 4)
	r.buffer.Write(subBytes)
}
func (r *RPCEncode) encodeString(value string) {
	r.buffer.WriteByte(Delimiter)
	r.buffer.WriteByte(BasicType)
	r.buffer.WriteByte(String)
	r.putInteger(len(value), 4)
	r.buffer.Write([]byte(value))
}

func (r *RPCEncode) encodeAtom(value string) {
	r.buffer.WriteByte(Delimiter)
	r.buffer.WriteByte(BasicType)
	r.buffer.WriteByte(Atom)
	r.putInteger(len(value), 4)
	r.buffer.Write([]byte(value))
}

func (r *RPCEncode) encodeInteger(value int) {
	r.buffer.WriteByte(Delimiter)
	r.buffer.WriteByte(BasicType)
	r.buffer.WriteByte(Integer)
	r.putInteger(8, 4)
	r.putInteger(value, 8)
}

func (r *RPCEncode) encodeFloat(value float64) {
	r.buffer.WriteByte(Delimiter)
	r.buffer.WriteByte(BasicType)
	r.buffer.WriteByte(Float)
	bits := math.Float64bits(value)
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, bits)
	r.putInteger(8, 4)
	r.buffer.Write(data)
}

func (r *RPCEncode) execute(conn *Conn) []interface{} {
	r.buffer.Write([]byte("\r\n"))
	line := r.buffer.Bytes()
	size := r.getLength(uint32(len(line)))
	line = append(size, line...)
	conn.WriteOnce(line)
	body := conn.ReadOnce()
	d := &RPCDecode{data: body}
	return d.extract()
}

func main() {
	times, _ := strconv.Atoi(os.Args[1])
	p := NewPool(10, "127.0.0.1:4004", DialTCP)
	wg.Add(times)
	start := time.Now()
	for i := 0; i < times; i++ {
		k := strconv.Itoa(i % 3000)
		go Choke(k, 8, 0.1, p)
	}
	wg.Wait()
	log.Println(time.Since(start))
}
