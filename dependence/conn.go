package main

import (
	"bufio"
	"net"
	"sync/atomic"
	"time"
)

type Conn struct {
	usedAt  int64 // atomic
	netConn net.Conn

	reader     *Reader
	buffWriter *bufio.Writer
	writer     *Writer

	Inited    bool
	pooled    bool
	createdAt time.Time
}

func New(netConn net.Conn) *Conn {
	conn := &Conn{
		netConn:   netConn,
		createdAt: time.Now(),
	}
	conn.reader = NewReader(netConn)
	conn.buffWriter = bufio.NewWriter(netConn)
	conn.writer = NewWriter(conn.buffWriter)
	conn.SetUsedAt(time.Now())
	return conn
}
func (cn *Conn) UsedAt() time.Time {
	unix := atomic.LoadInt64(&cn.usedAt)
	return time.Unix(unix, 0)
}

func (cn *Conn) SetUsedAt(tm time.Time) {
	atomic.StoreInt64(&cn.usedAt, tm.Unix())
}

func (cn *Conn) SetNetConn(netConn net.Conn) {
	cn.netConn = netConn
	cn.reader.Reset(netConn)
	cn.buffWriter.Reset(netConn)
}

func (cn *Conn) Write(b []byte) (int, error) {
	return cn.netConn.Write(b)
}

func (cn *Conn) RemoteAddr() net.Addr {
	if cn.netConn != nil {
		return cn.netConn.RemoteAddr()
	}
	return nil
}

func (cn *Conn) Close() error {
	return cn.netConn.Close()
}
