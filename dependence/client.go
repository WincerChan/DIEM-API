package main

import (
	"bufio"
	"bytes"
	"context"
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
	DELIMITER = 18
	STRING    = 0
	ATOM      = 1
	INTEGER   = 2
	FLOAT     = 3
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

// func GetCONN() *RPCConn {
// 	rpcConn = new(RPCConn)
// 	rpcConn.signal = make(chan error)
// 	addr, _ := net.ResolveTCPAddr("tcp", "10.0.0.86:4004")
// 	rpcConn.addr = addr
// 	rpcConn.connect()
// 	go rpcConn.checkConnect()
// 	return rpcConn
// }

func newConn() *RPCConn {
	addr, _ := net.ResolveTCPAddr("tcp", "10.0.0.86:4004")
	rpcConn = &RPCConn{
		addr: addr,
	}
	rpcConn.connect()
	return rpcConn
}

func (r *RPCEncode) setLength(value uint32) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, value)
	r.buffer.Write(data)
}

func (r *RPCEncode) encodeString(value string) {
	r.buffer.WriteByte(DELIMITER)
	r.setLength(uint32(len(value)))
	r.buffer.WriteByte(STRING)
	r.buffer.Write([]byte(value))
}

func (r *RPCEncode) encodeAtom(value string) {
	r.buffer.WriteByte(DELIMITER)
	r.setLength(uint32(len(value)))
	r.buffer.WriteByte(ATOM)
	r.buffer.Write([]byte(value))
}

func (r *RPCEncode) encodeInteger(value uint32) {
	r.buffer.WriteByte(DELIMITER)
	r.setLength(4)
	r.buffer.WriteByte(INTEGER)
	r.setLength(value)
}

func (r *RPCEncode) encodeFloat(value float64) {
	r.buffer.WriteByte(DELIMITER)
	bits := math.Float64bits(value)
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, bits)
	r.setLength(8)
	r.buffer.WriteByte(FLOAT)
	r.buffer.Write(data)
}

func (c *RPCConn) connect() {
	// c.conn, err = net.DialTCP("tcp", nil, c.addr)
	conn, err := net.Dial("tcp", "10.0.0.86:4004")
	c.conn = conn
	c.conn.Write([]byte("fhdkfd"))
	c.reader = bufio.NewReader(c.conn)
	if err != nil {
		log.Fatal("fail to conneciton")
	}
}

// func (c *RPCConn) checkConnect() {
// 	for {
// 		select {
// 		case err := <-c.signal:
// 			fmt.Println("connection failed", err)
// 			c.connect()
// 		case <-time.After(time.Second * 10):
// 			fmt.Println("timeout, still alive")
// 			c.sendUntilSucceed(&c.tmp)
// 		}
// 	}
// }

func (c *RPCConn) sendUntilSucceed(buf *bytes.Buffer, conn net.Conn, reader *bufio.Reader) {
	// buf.WriteTo(c.conn)
	// r := *reader
	// ret, n, err := r.ReadLine()
	// ret, _, err := c.reader.ReadLine()
	// ret := make([]byte, 128)
	// log.Println(ret)
	// n, err := conn.Read(ret)
	// var ret []byte
	// var err error
	// for ret, _, err = c.reader.ReadLine(); len(ret) == 0 && err != nil; {
	// 	log.Println(ret)
	// 	c.signal <- err
	// 	// sleep 2ms
	// 	time.Sleep(time.Second * 10)
	// 	buf.WriteTo(c.conn)
	// }

	// rpcDecode := new(RPCDecode)
	// rpcDecode.data = ret
	// rpcDecode.extract()
}

func (r *RPCEncode) send(conn *Conn) {
	// c := new(RPCConn)
	// c := GetCONN()
	// var err error
	// c := new(RPCConn)
	// c.conn, err = net.Dial("tcp", "10.0.0.86:4004")
	// if err != nil {
	// 	log.Println("send error", err)
	// }
	r.buffer.Write([]byte("\r\n"))
	line := r.buffer.Bytes()
	// size := make([]byte, 4)
	conn.netConn.Write(line)
	// io.ReadFull(conn, size)
	// len := binary.BigEndian.Uint32(size)
	// body := make([]byte, len)
	// io.ReadFull(conn, body)
	_, err := conn.reader.ReadLine()
	if err != nil {
		log.Println("fhsk", err)
	}
	// c.sendUntilSucceed(&r.buffer, conn, reader)
}

func Choke(key string, total int, speed float64, p *ConnPool) {
	ctx := context.Background()
	rpc := new(RPCEncode)
	rpc.encodeAtom("choke")
	rpc.encodeString(key)
	rpc.encodeInteger(uint32(total))
	rpc.encodeFloat(speed)
	conn, err := p.Get(ctx)
	if err != nil {
		log.Println("ghfjikdjgnkh", err)
	}
	rpc.send(conn)
	log.Println("conns: len(conns)", p.idleConnsLen)
	defer func() {
		p.Put(ctx, conn)
		wg.Done()
	}()
}

func dummyDialer(context.Context) (net.Conn, error) {
	c, err := net.Dial("tcp", "10.0.0.86:4004")
	if err != nil {
		log.Println(err)
	}
	return c, nil
}

func main() {
	times, _ := strconv.Atoi(os.Args[1])
	start := time.Now()
	p := NewConnPool(&Options{
		Dialer:             dummyDialer,
		PoolSize:           20,
		PoolTimeout:        time.Hour,
		IdleTimeout:        time.Millisecond,
		IdleCheckFrequency: time.Millisecond,
	})
	wg.Add(times)
	// reader := bufio.NewReader(conn)
	for i := 0; i < times; i++ {
		go Choke("10.0.9.8", 3, 0.1, p)
	}
	wg.Wait()
	log.Println(time.Since(start))
}
