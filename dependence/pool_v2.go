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
	// idle connections length
	idleConnsLen int
	// idle connection linkedlist
	headConn, lastCoon *ConnNode
	// all connections in pool
	allConns []*Conn
}

type ConnNode struct {
	c    *Conn
	next *ConnNode
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
		queue:    make(chan struct{}, poolSize),
		allConns: make([]*Conn, 0),
	}
	// add conn to pool.
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
		return
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
	if p.idleConnsLen == 0 && len(p.allConns) < p.ops.PoolSize {
		c := p.newConn()
		p.allConns = append(p.allConns, c)
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
	cn := &ConnNode{c: c}
	if p.headConn == nil {
		p.headConn = cn
	} else {
		p.lastCoon.next = cn
	}
	p.lastCoon = cn
	p.idleConnsLen++
	p.mutex.Unlock()
	p.queue <- struct{}{}
}

func (p *Pool) popIdle() *Conn {
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
	// All Conns in the pool has been timeout
	// Create a new one and fill it in the pool
	case <-time.After(time.Second * 1):
		log.Println("Creating a new connection")
		c := p.newConn()

		p.replaceFirstConn(c)
		return c
	}
}

func (p *Pool) replaceFirstConn(c *Conn) {
	p.mutex.Lock()
	drop := p.allConns[0]
	p.allConns = append(p.allConns[1:], c)
	p.mutex.Unlock()
	drop.closed = true
	drop.netConn.CloseRead()
	drop.netConn.CloseWrite()
	drop.netConn.Close()
}

func (p *Pool) newConn() *Conn {
	netConn, err := p.ops.Dial(p.ops.Addr)
	if err != nil {
		log.Println(err)
		return nil
	}
	conn := &Conn{
		netConn: netConn,
		reader:  bufio.NewReader(netConn),
		writer:  bufio.NewWriter(netConn),
	}
	return conn
}
