package rpcserver

import (
	"bufio"
	"io"
)

type Reader struct {
	rd   *bufio.Reader
	_buf []byte
}

func NewReader(reader io.Reader) *Reader {
	return &Reader{
		rd:   bufio.NewReader(reader),
		_buf: make([]byte, 64),
	}
}

func (r *Reader) Buffered() int {
	return r.rd.Buffered()
}

func (r *Reader) Peek(n int) ([]byte, error) {
	return r.rd.Peek(n)
}

func (r *Reader) Reset(rd io.Reader) {
	r.rd.Reset(rd)
}

func (r *Reader) ReadLine() ([]byte, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}
	return line, nil
}

func (r *Reader) readLine() ([]byte, error) {
	b, err := r.rd.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	return b, nil
}

type Writer struct {
	writer

	lenBuf []byte
	numBuf []byte
}
type writer interface {
	io.Writer
	io.ByteWriter
	// io.StringWriter
	WriteString(s string) (n int, err error)
}

func NewWriter(wr writer) *Writer {
	return &Writer{
		writer: wr,

		lenBuf: make([]byte, 64),
		numBuf: make([]byte, 64),
	}
}

func (w *Writer) WriteLine(line []byte) error {
	_, err := w.Write(line)
	return err
}
