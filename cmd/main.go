package cmd

import (
	"math/rand"
	"p2p-tun/cmd/context"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
)

var (
	log     = logging.Logger("p2p-tun")
	id_seed int64
)

var MainFlags = []cli.Flag{
	&cli.Int64Flag{
		Name:  "id",
		Usage: "id seed",
		Value: 0,
	},
	&cli.IntFlag{
		Name:  "listen-port",
		Usage: "p2p listen port",
		Value: 0,
	},
}

func Common(c *cli.Context) error {
	id_seed = c.Int64("id")
	if id_seed == 0 {
		id_seed = rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
	}

	logging.SetLogLevel("p2p-tun", "info")

	c.Context = context.SetLogger(c.Context, log)

	return nil
}
