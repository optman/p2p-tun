package cmd

import (
	"fmt"
	"p2p-tun/auth"
	"p2p-tun/cmd/context"
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

	conf := host.ClientConfig{
		Ctx:  c.Context,
		Port: c.Int("listen-port"),
		Seed: id_seed,
	}

	secret := c.String("secret")
	if len(secret) > 0 {
		conf.Auth = auth.NewAuthenticator(secret)
	}

	log.Info("connecting")
	client, err := host.NewClient(conf)
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
