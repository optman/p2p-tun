package port

import (
	"p2p-tun/cmd/context"
	"p2p-tun/port"

	"github.com/urfave/cli/v2"
)

func ServerCmd() *cli.Command {
	return &cli.Command{
		Name:  "port",
		Usage: "start port forward server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "forward-addr",
				Usage:    "forward to address",
				Required: true,
			}},
		Action: doForward,
	}
}

func doForward(c *cli.Context) error {
	ctx := context.Context{c.Context}

	ctx.Logger().Info("forward-addr:", c.String("forward-addr"))

	ctx.Server().HandleStream(port.ProtocolID, port.HandleStream(c.String("forward-addr")))

	select {
	case <-ctx.Done():
	}

	return nil
}
