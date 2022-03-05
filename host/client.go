package host

import (
	"context"
	"io"
	"log"
	"p2p-tun/host/p2p"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

type client struct {
	h         host.Host
	ctx       context.Context
	target_id peer.ID
}

func NewClient(ctx context.Context, port int, seed int64) (*client, error) {

	h, err := p2p.NewClientNode(ctx, port, seed)
	if err != nil {
		return nil, err
	}

	return &client{
		h:   h,
		ctx: ctx,
	}, nil
}

func (self *client) Connect(id peer.ID) {
	self.target_id = id

	ctx := self.ctx
	var s network.Stream
	var err error
	retry_wait := 10 * time.Second

	for {
		s, err = self.h.NewStream(ctx, id, protocol.TestingID)
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

}

func (self *client) CreateStream(proto protocol.ID) func(context.Context) (io.ReadWriteCloser, error) {
	return func(ctx context.Context) (io.ReadWriteCloser, error) {

		s, err := self.h.NewStream(ctx, self.target_id, proto)
		if err != nil {
			return nil, err
		}
		if isRelayAddress(s.Conn().RemoteMultiaddr()) {
			log.Println("through relay")
		}

		return s, nil

	}
}
