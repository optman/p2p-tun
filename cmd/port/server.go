package port

import (
	"p2p-tun/host"
	"p2p-tun/port"

	logging "github.com/ipfs/go-log/v2"
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
	readyChan := c.Context.Value("ready").(chan struct{})
	select {
	case <-readyChan:
	case <-c.Context.Done():
		return nil
	}

	server := c.Context.Value("server").(*host.Server)
	log := c.Context.Value("logger").(logging.StandardLogger)

	log.Info("forward-addr:", c.String("forward-addr"))

	server.HandleStream(port.ProtocolID, port.HandleStream(c.String("forward-addr")))

	select {
	case <-c.Context.Done():
	}

	return nil
}
