package p2p

import (
	"context"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"

	quic "github.com/optman/rndz-libp2p-quic-transport"
	tcp "github.com/optman/rndz-tcp-transport"
)

var (
	log = logging.Logger("p2p-tun")
)

func NewServerNode(ctx context.Context, listenAddrs []string, privKey crypto.PrivKey) (host.Host, error) {

	h, err := libp2p.New(
		libp2p.ListenAddrStrings(listenAddrs...),
		libp2p.Identity(privKey),
		libp2p.Transport(tcp.NewTransport),
		libp2p.Transport(quic.NewTransport),
		libp2p.DisableRelay(),
	)

	if err != nil {
		return nil, err
	}

	h.SetStreamHandler(protocol.TestingID, func(s network.Stream) {
		buf := make([]byte, 5)
		n, err := s.Read(buf)
		if err == nil && string(buf[:n]) == "hello" {
			s.Write([]byte("world"))
		}

	})

	return h, nil
}

func NewClientNode(ctx context.Context, privKey crypto.PrivKey) (host.Host, error) {
	h, err := libp2p.New(
		libp2p.Identity(privKey),
		libp2p.Transport(tcp.NewTransport),
		libp2p.Transport(quic.NewTransport),
		libp2p.DisableRelay(),
	)
	if err != nil {
		return nil, err
	}

	return h, nil
}
