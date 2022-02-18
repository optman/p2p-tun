package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/event"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	ma "github.com/multiformats/go-multiaddr"
)

var (
	log = logging.Logger("p2p-tun")
)

var peer_id = flag.Int64("i", 0, "id seed")
var listen_port = flag.Int("p", 0, "listen port")

func main() {
	logging.SetLogLevel("p2p-tun", "info")

	flag.Parse()

	if *peer_id == 0 {
		*peer_id = rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(*forward_addr) == 0 {

		id, err := peer.Decode(*server_id)
		if err != nil {
			panic(fmt.Errorf("invalid server id, %s", err))
		}

		h, err := client_node(ctx, id)
		if err != nil {
			return
		}

		go func() {
			var s network.Stream
			var err error
			retry_wait := 10 * time.Second

			log.Info("connecting")
			for {
				s, err = h.NewStream(ctx, id, protocol.TestingID)
				if err == nil {
					break
				}

				select {
				case <-time.After(retry_wait):
				case <-ctx.Done():
					return
				}

				if retry_wait < 5*time.Minute {
					retry_wait *= 2
				}

			}
			defer s.Close()

			s.Write([]byte("hello"))
			buf := make([]byte, 5)
			s.Read(buf)
			if string(buf) != "world" {
				panic("invalid peer")
			}

			log.Info("connected")
		}()

		run_client(ctx, func(ctx context.Context) (io.ReadWriteCloser, error) {
			s, err := h.NewStream(ctx, id, "/pfw")
			if err != nil {
				return nil, err
			}
			if isRelayAddress(s.Conn().RemoteMultiaddr()) {
				log.Warn("through relay")
			}
			return s, nil
		})
	} else {

		h, err := server_node(ctx)
		if err != nil {
			return
		}

		go func() {
			subReachability, _ := h.EventBus().Subscribe(new(event.EvtLocalReachabilityChanged))
			defer subReachability.Close()

			log.Info("wait public or relay addresses ready...")

		loop:
			for {
				if containsRelayAddr(h.Addrs()) {
					break loop
				}
				select {
				case ev, ok := <-subReachability.Out():
					if !ok {
						return
					}
					evt := ev.(event.EvtLocalReachabilityChanged)
					if evt.Reachability == network.ReachabilityPublic {
						break loop
					}

				case <-time.After(5 * time.Second):
				case <-ctx.Done():
					return
				}
			}

			log.Info("ready be connect, addrs:", h.Addrs())
		}()

		h.SetStreamHandler("/pfw", func(s network.Stream) {
			handle_stream(s)
		})

		run_server(h.ID().String())
	}
}

func listen_addrs() []string {
	addrs := []string{
		"/ip4/0.0.0.0/tcp/%d",
		"/ip4/0.0.0.0/udp/%d/quic",
		"/ip6/::/tcp/%d",
		"/ip6/::/udp/%d/quic",
	}

	for i, a := range addrs {
		addrs[i] = fmt.Sprintf(a, *listen_port)
	}

	return addrs
}

func crypto_key() crypto.PrivKey {
	r := rand.New(rand.NewSource(*peer_id))
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, -1, r)
	if err != nil {
		panic(err)
	}
	return priv
}

func client_node(ctx context.Context, id peer.ID) (host.Host, error) {
	ds := dssync.MutexWrap(ds.NewMapDatastore())

	h, err := libp2p.New(
		libp2p.Identity(crypto_key()),
		libp2p.ListenAddrStrings(listen_addrs()...),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			return dht.NewDHTClient(ctx, h, ds), nil
		}),

		libp2p.EnableHolePunching(),
	)
	if err != nil {
		panic(err)
	}

	if err := bootstrap(ctx, h); err != nil {
		h.Close()
		return nil, err
	}

	return h, nil
}

func server_node(ctx context.Context) (host.Host, error) {

	ds := dssync.MutexWrap(ds.NewMapDatastore())

	h, err := libp2p.New(
		libp2p.ListenAddrStrings(listen_addrs()...),
		libp2p.Identity(crypto_key()),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			return dht.NewDHTClient(ctx, h, ds), nil
		}),
		libp2p.EnableHolePunching(),
		libp2p.EnableAutoRelay(),
	)

	if err != nil {
		panic(err)
	}

	h.SetStreamHandler(protocol.TestingID, func(s network.Stream) {
		buf := make([]byte, 5)
		n, err := s.Read(buf)
		if err == nil && string(buf[:n]) == "hello" {
			s.Write([]byte("world"))
		}

	})

	if err = bootstrap(ctx, h); err != nil {
		h.Close()
		return nil, err
	}

	return h, nil
}

func bootstrap(ctx context.Context, h host.Host) error {
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
		return errors.New("bootstrap fail")
	}

	return nil
}

func isRelayAddress(a ma.Multiaddr) bool {
	_, err := a.ValueForProtocol(ma.P_CIRCUIT)
	return err == nil
}

func containsRelayAddr(addrs []ma.Multiaddr) bool {
	for _, addr := range addrs {
		if isRelayAddress(addr) {
			return true
		}
	}
	return false
}
