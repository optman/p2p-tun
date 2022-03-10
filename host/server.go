package host

import (
	"context"
	"fmt"
	"io"
	"p2p-tun/auth"
	"p2p-tun/host/p2p"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/event"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
)

var (
	log = logging.Logger("p2p-tun")
)

type Server struct {
	h    host.Host
	ctx  context.Context
	auth *auth.Authenticator
}

func NewServer(ctx context.Context, conf *NodeConfig) (*Server, error) {
	h, err := p2p.NewServerNode(ctx, conf.ListenPort, conf.PrivateKey)
	if err != nil {
		return nil, err
	}

	return &Server{
		h:    h,
		ctx:  ctx,
		auth: conf.Auth,
	}, nil
}

func (self *Server) Host() host.Host {
	return self.h
}

func (self *Server) Start() error {
	subReachability, _ := self.h.EventBus().Subscribe(new(event.EvtLocalReachabilityChanged))
	defer subReachability.Close()

	log.Info("wait public or relay addresses ready...")

loop:
	for {
		if containsRelayAddr(self.h.Addrs()) {
			break loop
		}
		select {
		case ev, ok := <-subReachability.Out():
			if !ok {
				return fmt.Errorf("Unreachable!")
			}
			evt := ev.(event.EvtLocalReachabilityChanged)
			if evt.Reachability == network.ReachabilityPublic {
				break loop
			}

		case <-time.After(5 * time.Second):
		case <-self.ctx.Done():
			return nil
		}
	}

	log.Info("ready be connect, addrs:", self.h.Addrs())

	return nil
}

func (self *Server) HandleStream(proto protocol.ID, f func(io.ReadWriteCloser)) {
	self.h.SetStreamHandler(proto, func(s network.Stream) {
		if self.auth != nil {
			if ok, err := self.auth.Read(s); err != nil || !ok {
				log.Warnf("authenticate %s fail!", s.Conn().RemotePeer())
				s.Close()
				return
			}
		}
		f(s)
	})
}
