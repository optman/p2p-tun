package cmd

import (
	"fmt"
	"github.com/optman/p2p-tun/cmd/context"
	"github.com/optman/p2p-tun/cmd/port"
	"github.com/optman/p2p-tun/cmd/tun"
	"github.com/optman/p2p-tun/host"

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
				EnvVars:  []string{"SERVER_ID"},
				Required: true,
			},
		},
		Subcommands: []*cli.Command{
			port.ClientCmd(),
			tun.TunCmd(),
		},
		Before: func(c *cli.Context) error {
			if err := common(c); err != nil {
				return err
			}

			return connect(c)
		},
	}
}

func connect(c *cli.Context) error {
	ctx := context.Context{c.Context}
	conf := ctx.NodeConfig()

	server_id, err := peer.Decode(c.String("server-id"))
	if err != nil {
		return fmt.Errorf("invalid server id, %s", err)
	}

	log.Info("connecting")
	client, err := host.NewClient(c.Context, conf)
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

	c.Context = context.SetClient(c.Context, client)
	c.Context = context.SetHostReady(c.Context, readyChan)

	return nil
}
