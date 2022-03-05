package main

import (
	"fmt"
	"os"
	"p2p-tun/cmd"

	"github.com/urfave/cli"
)

func main() {

	app := &cli.App{
		Name:  "p2p-tun",
		Flags: cmd.MainFlags,
		Commands: []cli.Command{
			cmd.Client(),
			cmd.Server(),
		},
		Before: cmd.Common,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
