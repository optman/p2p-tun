package util

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func SplitListenAddr(addr ma.Multiaddr) (localAddr ma.Multiaddr, rndzAddr ma.Multiaddr) {
	localAddr, rndzAddr = ma.SplitFunc(addr, func(c ma.Component) bool {
		return c.Protocol().Code == ma.P_CIRCUIT
	})

	if rndzAddr != nil {
		_, rndzAddr = ma.SplitFirst(rndzAddr)
	}

	return
}

func NewServerAddr(rndzServer ma.Multiaddr, peerId peer.ID) ma.Multiaddr {
	p2pPart, err := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", peerId))
	if err != nil {
		panic(err)
	}

	return rndzServer.Encapsulate(p2pPart)
}
