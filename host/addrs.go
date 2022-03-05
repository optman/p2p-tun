package host

import (
	ma "github.com/multiformats/go-multiaddr"
)

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
