package cmd

import (
	"context"
	"fmt"
	"p2p-tun/cmd/port"
	"p2p-tun/host"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/urfave/cli/v2"
)

func ClientCmd() *cli.Command {
	return &cli.Command{
		Name:  "client",
		Usage: "start client node",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "server-id",
				Usage:    "server peer id",
				Required: true,
			},
		},
		Subcommands: []*cli.Command{
			port.ClientCmd(),
		},
		Before: connect,
	}
}

func connect(c *cli.Context) error {
	server_id, err := peer.Decode(c.String("server-id"))
	if err != nil {
		return fmt.Errorf("invalid server id, %s", err)
	}

	log.Info("connecting")
	client, err := host.NewClient(c.Context, c.Int("listen-port"), id_seed)
	if err != nil {
		panic(err)
	}
	client.Connect(server_id)

	log.Info("connected")

	c.Context = context.WithValue(c.Context, "client", client)

	return nil
}
