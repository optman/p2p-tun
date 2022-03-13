package tun

import (
	"fmt"
	"net"
)

func NewNetstack(fd int, mtu uint32, tcp_stream_handler func(*net.TCPAddr, Stream), udp_stream_handler func(*net.UDPAddr, Stream)) error {
	return fmt.Errorf("not support on Windows")
}
