package tun

import (
	"gvisor.dev/gvisor/pkg/tcpip/link/rawfile"
	"gvisor.dev/gvisor/pkg/tcpip/link/tun"
)

func NewTun(tun_name string) (fd int, mtu uint32, err error) {
	fd, err = tun.Open(tun_name)
	if err != nil {
		return
	}

	mtu, err = rawfile.GetMTU(tun_name)
	if err != nil {
		return
	}

	return
}
