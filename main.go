package main

import (
	"fmt"
	"os"
	"p2p-tun/cmd"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "p2p-tun",
		Usage: "port forward and tun2socks through libp2p",
		Flags: cmd.MainFlags,
		Commands: []*cli.Command{
			cmd.ClientCmd(),
			cmd.ServerCmd(),
		},
		Before: cmd.Common,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
