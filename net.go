package pnet

import (
	"net"
)

func ListenTCP(network string, laddr *net.TCPAddr) (*net.TCPListener, error) {
	ln, err := getPoolListen(laddr)
	if err == nil {
		return ln.(*net.TCPListener), nil
	}
	l, err := net.ListenTCP(network, laddr)
	if err != nil {
		return nil, err
	}
	if err := register(laddr, l); err != nil {
		panic(err)
	}
	return l, nil
}

func ListenUnix(network string, laddr *net.UnixAddr) (*net.UnixListener, error) {
	ln, err := getPoolListen(laddr)
	if err == nil {
		return ln.(*net.UnixListener), nil
	}
	l, err := net.ListenUnix(network, laddr)
	if err != nil {
		return nil, err
	}
	if err := register(laddr, l); err != nil {
		panic(err)
	}
	return l, nil
}

func ListenUDP(network string, laddr *net.UDPAddr) (*net.UDPConn, error) {
	ln, err := getPoolListen(laddr)
	if err == nil {
		return ln.(*net.UDPConn), nil
	}
	l, err := net.ListenUDP(network, laddr)
	if err != nil {
		return nil, err
	}
	if err := register(laddr, l); err != nil {
		panic(err)
	}
	return l, nil
}

func ListenIP(network string, laddr *net.IPAddr) (*net.IPConn, error) {
	ln, err := getPoolListen(laddr)
	if err == nil {
		return ln.(*net.IPConn), nil
	}
	l, err := net.ListenIP(network, laddr)
	if err != nil {
		return nil, err
	}
	if err := register(laddr, l); err != nil {
		panic(err)
	}
	return l, nil
}

func ListenUnixgram(network string, laddr *net.UnixAddr) (*net.UnixConn, error) {
	ln, err := getPoolListen(laddr)
	if err == nil {
		return ln.(*net.UnixConn), nil
	}
	l, err := net.ListenUnixgram(network, laddr)
	if err != nil {
		return nil, err
	}
	if err := register(laddr, l); err != nil {
		panic(err)
	}
	return l, nil
}

func Listen(network string, laddr string) (net.Listener, error) {
	addr, err := ResolveAddr(network, laddr)
	ln, err := getPoolListen(addr)
	if err == nil {
		return ln.(net.Listener), nil
	}
	l, err := net.Listen(network, laddr)
	if err != nil {
		return nil, err
	}
	if err := register(addr, l); err != nil {
		panic(err)
	}
	return l, nil
}

func ListenPacket(network string, laddr string) (net.PacketConn, error) {
	addr, err := ResolveAddr(network, laddr)
	ln, err := getPoolListen(addr)
	if err == nil {
		return ln.(net.PacketConn), nil
	}
	l, err := net.ListenPacket(network, laddr)
	if err != nil {
		return nil, err
	}
	if err := register(addr, l); err != nil {
		panic(err)
	}
	return l, nil
}

func ResolveAddr(network string, laddr string) (net.Addr, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		return net.ResolveTCPAddr(network, laddr)
	case "udp", "udp4", "udp6":
		return net.ResolveUDPAddr(network, laddr)
	case "ip", "ip4", "ip6":
		return net.ResolveIPAddr(network, laddr)
	case "unix", "unixgram", "unixpacket":
		return net.ResolveUnixAddr(network, laddr)
	default:
		return nil, net.UnknownNetworkError(network)
	}
}
