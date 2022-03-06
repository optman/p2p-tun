package cmd

import (
	"context"
	"p2p-tun/cmd/port"
	"p2p-tun/host"

	"github.com/urfave/cli/v2"
)

func ServerCmd() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "start server node",
		Subcommands: []*cli.Command{
			port.ServerCmd(),
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

	if err := server.Start(); err != nil {
		return err
	}

	c.Context = context.WithValue(c.Context, "server", server)

	return nil
}
