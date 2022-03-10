package port

import (
	"context"
	"io"
	"net"
	"github.com/optman/p2p-tun/util"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/protocol"
)

var (
	log = logging.Logger("p2p-tun")
)

const ProtocolID = protocol.ID("/pfw")

func handle_stream(src io.ReadWriteCloser, forward_addr string) {

	if dst, err := net.Dial("tcp", forward_addr); err == nil {
		util.ConcatStream(src, dst)

	} else {
		src.Close()
		log.Warn(err)
	}
}

func HandleStream(forward_addr string) func(s io.ReadWriteCloser) {

	return func(s io.ReadWriteCloser) {
		handle_stream(s, forward_addr)
	}

}

type NewStream func(ctx context.Context) (io.ReadWriteCloser, error)

func RunClient(ctx context.Context, local_addr string, newStream NewStream) error {

	ln, err := net.Listen("tcp", local_addr)
	if err != nil {
		return err
	}

	log.Info("local_addr: ", ln.Addr())

	for {
		src, err := ln.Accept()
		if err != nil {
			return err
		}

		go func() {
			dst, err := newStream(ctx)
			if err != nil {
				src.Close()
				log.Warn("stream open fail")
				return
			}

			util.ConcatStream(src, dst)

		}()
	}

	return nil

}
