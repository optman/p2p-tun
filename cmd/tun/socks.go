package tun

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/armon/go-socks5"
	"github.com/optman/p2p-tun/cmd/context"
	"github.com/optman/p2p-tun/host"
)

const UDPTIMEOUT time.Duration = 30 * time.Second

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

func socksConnect2(conn io.ReadWriter, ip net.IP, port int) error {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, ip)
	binary.Write(buf, binary.BigEndian, uint16(port))
	_, err := conn.Write(buf.Bytes())
	return err
}

func handleSocks5UdpStreamFunc(ctx context.Context) func(conn host.Stream) {

	log := ctx.Logger()

	return func(src host.Stream) {

		src.SetDeadline(time.Now().Add(UDPTIMEOUT))

		head := make([]byte, 6)
		if _, err := io.ReadAtLeast(src, head, len(head)); err != nil {
			return
		}
		addr := net.UDPAddr{
			IP:   net.IP(head[:4]),
			Port: int(binary.BigEndian.Uint16(head[4:6])),
		}

		log.Debug("tcp ", addr)

		buf, err := io.ReadAll(src)
		if err != nil {
			return
		}

		dst, err := net.Dial("udp", addr.String())
		if err != nil {
			return
		}

		if _, err := dst.Write(buf); err != nil {
			return
		}

		dst.SetDeadline(time.Now().Add(UDPTIMEOUT))

		buf = make([]byte, 4096)
		var n int
		n, err = dst.Read(buf)
		if err != nil {
			return
		}

		if _, err := src.Write(buf[:n]); err != nil {
			return
		}

		src.CloseWrite()
	}

}

func handleSocks5TcpStreamFunc(ctx context.Context, svr *socks5.Server) func(conn host.Stream) {

	return func(conn host.Stream) {
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
