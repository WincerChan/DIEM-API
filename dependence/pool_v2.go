package main

import (
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
	c    *net.TCPConn
	next *ConnNode
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

func (p *Pool) Get() *net.TCPConn {
	c := p.popIdle()
	return c
}

func (p *Pool) Put(c *net.TCPConn) {
	p.pushIdle(c)
}

func (p *Pool) fillUp() {
	for p.idleConnsLen < p.ops.PoolSize {
		c := p.newConn()
		p.pushIdle(c)
	}
}

func (p *Pool) pushIdle(c *net.TCPConn) {
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

func (p *Pool) popIdle() *net.TCPConn {
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

func (p *Pool) newConn() *net.TCPConn {
	tcpConn, err := p.ops.Dial(p.ops.Addr)
	if err != nil {
		log.Println(err)
	}
	return tcpConn
}
