package tun

import (
	"io"
	"net"
	"time"

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

	if err := tun.NewNetstack(fd, mtu, handleTcpConn(ctx), handleUdpConn(ctx)); err != nil {
		return err
	}

	ctx.Logger().Info("tun ready")

	select {
	case <-ctx.Done():
		return nil
	}

	return nil
}

func handleTcpConn(ctx context.Context) func(target_addr *net.TCPAddr, src tun.Stream) {

	createStream := ctx.Client().CreateStream(tun.TcpProtocolID)
	log := ctx.Logger()

	return func(target_addr *net.TCPAddr, src tun.Stream) {
		defer src.Close()

		log.Debug("tcp ", target_addr)

		dst, err := createStream(ctx)
		if err != nil {
			return
		}
		defer dst.Close()

		if err := socksConnect(dst, target_addr); err != nil {
			return
		}

		util.ConcatStream(src, dst)
	}
}

func handleUdpConn(ctx context.Context) func(target_addr *net.UDPAddr, src tun.Stream) {

	createStream := ctx.Client().CreateStream(tun.UdpProtocolID)
	log := ctx.Logger()

	return func(target_addr *net.UDPAddr, src tun.Stream) {
		defer src.Close()

		log.Debug("udp ", target_addr)

		src.SetDeadline(time.Now().Add(UDPTIMEOUT))

		dst, err := createStream(ctx)
		if err != nil {
			return
		}
		defer dst.Close()

		dst.SetDeadline(time.Now().Add(UDPTIMEOUT))

		if err := socksConnect2(dst, target_addr.IP, target_addr.Port); err != nil {
			return
		}

		buf := make([]byte, 4096)
		n, err := src.Read(buf)
		if err != nil {
			return
		}

		if _, err := dst.Write(buf[:n]); err != nil {
			return
		}

		dst.CloseWrite()

		resp, err := io.ReadAll(dst)
		if err != nil {
			return
		}

		if _, err := src.Write(resp); err != nil {
			return
		}
	}
}
