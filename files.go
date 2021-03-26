package pnet

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	envKey = "pnet_fds"
	// network://address@fd
	envItemFormat = "%s://%s@%d"
)

type listenFile struct {
	file    *os.File
	network string
	address string
}

type filer interface {
	File() (*os.File, error)
}

func (l *listenFile) String() string {
	return toListenKey(l.network, l.address)
}

func toListenKey(network, addr string) string {
	return fmt.Sprintf("%s://%s", network, addr)
}

var listenPool = make(map[string]*listenFile)

func getPoolListen(addr net.Addr) (interface{}, error) {
	key := toListenKey(addr.Network(), addr.String())
	file, ok := listenPool[key]
	if !ok {
		return nil, fmt.Errorf("file not found")
	}
	switch addr.Network() {
	case "tcp", "tcp4", "tcp6",
		"unix", "unixpacket":
		ln, err := net.FileListener(file.file)
		if err != nil {
			return nil, err
		}
		return ln, nil
	case "udp", "udp4", "udp6",
		"ip", "ip4", "ip6",
		"unixgram":
		ln, err := net.FilePacketConn(file.file)
		if err != nil {
			return nil, err
		}
		return ln, nil
	default:
		return nil, fmt.Errorf("unknown network %s", addr.Network())
	}
}

func register(addr net.Addr, i interface{}) error {
	f, ok := i.(filer)
	if !ok {
		return fmt.Errorf("not get file")
	}
	file, err := f.File()
	if err != nil {
		return err
	}
	lf := listenFile{
		file:    file,
		network: addr.Network(),
		address: addr.String(),
	}
	listenPool[lf.String()] = &lf
	return nil
}

func unregister(addr net.Addr) error {
	key := toListenKey(addr.Network(), addr.String())
	lf, ok := listenPool[key]
	if !ok {
		return fmt.Errorf("not found")
	}
	delete(listenPool, key)
	return lf.file.Close()

}

// InjectionNetFiles 将 fd 传入下一个启动的程序
func InjectionNetFiles(cmd *exec.Cmd) error {
	if cmd == nil {
		return fmt.Errorf("cmd is nil")
	}
	var envs []string
	startIndex := len(cmd.ExtraFiles) + 3 // stdin, stdout, and stderr are 0, 1,and 2
	for _, file := range listenPool {
		cmd.ExtraFiles = append(cmd.ExtraFiles, file.file)
		envs = append(envs, fmt.Sprintf("%s@%d", file.String(), startIndex))
		startIndex++
	}
	if envs != nil {
		cmd.Env = append(cmd.Env,
			fmt.Sprintf("%s=%s", envKey, strings.Join(envs, ",")),
		)
	}
	return nil
}

func loadEnvFiles() error {
	env := os.Getenv(envKey)
	if logger != nil {
		logger.Debugln("env:", env)
	}
	items := strings.Split(env, ",")
	for _, item := range items {
		l, err := parseFileString(item)
		if err != nil {
			if logger != nil {
				logger.Warn(err)
			}
			continue
		}
		listenPool[l.String()] = l
	}
	return nil
}

// parse from envItemFormat
func parseFileString(s string) (*listenFile, error) {
	index1 := strings.Index(s, "://")
	if index1 == -1 {
		return nil, fmt.Errorf("invalidate env item:%s", s)
	}
	network := s[:index1]
	index2 := strings.Index(s, "@")
	if index2 == -1 {
		return nil, fmt.Errorf("invalidate env item:%s", s)
	}
	addr := s[index1+3 : index2]
	fd, err := strconv.ParseUint(s[index2+1:], 10, 64)
	if err != nil {
		return nil, err
	}
	if logger != nil {
		logger.Debugln("parseFileString:", network, addr, fd)
	}
	file := os.NewFile(uintptr(fd), s)
	return &listenFile{
		file:    file,
		network: network,
		address: addr,
	}, nil
}

func init() {
	err := loadEnvFiles()
	if err != nil && logger != nil {
		logger.Errorf("loadEnvFiles:%v", err)
	}
}
