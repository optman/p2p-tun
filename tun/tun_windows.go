package tun

import "fmt"

func NewTun(tun_name string) (fd int, mtu uint32, err error) {
	return 0, 0, fmt.Errorf("tun mod not supported on Windows")
}
