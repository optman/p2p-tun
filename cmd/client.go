package cmd

import (
	"fmt"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/optman/p2p-tun/cmd/context"
	"github.com/optman/p2p-tun/cmd/port"
	"github.com/optman/p2p-tun/cmd/tun"
	"github.com/optman/p2p-tun/host"

	"github.com/urfave/cli/v2"
)

func ClientCmd() *cli.Command {
	return &cli.Command{
		Name:  "client",
		Usage: "start client node",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "server-addr",
				Usage:    "server multiaddr",
				EnvVars:  []string{"SERVER_ADDR"},
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

	serverAddr, err := ma.NewMultiaddr(c.String("server-addr"))
	if err != nil {
		return fmt.Errorf("invalid server addr, %s", err)
	}

	log.Info("connecting")
	client, err := host.NewClient(c.Context, conf)
	if err != nil {
		panic(err)
	}

	readyChan := make(chan struct{})

	go func() {
		if err := client.Connect(serverAddr); err != nil {
			log.Fatal(err)
		}
		log.Info("connected")
		close(readyChan)
	}()

	c.Context = context.SetClient(c.Context, client)
	c.Context = context.SetHostReady(c.Context, readyChan)

	return nil
}
