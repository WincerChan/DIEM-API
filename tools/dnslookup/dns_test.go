package dnslookup

import (
	"net"
	"testing"
)

func TestResolveAddr(t *testing.T) {
	// normal
	var (
		in       = "itswincer.com:80"
		expected = "104.24.125.13:80"
	)
	out := ResolveAddr(in)
	if out != expected {
		t.Errorf("ResolveAddr(%s) = %s; expected %s", in, out, expected)
	}
}

func TestResolveOne2(t *testing.T) {
	var (
		in = "www.baidu.com"
	)
	out := ResolveOne(in)
	if net.ParseIP(out) == nil {
		t.Errorf("ResolveOne(%s) = %s; not a valid ip address", in, out)
	}
}

func TestResolveOne(t *testing.T) {
	var (
		in       = "itswincer.com"
		expected = "104.24.125.13"
	)
	out := ResolveOne(in)
	if out != expected {
		t.Errorf("ResolveAddr(%s) = %s; expected %s", in, out, expected)
	}
}
