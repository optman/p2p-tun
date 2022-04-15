package cmd

import (
	"errors"
	"fmt"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/optman/p2p-tun/cmd/context"
	"github.com/optman/p2p-tun/cmd/port"
	"github.com/optman/p2p-tun/cmd/tun"
	"github.com/optman/p2p-tun/host"

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
		Before: func(c *cli.Context) error {
			if err := common(c); err != nil {
				return err
			}

			return startServer(c)
		},
	}
}

func startServer(c *cli.Context) error {
	ctx := context.Context{c.Context}
	conf := ctx.NodeConfig()

	if conf.Auth == nil {
		log.Warn("Danger!!! No secret set!!!")
	}

	server, err := host.NewServer(c.Context, conf)
	if err != nil {
		return err
	}

	serverAddr, err := ma.NewMultiaddr(conf.RndzServer)
	if err != nil {
		return errors.New("invalid rndz server addr")
	}

	p2pPart, err := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", server.Host().ID()))
	if err != nil {
		panic(err)
	}

	serverAddr = serverAddr.Encapsulate(p2pPart)

	log.Infof("server addr %s", serverAddr)

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
