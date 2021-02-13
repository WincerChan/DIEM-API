package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Options struct {
	Dial     func(string) (*net.TCPConn, error)
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
	netConn *net.TCPConn
	reader  *bufio.Reader
	writer  *bufio.Writer
	closed  bool
}

func NewPool(poolSize int, Addr string, dial func(string) (*net.TCPConn, error)) *Pool {
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

func DialTCP(addr string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Conn) WriteLine(line []byte) {
	_, err := c.writer.Write(line)
	if err != nil {
		log.Println(err)
	}
	c.writer.Flush()
}

func (c *Conn) ReadLine() []byte {
	prefix, err := c.reader.Peek(4)
	if err != nil {
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
	if c := p.fillToPool(); c != nil {
		return c
	}
	return p.popIdle()
}

func (p *Pool) fillToPool() *Conn {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if len(p.usedConns) < cap(p.usedConns) {
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
		return
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
	netConn, err := p.ops.Dial(p.ops.Addr)
	if err != nil {
		log.Panicln(err)
	}
	conn := &Conn{
		netConn: netConn,
		reader:  bufio.NewReader(netConn),
		writer:  bufio.NewWriter(netConn),
	}
	return conn
}
