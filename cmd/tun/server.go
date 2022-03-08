package tun

import (
	"p2p-tun/cmd/context"
	"p2p-tun/tun"

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

	svr, err := newSocks5Server()
	if err != nil {
		return err
	}

	ctx.Server().HandleStream(tun.ProtocolID, handleSocks5StreamFunc(ctx, svr))

	select {
	case <-ctx.Done():
	}

	return nil
}
