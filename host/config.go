package host

import (
	"p2p-tun/auth"

	"github.com/libp2p/go-libp2p-core/crypto"
)

type NodeConfig struct {
	PrivateKey crypto.PrivKey
	ListenPort int
	Auth       *auth.Authenticator
}
