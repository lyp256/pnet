package pnet

import (
	"net"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListen(t *testing.T) {
	addr := ":10999"
	network := "tcp"
	l1, err := Listen(network, addr)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, l1)
	l2, err := Listen(network, addr)
	if !assert.NoError(t, err) {
		return
	}
	tl1, ok := l1.(*net.TCPListener)
	if !assert.True(t, ok) {
		return
	}

	tl2, ok := l2.(*net.TCPListener)
	if !assert.True(t, ok) {
		return
	}
	assert.Equal(t, tl1.Addr(), tl2.Addr())
	assert.NoError(t, l1.Close())
	l1.Close()
	_, err = net.Listen(network, addr)
	assert.Error(t, err)
	l2.Close()
	laddr, _ := ResolveAddr(network, addr)
	unregister(laddr)
	l, err := net.Listen(network, addr)
	assert.NoError(t, err)
	assert.NoError(t, l.Close())
}

func TestListenPacket(t *testing.T) {
	addr := ":10999"
	network := "udp"
	l1, err := ListenPacket(network, addr)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, l1)
	l2, err := ListenPacket(network, addr)
	if !assert.NoError(t, err) {
		return
	}
	tl1, ok := l1.(*net.UDPConn)
	if !assert.True(t, ok) {
		return
	}

	tl2, ok := l2.(*net.UDPConn)
	if !assert.True(t, ok) {
		return
	}
	assert.Equal(t, tl1.LocalAddr(), tl2.LocalAddr())
	assert.NoError(t, l1.Close())
	l1.Close()
	_, err = net.ListenPacket(network, addr)
	assert.Error(t, err)
	l2.Close()
	laddr, _ := ResolveAddr(network, addr)
	unregister(laddr)
	l, err := net.ListenPacket(network, addr)
	assert.NoError(t, err)
	assert.NoError(t, l.Close())
}

func TestListenIP(t *testing.T) {
	// todo
}

func TestListenTCP(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", ":10999")
	assert.NoError(t, err)
	l1, err := ListenTCP("tcp", addr)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, l1)
	l2, err := ListenTCP("tcp", addr)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, l1.Addr(), l2.Addr())
	assert.NoError(t, l1.Close())
	l1.Close()
	_, err = net.ListenTCP("tcp", addr)
	assert.Error(t, err)
	l2.Close()
	unregister(addr)
	l, err := net.ListenTCP("tcp", addr)
	assert.NoError(t, err)
	assert.NoError(t, l.Close())
}

func TestListenUDP(t *testing.T) {
	network := "udp"
	addr, err := net.ResolveUDPAddr(network, ":10999")
	assert.NoError(t, err)
	l1, err := ListenUDP(network, addr)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, l1)
	l2, err := ListenUDP(network, addr)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, l1.LocalAddr(), l2.LocalAddr())
	assert.NoError(t, l1.Close())
	l1.Close()
	_, err = net.ListenUDP(network, addr)
	assert.Error(t, err)
	l2.Close()
	unregister(addr)
	l, err := net.ListenUDP(network, addr)
	assert.NoError(t, err)
	assert.NoError(t, l.Close())
}

func TestListenUnix(t *testing.T) {
	network := "unix"
	laddr := path.Join(os.TempDir(), network)
	defer os.RemoveAll(laddr)
	addr, err := net.ResolveUnixAddr(network,
		laddr)
	assert.NoError(t, err)
	l1, err := ListenUnix(network, addr)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, l1)
	l2, err := ListenUnix(network, addr)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, l1.Addr(), l2.Addr())
	_, err = net.ListenUnix(network, addr)
	assert.Error(t, err)

	assert.NoError(t, l1.Close())
	assert.NoError(t, l2.Close())
	err = unregister(addr)
	l, err := net.ListenUnix(network, addr)
	assert.NoError(t, err)
	assert.NoError(t, l.Close())

}

func TestListenUnixgram(t *testing.T) {
	network := "unixgram"
	laddr := path.Join(os.TempDir(), network)

	addr, err := net.ResolveUnixAddr(network,
		laddr)
	assert.NoError(t, err)
	l1, err := ListenUnixgram(network, addr)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, l1)
	l2, err := ListenUnixgram(network, addr)
	if !assert.NoError(t, err) {
		return
	}
	_, err = net.ListenUnixgram(network, addr)
	assert.Error(t, err)

	assert.Equal(t, l1.LocalAddr(), l2.LocalAddr())
	assert.NoError(t, l1.Close())
	assert.NoError(t, l2.Close())
	assert.NoError(t, unregister(addr))
	assert.NoError(t, os.RemoveAll(laddr))
	l, err := net.ListenUnixgram(network, addr)
	assert.NoError(t, err)
	assert.NoError(t, l.Close())
	assert.NoError(t, os.RemoveAll(laddr))
}
