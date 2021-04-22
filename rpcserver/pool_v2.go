package rpcserver

import (
	T "DIEM-API/tools"
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Options struct {
	Dial     func(string) (*CompoundConn, error)
	PoolSize int
	Addr     string
}

type Pool struct {
	ops   *Options
	mutex sync.Mutex
	queue chan struct{}
	// idle connections in pool
	idleConns []*Conn
	// all connections in pool
	usedConns []*Conn
}

type Conn struct {
	netConn *CompoundConn
	reader  *bufio.Reader
	writer  *bufio.Writer
	closed  bool
}

func NewPool(poolSize int, Addr string, dial func(string) (*CompoundConn, error)) *Pool {
	p := &Pool{
		ops: &Options{
			Dial:     dial,
			PoolSize: poolSize,
			Addr:     Addr,
		},
		queue:     make(chan struct{}, poolSize),
		usedConns: make([]*Conn, 0, poolSize),
		idleConns: make([]*Conn, 0, poolSize),
	}
	return p
}

func DialUDS(addr string) (*CompoundConn, error) {
	sockFile, err := net.ResolveUnixAddr("unix", addr)
	if err != nil {
		T.CheckException(err, "resolve unix addr failed.")
	}
	conn, err := net.DialUnix("unix", nil, sockFile)
	if err != nil {
		T.CheckException(err, "dial unix addr failed.")
	}
	return &CompoundConn{UDSConn: conn}, nil
}

func DialTCP(addr string) (*CompoundConn, error) {
	sockConn, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		T.CheckException(err, "resolve tcp addr failed.")
	}
	conn, err := net.DialTCP("tcp", nil, sockConn)
	if err != nil {
		T.CheckException(err, "dial tcp addr failed.")
	}
	return &CompoundConn{TcpConn: conn}, nil
}

func (c *Conn) WriteLine(line []byte) {
	line = append(line, 10)
	_, err := c.writer.Write(line)
	if err != nil {
		c.closed = true
		log.Println(err)
	}
	c.writer.Flush()
}

func (c *Conn) ReadLine() []byte {
	k, err := c.reader.ReadBytes(10)
	if err != nil {
		c.closed = true
		log.Println(err)
	}
	return k
}

func (c *Conn) WriteOnce(line []byte) {
	_, err := c.writer.Write(line)
	if err != nil {
		c.closed = true
		log.Println(err)
	}
	c.writer.Flush()
}

func (c *Conn) ReadOnce() []byte {
	prefix, err := c.reader.Peek(4)
	if err != nil {
		c.closed = true
		log.Println(err)
		return nil
	}
	size := binary.BigEndian.Uint32(prefix)
	data := make([]byte, size+4)
	_, err = io.ReadFull(c.reader, data)
	if err != nil {
		log.Println(err)
	}
	return data[4:]
}

func (p *Pool) Get() *Conn {
	// log.Println("get a conn")
	if c := p.fillToPool(); c != nil {
		return c
	}
	return p.popIdle()
}

func (p *Pool) fillToPool() *Conn {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if len(p.usedConns) < cap(p.usedConns) && len(p.idleConns) == 0 {
		c := p.newConn()
		p.usedConns = append(p.usedConns, c)
		return c
	}
	return nil
}

func (p *Pool) Put(c *Conn) {
	p.mutex.Lock()
	closed := c.closed
	p.mutex.Unlock()
	if closed {
		log.Println("error: closed connection")
		c = p.newConn()
	}
	p.pushIdle(c)
}

func (p *Pool) pushIdle(c *Conn) {
	p.mutex.Lock()
	p.idleConns = append(p.idleConns, c)
	p.mutex.Unlock()
	p.queue <- struct{}{}
}

func (p *Pool) popIdle() *Conn {
	select {
	case <-p.queue:
		p.mutex.Lock()
		c := p.idleConns[0]
		p.idleConns = p.idleConns[1:]
		p.mutex.Unlock()
		return c
	// All Conns in the pool has been timeout
	// Create a new one and fill it in the pool
	case <-time.After(time.Second * 30):
		log.Println("All Conn Timeout. Creating a new connection")
		p.mutex.Lock()
		c := p.newConn()
		drop := p.usedConns[0]
		p.usedConns = append(p.usedConns[1:], c)
		p.mutex.Unlock()
		drop.closed = true
		drop.netConn.Close()
		return c
	}
}

func (p *Pool) newConn() *Conn {
	// log.Println("here new")
	netConn, err := p.ops.Dial(p.ops.Addr)
	if err != nil {
		log.Panicln(err)
	}
	conn := &Conn{
		netConn: netConn,
		reader:  netConn.Reader(),
		writer:  netConn.Writer(),
	}
	return conn
}
