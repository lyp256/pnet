package pnet

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseFileString(t *testing.T) {
	l, err := net.Listen("tcp", ":8800")
	assert.NoError(t, err)
	defer l.Close()
	addr := l.Addr()
	file, err := l.(*net.TCPListener).File()
	assert.NoError(t, err)
	item := fmt.Sprintf(envItemFormat, addr.Network(), addr.String(), file.Fd())
	lf, err := parseFileString(item)
	assert.NoError(t, err)
	defer lf.file.Close()
	assert.Equal(t, lf.network, addr.Network())
	assert.Equal(t, lf.address, addr.String())

}
