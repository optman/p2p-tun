package p2p

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

var (
	log = logging.Logger("p2p-tun")
)

func NewServerNode(ctx context.Context, port int, seed int64) (host.Host, error) {

	h, err := libp2p.New(
		libp2p.ListenAddrStrings(listen_addrs(port)...),
		libp2p.Identity(crypto_key(seed)),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			return dht.New(ctx, h, dht.Mode(dht.ModeClient))
		}),
		libp2p.EnableHolePunching(),
		libp2p.EnableAutoRelay(),
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

	go bootstrap(ctx, h)

	return h, nil
}

func NewClientNode(ctx context.Context, port int, seed int64) (host.Host, error) {

	h, err := libp2p.New(
		libp2p.Identity(crypto_key(seed)),
		libp2p.ListenAddrStrings(listen_addrs(port)...),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			return dht.New(ctx, h, dht.Mode(dht.ModeClient))
		}),

		libp2p.EnableHolePunching(),
	)
	if err != nil {
		return nil, err
	}

	go bootstrap(ctx, h)

	return h, nil
}

func bootstrap(ctx context.Context, h host.Host) {

	bootstrap := func() {
		w := sync.WaitGroup{}
		success := 0
		for _, p := range dht.DefaultBootstrapPeers {
			p, _ := peer.AddrInfoFromP2pAddr(p)
			w.Add(1)
			go func() {
				defer w.Done()
				if err := h.Connect(ctx, *p); err == nil {
					success += 1
				}
			}()
		}

		w.Wait()

		if success == 0 {
			log.Error("bootstrap fail")
		}
	}

	//re-bootstrap after network recovery
	for {
		if len(h.Network().Conns()) < 1 {
			bootstrap()
		}

		select {
		case <-time.After(5 * time.Minute):
		case <-ctx.Done():
			return
		}
	}
}

func listen_addrs(port int) []string {
	addrs := []string{
		"/ip4/0.0.0.0/tcp/%d",
		"/ip4/0.0.0.0/udp/%d/quic",
		"/ip6/::/tcp/%d",
		"/ip6/::/udp/%d/quic",
	}

	for i, a := range addrs {
		addrs[i] = fmt.Sprintf(a, port)
	}

	return addrs
}

func crypto_key(seed int64) crypto.PrivKey {
	r := rand.New(rand.NewSource(seed))
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, -1, r)
	if err != nil {
		panic(err)
	}
	return priv
}
