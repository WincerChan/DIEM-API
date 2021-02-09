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
	ops                *Options
	mutex              sync.Mutex
	queue              chan struct{}
	idleConnsLen       int
	headConn, lastCoon *ConnNode
}

type ConnNode struct {
	c    *Conn
	next *ConnNode
}

type Conn struct {
	netConn *net.TCPConn
	reader  *bufio.Reader
	writer  *bufio.Writer
}

func NewPool(poolSize int, Addr string, dial func(string) (*net.TCPConn, error)) *Pool {
	ops := &Options{
		Dial:     dial,
		PoolSize: poolSize,
		Addr:     Addr,
	}
	p := &Pool{ops: ops, queue: make(chan struct{}, ops.PoolSize)}
	// add conn to pool.
	p.fillUp()
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
	c.writer.Write(line)
	c.writer.Flush()
}

func (c *Conn) ReadLine() []byte {
	prefix, _ := c.reader.Peek(4)
	size := binary.BigEndian.Uint32(prefix)
	data := make([]byte, size+4)
	_, err := io.ReadFull(c.reader, data)
	if err != nil {
		log.Println(err)
	}
	log.Println(data)
	return data[4:]
}

func (p *Pool) Get() *Conn {
	c := p.popIdle()
	return c
}

func (p *Pool) Put(c *Conn) {
	p.pushIdle(c)
}

func (p *Pool) fillUp() {
	for p.idleConnsLen < p.ops.PoolSize {
		c := p.newConn()
		p.pushIdle(c)
	}
}

func (p *Pool) pushIdle(c *Conn) {
	cn := &ConnNode{c: c}
	p.mutex.Lock()
	if p.headConn == nil {
		p.headConn = cn
	} else {
		p.lastCoon.next = cn
	}
	p.lastCoon = cn
	p.idleConnsLen++
	p.queue <- struct{}{}
	p.mutex.Unlock()
}

func (p *Pool) popIdle() *Conn {
	for {
		select {
		case <-p.queue:
			p.mutex.Lock()
			cn := p.headConn
			if p.headConn == p.lastCoon {
				p.lastCoon = p.lastCoon.next
			}
			p.headConn = p.headConn.next
			p.idleConnsLen--
			p.mutex.Unlock()
			return cn.c
		// connection full
		case <-time.After(time.Second * 10):
			log.Println("connection pool full")
		}
	}
}

func (p *Pool) newConn() *Conn {
	netConn, err := p.ops.Dial(p.ops.Addr)
	if err != nil {
		log.Println(err)
	}
	conn := &Conn{
		netConn: netConn,
		reader:  bufio.NewReader(netConn),
		writer:  bufio.NewWriter(netConn),
	}
	return conn
}
