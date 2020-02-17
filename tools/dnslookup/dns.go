package dnslookup

import (
	"net"
	"strings"
)

func ResolveOne(hostname string) string {
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		panic(err)
	}
	return addrs[0]
}

func ResolveAddr(hostAndPort string) string {
	addrs := strings.Split(hostAndPort, ":")
	addrs[0] = ResolveOne(addrs[0])
	return strings.Join(addrs, ":")
}
