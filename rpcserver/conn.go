package rpcserver

import (
	"bufio"
	"net"
)

type CompoundConn struct {
	TcpConn *net.TCPConn
	UDSConn *net.UnixConn
}

func (c *CompoundConn) Reader() *bufio.Reader {
	if c.TcpConn != nil {
		return bufio.NewReader(c.TcpConn)
	}
	return bufio.NewReader(c.UDSConn)
}

func (c *CompoundConn) Writer() *bufio.Writer {
	if c.TcpConn != nil {
		return bufio.NewWriter(c.TcpConn)
	}
	return bufio.NewWriter(c.UDSConn)
}

func (c *CompoundConn) Close() {
	if c.TcpConn != nil {
		c.TcpConn.Close()
		return
	}
	c.UDSConn.Close()
}
