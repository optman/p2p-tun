package port

import (
	"context"
	"io"
	"log"
	"net"
	"sync"

	"github.com/libp2p/go-libp2p-core/protocol"
)

const ProtocolID = protocol.ID("/pfw")

func concat_conn(src, dst io.ReadWriteCloser) {
	defer src.Close()
	defer dst.Close()

	var wg sync.WaitGroup

	cp := func(dst, src io.ReadWriteCloser) {
		defer wg.Done()

		io.Copy(dst, src)
		dst.Close()
	}

	wg.Add(2)

	go cp(dst, src)
	go cp(src, dst)

	wg.Wait()
}

func handle_stream(src io.ReadWriteCloser, forward_addr string) {

	if dst, err := net.Dial("tcp", forward_addr); err == nil {
		concat_conn(src, dst)

	} else {
		src.Close()
		log.Println(err)
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

	log.Println("client local_addr: ", ln.Addr())

	for {
		src, err := ln.Accept()
		if err != nil {
			return err
		}

		go func() {
			dst, err := newStream(ctx)
			if err != nil {
				src.Close()
				log.Println("stream open fail")
				return
			}

			concat_conn(src, dst)

		}()
	}

	return nil

}
