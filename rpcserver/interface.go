package rpcserver

import "bytes"

func Choke(key string, total int, speed float64, p *Pool) []interface{} {
	bf := new(bytes.Buffer)
	encodeAtom(bf, "choke")
	encodeString(bf, key)
	encodeInteger(bf, total)
	encodeFloat(bf, speed)
	conn := p.Get()
	defer p.Put(conn)
	return execute(bf, conn)
}

func Search(pages, ranges []int, terms, q []string, p *Pool) []interface{} {
	bf := new(bytes.Buffer)
	encodeIntegerList(bf, pages)
	encodeIntegerList(bf, ranges)
	encodeStringList(bf, terms)
	encodeStringList(bf, q)
	conn := p.Get()
	defer p.Put(conn)
	return execute(bf, conn)
}
