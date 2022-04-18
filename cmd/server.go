package cmd

import (
	"errors"
	"fmt"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/optman/p2p-tun/cmd/context"
	"github.com/optman/p2p-tun/cmd/port"
	"github.com/optman/p2p-tun/cmd/tun"
	"github.com/optman/p2p-tun/host"
	ra "github.com/optman/rndz-multiaddr"

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

	var serverAddrs []ma.Multiaddr
	for _, a := range conf.ListenAddrs {
		addr, err := ma.NewMultiaddr(a)
		if err != nil {
			return fmt.Errorf("invalid listen address, %s", addr)
		}
		_, rndzServer := ra.SplitListenAddr(addr)
		if rndzServer == nil {
			return errors.New("invalid listen address, no rndz server set")
		}

		serverAddrs = append(serverAddrs, ra.NewDialAddr(rndzServer, server.Host().ID()))
	}

	log.Infof("server addr %v", serverAddrs)

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
