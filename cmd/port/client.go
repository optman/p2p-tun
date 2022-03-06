package port

import (
	"p2p-tun/host"
	"p2p-tun/port"

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
	readyChan := c.Context.Value("ready").(chan struct{})
	select {
	case <-readyChan:
	case <-c.Context.Done():
		return nil
	}
	client := c.Context.Value("client").(*host.Client)
	return port.RunClient(c.Context, c.String("local-address"), client.CreateStream(port.ProtocolID))
}
