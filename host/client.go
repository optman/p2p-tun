package host

import (
	"context"
	"fmt"
	"io"
	"time"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/optman/p2p-tun/auth"
	"github.com/optman/p2p-tun/host/p2p"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/protocol"
)

type Client struct {
	h         host.Host
	ctx       context.Context
	target_id peer.ID
	auth      *auth.Authenticator
}

type Stream interface {
	io.Reader
	io.Writer
	io.Closer
	CloseWrite() error
	SetDeadline(time.Time) error
}

func NewClient(ctx context.Context, conf *NodeConfig) (*Client, error) {

	h, err := p2p.NewClientNode(ctx, conf.PrivateKey)
	if err != nil {
		return nil, err
	}

	return &Client{
		h:    h,
		ctx:  ctx,
		auth: conf.Auth,
	}, nil
}

func (self *Client) Connect(serverAddr ma.Multiaddr) error {

	addr, id := peer.SplitAddr(serverAddr)
	self.h.Peerstore().AddAddr(id, addr, peerstore.PermanentAddrTTL)

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
			return nil
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
		return fmt.Errorf("invalid peer")
	}

	return nil
}

func (self *Client) CreateStream(proto protocol.ID) func(context.Context) (Stream, error) {
	return func(ctx context.Context) (Stream, error) {

		s, err := self.h.NewStream(ctx, self.target_id, proto)
		if err != nil {
			return nil, err
		}
		if isRelayAddress(s.Conn().RemoteMultiaddr()) {
			log.Debug("through relay")
		}

		if self.auth != nil {
			if err := self.auth.Write(s); err != nil {
				return nil, err
			}
		}

		return s, nil

	}
}
