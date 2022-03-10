package cmd

import (
	"crypto/rand"
	"fmt"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/urfave/cli/v2"
)

func GenkeyCmd() *cli.Command {
	return &cli.Command{
		Name:  "genkey",
		Usage: "generate a private key",
		Action: func(c *cli.Context) error {
			privKey, err := genkey()
			if err != nil {
				return err
			}

			b, err := crypto.MarshalPrivateKey(privKey)
			if err != nil {
				return err
			}

			str := crypto.ConfigEncodeKey(b)

			fmt.Println(str)

			return nil
		},
	}
}

func genkey() (crypto.PrivKey, error) {
	privKey, _, err := crypto.GenerateSecp256k1Key(rand.Reader)
	return privKey, err
}
