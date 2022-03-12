package tun

import (
	"fmt"
	"io"
	"net"
)

func NewNetstack(fd int, mtu uint32, tcp_stream_handler func(*net.TCPAddr, io.ReadWriteCloser), udp_stream_handler func(*net.UDPAddr, io.ReadWriteCloser)) error {
	return fmt.Errorf("not support on Windows")
}
