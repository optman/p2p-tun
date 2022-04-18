package host

import (
	"github.com/optman/p2p-tun/auth"

	"github.com/libp2p/go-libp2p-core/crypto"
)

type NodeConfig struct {
	PrivateKey  crypto.PrivKey
	ListenAddrs []string
	Auth        *auth.Authenticator
}
