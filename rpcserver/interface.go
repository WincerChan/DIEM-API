package rpcserver

func Choke(key string, total int, speed float64, p *Pool) []interface{} {
	rpc := new(RPCEncode)
	rpc.encodeAtom("choke")
	rpc.encodeString(key)
	rpc.encodeInteger(total)
	rpc.encodeFloat(speed)
	conn := p.Get()
	defer p.Put(conn)
	return rpc.execute(conn)
}

func Search(pages, ranges []int, terms, q []string, p *Pool) []interface{} {
	rpc := new(RPCEncode)
	rpc.encodeIntegerList(pages)
	rpc.encodeIntegerList(ranges)
	rpc.encodeStringList(terms)
	rpc.encodeStringList(q)
	conn := p.Get()
	defer p.Put(conn)
	return rpc.execute(conn)
}
