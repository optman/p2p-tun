package main

import (
	"fmt"
	"os"
	"github.com/optman/p2p-tun/cmd"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:     "p2p-tun",
		Usage:    "port forward and tun2socks through libp2p",
		Flags:    cmd.Flags,
		Commands: cmd.Commands,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
