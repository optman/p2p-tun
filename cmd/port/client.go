package port

import (
	"github.com/optman/p2p-tun/cmd/context"
	"github.com/optman/p2p-tun/port"

	"github.com/urfave/cli/v2"
)

func ClientCmd() *cli.Command {
	return &cli.Command{
		Name:  "port",
		Usage: "start port forward client",
		Flags: []cli.Flag{
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
	ctx := context.Context{c.Context}
	select {
	case <-ctx.HostReady():
	case <-ctx.Done():
		return nil
	}
	return port.RunClient(ctx, c.String("local-address"), ctx.Client().CreateStream(port.ProtocolID))
}
