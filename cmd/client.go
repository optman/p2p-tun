package cmd

import (
	"context"
	"fmt"
	"log"
	"p2p-tun/host"
	"p2p-tun/port"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/urfave/cli"
)

func Client() cli.Command {
	return cli.Command{
		Name:  "client",
		Usage: "start client node",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "server-id",
				Usage:    "server peer id",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "local-address",
				Usage: "local listen port",
				Value: "0.0.0.0:0",
			},
		},
		Action: doClient,
	}
}

func doClient(c *cli.Context) error {
	server_id, err := peer.Decode(c.String("server-id"))
	if err != nil {
		return fmt.Errorf("invalid server id, %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("connecting")
	client, err := host.NewClient(ctx, c.Int("listen-port"), id_seed)
	if err != nil {
		panic(err)
	}
	client.Connect(server_id)

	log.Println("connected")

	return port.RunClient(ctx, c.String("local-address"), client.CreateStream(port.ProtocolID))
}
