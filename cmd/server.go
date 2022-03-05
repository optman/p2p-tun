package cmd

import (
	"context"
	"log"
	"p2p-tun/host"
	"p2p-tun/port"

	"github.com/urfave/cli"
)

func Server() cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "start server node",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "forward-addr",
				Usage:    "forward to address",
				Required: true,
			}},
		Action: doServer,
	}
}

func doServer(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := host.NewServer(ctx, c.Int("listen-port"), id_seed)
	if err != nil {
		return err
	}

	forward_addr := c.String("forward-addr")

	log.Printf("server id %s, forward_addr:%s", server.Host().ID(), forward_addr)

	server.Start()

	server.HandleStream(port.ProtocolID, port.HandleStream(forward_addr))

	select {
	case <-ctx.Done():
	}

	return nil
}
