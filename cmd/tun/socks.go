package tun

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/armon/go-socks5"
)

func socksConnect(conn io.ReadWriter, target_addr *net.TCPAddr) error {
	conn.Write([]byte{0x05, 0x01, 0x00})
	resp := make([]byte, 2)
	if _, err := io.ReadFull(conn, resp); err != nil {
		return fmt.Errorf("socks auth fail")
	}

	var addrType byte
	var respLen int
	if len(target_addr.IP) == net.IPv4len /*ipv4*/ {
		addrType = 1
		respLen = 10
	} else /*ipv6*/ {
		addrType = 4
		respLen = 22
	}

	buf := bytes.NewBuffer([]byte{0x05, 0x01, 0x00, addrType})
	binary.Write(buf, binary.BigEndian, target_addr.IP)
	binary.Write(buf, binary.BigEndian, uint16(target_addr.Port))
	conn.Write(buf.Bytes())
	resp = make([]byte, respLen)
	if _, err := io.ReadFull(conn, resp); err != nil {
		return fmt.Errorf("socks connect fail")
	}

	return nil
}

func handleSocks5StreamFunc(ctx context.Context, svr *socks5.Server) func(conn io.ReadWriteCloser) {

	return func(conn io.ReadWriteCloser) {
		svr.ServeConn(&fake_conn{conn})
	}
}

func newSocks5Server() (*socks5.Server, error) {
	conf := &socks5.Config{}
	return socks5.New(conf)
}

type fake_conn struct {
	io.ReadWriteCloser
}

func (c *fake_conn) LocalAddr() net.Addr                { return nil }
func (c *fake_conn) RemoteAddr() net.Addr               { return nil }
func (c *fake_conn) SetDeadline(t time.Time) error      { return nil }
func (c *fake_conn) SetReadDeadline(t time.Time) error  { return nil }
func (c *fake_conn) SetWriteDeadline(t time.Time) error { return nil }
