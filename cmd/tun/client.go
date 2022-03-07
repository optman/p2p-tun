package tun

import (
	"io"
	"net"
	"p2p-tun/cmd/context"
	"p2p-tun/tun"
	"p2p-tun/util"

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

	ctx := context.Context{c.Context}

	//TODO: setup tun

	select {
	case <-ctx.HostReady():
	case <-ctx.Done():
		return nil
	}
	//TODO: setup netstack
	if err := tun.SetupTun(c.String("tun-name"), handleStreamFunc(ctx)); err != nil {
		return err
	}

	ctx.Logger().Info("tun ready")

	select {
	case <-ctx.Done():
		return nil
	}

	return nil
}

func handleStreamFunc(ctx context.Context) func(target_addr *net.TCPAddr, src io.ReadWriteCloser) {

	createStream := ctx.Client().CreateStream(tun.ProtocolID)
	log := ctx.Logger()

	return func(target_addr *net.TCPAddr, src io.ReadWriteCloser) {
		dst, err := createStream(ctx)
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
