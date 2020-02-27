package dnslookup

import (
	T "DIEM-API/tools"
	"net"
	"strings"
)

// Resolve host use default dns lookup,
// return the first address.
func ResolveOne(hostname string) string {
	addresses, err := net.LookupHost(hostname)
	T.CheckFatalError(err, false)
	return addresses[0]
}

// Resolve host:port pair, then modify as
// ip_address:port.
func ResolveAddr(hostAndPort string) string {
	addresses := strings.Split(hostAndPort, ":")
	addresses[0] = ResolveOne(addresses[0])
	return strings.Join(addresses, ":")
}
