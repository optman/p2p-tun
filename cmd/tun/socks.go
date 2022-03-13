package tun

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/optman/p2p-tun/cmd/context"
	"github.com/optman/p2p-tun/host"
	"github.com/optman/p2p-tun/util"
)

const UDPTIMEOUT time.Duration = 30 * time.Second

const (
	IPv4Addr = uint8(0)
	IPv6Addr = uint8(4)
)

func socksConnect(conn io.ReadWriter, ip net.IP, port int) error {
	var addrType byte
	if len(ip) == net.IPv4len {
		addrType = IPv4Addr
	} else {
		addrType = IPv6Addr
	}

	buf := bytes.NewBuffer([]byte{addrType})
	binary.Write(buf, binary.BigEndian, ip)
	binary.Write(buf, binary.BigEndian, uint16(port))
	_, err := conn.Write(buf.Bytes())

	return err
}

func recvAddr(r io.Reader) (ip net.IP, port int, err error) {
	var addrType = []byte{0}
	if _, err = r.Read(addrType); err != nil {
		return
	}

	var addrLen int
	switch addrType[0] {
	case IPv4Addr:
		addrLen = 4
	case IPv6Addr:
		addrLen = 16
	default:
		return nil, 0, fmt.Errorf("invalid addr type")
	}
	buf := make([]byte, addrLen)
	if _, err = io.ReadAtLeast(r, buf, len(buf)); err != nil {
		return
	}
	ip = net.IP(buf)

	buf = []byte{0, 0}
	if _, err = io.ReadAtLeast(r, buf, len(buf)); err != nil {
		return
	}

	port = int(binary.BigEndian.Uint16(buf))

	return

}

func handleSocks5UdpStreamFunc(ctx context.Context) func(conn host.Stream) {

	log := ctx.Logger()

	return func(src host.Stream) {
		defer src.Close()

		src.SetDeadline(time.Now().Add(UDPTIMEOUT))

		ip, port, err := recvAddr(src)
		if err != nil {
			return
		}

		addr := net.UDPAddr{
			IP:   ip,
			Port: port,
		}

		log.Debugf("udp %s", addr)

		buf, err := io.ReadAll(src)
		if err != nil {
			return
		}

		dst, err := net.Dial(addr.Network(), addr.String())
		if err != nil {
			return
		}
		defer dst.Close()

		dst.SetDeadline(time.Now().Add(UDPTIMEOUT))

		if _, err := dst.Write(buf); err != nil {
			log.Debug(err)
			return
		}

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

func handleSocks5TcpStreamFunc(ctx context.Context) func(conn host.Stream) {

	log := ctx.Logger()

	return func(src host.Stream) {
		defer src.Close()

		ip, port, err := recvAddr(src)
		if err != nil {
			return
		}

		addr := net.TCPAddr{
			IP:   ip,
			Port: port,
		}

		log.Debugf("tcp %s", addr)

		dst, err := net.Dial(addr.Network(), addr.String())
		if err != nil {
			log.Debug(err)
			return
		}
		defer dst.Close()

		util.ConcatStream(src, dst)
	}
}
