package rpcserver

func Choke(key string, total int, speed float64, p *Pool) []interface{} {
	rpc := new(RPCEncode)
	rpc.encodeAtom("choke")
	rpc.encodeString(key)
	rpc.encodeInteger(uint32(total))
	rpc.encodeFloat(speed)
	conn := p.Get()
	defer p.Put(conn)
	return rpc.execute(conn)
}
