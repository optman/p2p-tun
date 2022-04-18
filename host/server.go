package host

import (
	"context"

	"github.com/optman/p2p-tun/auth"
	"github.com/optman/p2p-tun/host/p2p"

	logging "github.com/ipfs/go-log/v2"
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
	h, err := p2p.NewServerNode(ctx, conf.ListenAddrs, conf.PrivateKey)
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
	return nil
}

func (self *Server) HandleStream(proto protocol.ID, f func(Stream)) {
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
