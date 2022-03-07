package tun

import (
	"fmt"
	"io"
	"net"
	"strings"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/protocol"
	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/adapters/gonet"
	"gvisor.dev/gvisor/pkg/tcpip/link/fdbased"
	"gvisor.dev/gvisor/pkg/tcpip/link/rawfile"
	"gvisor.dev/gvisor/pkg/tcpip/link/tun"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv4"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv6"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
	"gvisor.dev/gvisor/pkg/tcpip/transport/tcp"
	"gvisor.dev/gvisor/pkg/waiter"
)

var (
	log = logging.Logger("p2p-tun")
)

const ProtocolID = protocol.ID("/tun")

func NewTun(tun_name string) (fd int, mtu uint32, err error) {
	fd, err = tun.Open(tun_name)
	if err != nil {
		return
	}

	mtu, err = rawfile.GetMTU(tun_name)
	if err != nil {
		return
	}

	return
}

func NewNetstack(fd int, mtu uint32, tcp_stream_handler func(*net.TCPAddr, io.ReadWriteCloser)) error {

	linkEP, err := fdbased.New(&fdbased.Options{
		FDs: []int{fd},
		MTU: mtu,
	})
	if err != nil {
		return err
	}

	s := stack.New(stack.Options{
		NetworkProtocols:   []stack.NetworkProtocolFactory{ipv4.NewProtocol, ipv6.NewProtocol},
		TransportProtocols: []stack.TransportProtocolFactory{tcp.NewProtocol}})

	if err := s.CreateNIC(1, linkEP); err != nil {
		return fmt.Errorf("create nic fail %v", err)
	}

	s.SetNICForwarding(1, ipv4.ProtocolNumber, true)
	s.SetPromiscuousMode(1, true)
	s.SetSpoofing(1, true)

	subnet, _ := tcpip.NewSubnet(tcpip.Address(strings.Repeat("\x00", 4)),
		tcpip.AddressMask(strings.Repeat("\x00", 4)))
	subnet6, _ := tcpip.NewSubnet(tcpip.Address(strings.Repeat("\x00", 16)),
		tcpip.AddressMask(strings.Repeat("\x00", 16)))

	s.SetRouteTable([]tcpip.Route{
		{
			Destination: subnet,
			NIC:         1,
		},
		{
			Destination: subnet6,
			NIC:         1,
		},
	})

	tcpFwd := tcp.NewForwarder(s, 0, 256, func(r *tcp.ForwarderRequest) {
		var wq waiter.Queue
		ep, err := r.CreateEndpoint(&wq)
		if err != nil {
			log.Error(err)
			r.Complete(true)
			return
		}

		r.Complete(false)

		conn := gonet.NewTCPConn(&wq, ep)
		defer conn.Close()

		tcp_stream_handler(conn.LocalAddr().(*net.TCPAddr), conn)
	})

	s.SetTransportProtocolHandler(tcp.ProtocolNumber, tcpFwd.HandlePacket)

	return nil
}
