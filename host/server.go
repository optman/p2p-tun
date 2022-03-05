package host

import (
	"context"
	"io"
	"log"
	"p2p-tun/host/p2p"
	"time"

	"github.com/libp2p/go-libp2p-core/event"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
)

type server struct {
	h   host.Host
	ctx context.Context
}

func NewServer(ctx context.Context, port int, seed int64) (*server, error) {
	h, err := p2p.NewServerNode(ctx, port, seed)
	if err != nil {
		return nil, err
	}

	return &server{
		h:   h,
		ctx: ctx,
	}, nil
}

func (self *server) Host() host.Host {
	return self.h
}

func (self *server) Start() {
	subReachability, _ := self.h.EventBus().Subscribe(new(event.EvtLocalReachabilityChanged))
	defer subReachability.Close()

	log.Println("wait public or relay addresses ready...")

loop:
	for {
		if containsRelayAddr(self.h.Addrs()) {
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
		case <-self.ctx.Done():
			return
		}
	}

	log.Println("ready be connect, addrs:", self.h.Addrs())
}

func (self *server) HandleStream(proto protocol.ID, f func(io.ReadWriteCloser)) {
	self.h.SetStreamHandler(proto, func(s network.Stream) {
		f(s)
	})
}
