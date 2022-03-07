package cmd

import (
	"p2p-tun/cmd/context"
	"p2p-tun/cmd/port"
	"p2p-tun/cmd/tun"
	"p2p-tun/host"

	"github.com/urfave/cli/v2"
)

func ServerCmd() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "start server node",
		Subcommands: []*cli.Command{
			port.ServerCmd(),
			tun.SocksCmd(),
		},
		Before: startServer,
	}
}

func startServer(c *cli.Context) error {

	server, err := host.NewServer(c.Context, c.Int("listen-port"), id_seed)
	if err != nil {
		return err
	}

	log.Infof("server id %s", server.Host().ID())

	readyChan := make(chan struct{})

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
		close(readyChan)
	}()

	c.Context = context.SetServer(c.Context, server)
	c.Context = context.SetHostReady(c.Context, readyChan)

	return nil
}
