package cmd

import (
	"context"
	"fmt"
	"p2p-tun/cmd/port"
	"p2p-tun/cmd/tun"
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
			tun.TunCmd(),
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

	readyChan := make(chan struct{})

	go func() {
		if err := client.Connect(server_id); err != nil {
			log.Fatal(err)
		}
		log.Info("connected")
		close(readyChan)
	}()

	c.Context = context.WithValue(c.Context, "client", client)
	c.Context = context.WithValue(c.Context, "ready", readyChan)

	return nil
}
