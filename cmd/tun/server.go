package tun

import (
	"p2p-tun/host"
	"p2p-tun/tun"

	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
)

func SocksCmd() *cli.Command {
	return &cli.Command{
		Name:   "socks",
		Usage:  "start socks5 server",
		Action: startSocks5Server,
	}
}

func startSocks5Server(c *cli.Context) error {

	readyChan := c.Context.Value("ready").(chan struct{})
	select {
	case <-readyChan:
	case <-c.Context.Done():
		return nil
	}

	server := c.Context.Value("server").(*host.Server)
	log := c.Context.Value("logger").(logging.StandardLogger)

	log.Info("start socks5 server")

	svr, err := newSocks5Server()
	if err != nil {
		return err
	}

	server.HandleStream(tun.ProtocolID, handleSocks5StreamFunc(c.Context, svr))

	select {
	case <-c.Context.Done():
	}

	return nil
}
