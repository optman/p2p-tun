package tun

import (
	"github.com/optman/p2p-tun/cmd/context"
	"github.com/optman/p2p-tun/tun"

	"github.com/urfave/cli/v2"
)

func SocksCmd() *cli.Command {
	return &cli.Command{
		Name:   "socks",
		Usage:  "start socks5 server",
		Action: startSocks5Server,
	}
}

func startSocks5Server(c *cli.Context) error {

	ctx := context.Context{c.Context}

	ctx.Logger().Info("start socks5 server")

	ctx.Server().HandleStream(tun.TcpProtocolID, handleSocks5TcpStreamFunc(ctx))
	ctx.Server().HandleStream(tun.UdpProtocolID, handleSocks5UdpStreamFunc(ctx))

	select {
	case <-ctx.Done():
	}

	return nil
}
