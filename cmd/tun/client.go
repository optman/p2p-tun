package tun

import (
	"io"
	"net"
	"github.com/optman/p2p-tun/cmd/context"
	"github.com/optman/p2p-tun/tun"
	"github.com/optman/p2p-tun/util"

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
				Value: "",
			},
		},
		Action: startTun,
	}
}

func startTun(c *cli.Context) error {

	fd, mtu, err := tun.NewTun(c.String("tun-name"))
	if err != nil {
		return err
	}

	ctx := context.Context{c.Context}
	select {
	case <-ctx.HostReady():
	case <-ctx.Done():
		return nil
	}

	if err := tun.NewNetstack(fd, mtu, handleStreamFunc(ctx)); err != nil {
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
