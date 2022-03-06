package tun

import (
	"context"
	"io"
	"net"
	"p2p-tun/host"
	"p2p-tun/tun"
	"p2p-tun/util"

	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
)

func TunCmd() *cli.Command {
	return &cli.Command{
		Name:  "tun",
		Usage: "start tun mode ",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tun-name",
				Usage: "tun interface name",
				Value: "tun0",
			},
		},
		Action: startTun,
	}
}

func startTun(c *cli.Context) error {

	//TODO: setup tun

	readyChan := c.Context.Value("ready").(chan struct{})
	select {
	case <-readyChan:
	case <-c.Context.Done():
		return nil
	}
	//TODO: setup netstack
	if err := tun.SetupTun(c.String("tun-name"), handleStreamFunc(c.Context)); err != nil {
		return err
	}

	log := c.Context.Value("logger").(logging.StandardLogger)
	log.Info("tun ready")

	select {
	case <-c.Context.Done():
		return nil
	}

	return nil
}

func handleStreamFunc(ctx context.Context) func(target_addr *net.TCPAddr, src io.ReadWriteCloser) {

	createStream := ctx.Value("client").(*host.Client).CreateStream(tun.ProtocolID)
	log := ctx.Value("logger").(logging.StandardLogger)

	return func(target_addr *net.TCPAddr, src io.ReadWriteCloser) {
		dst, err := createStream(context.Background())
		if err != nil {
			log.Debug("create socks stream fail ", err)
			return
		}
		defer dst.Close()

		if err := socksConnect(dst, target_addr); err != nil {
			log.Debugf("socks connect %v fail", target_addr)
			return
		}

		util.ConcatStream(src, dst)
	}
}
