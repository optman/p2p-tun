package main

import (
	"context"
	"flag"
	"io"
	"net"
	"sync"
)

var local_addr = flag.String("l", "", "local addr")
var server_id = flag.String("s", "", "local addr")
var forward_addr = flag.String("f", "", "forward addr")

func loop_accept(f func(conn net.Conn)) {
	ln, err := net.Listen("tcp", *local_addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go f(conn)
	}
}

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

func handle_stream(src io.ReadWriteCloser) {

	if dst, err := net.Dial("tcp", *forward_addr); err == nil {
		concat_conn(src, dst)

	} else {
		src.Close()
		log.Error(err)
	}
}

func run_server(id string) {

	log.Infof("server id %s, forward_addr:%s", id, *forward_addr)

	select {}

}

type NewStream func(ctx context.Context) (io.ReadWriteCloser, error)

func run_client(ctx context.Context, newStream NewStream) {

	log.Infof("client local_addr:%s\n", *local_addr)

	loop_accept(func(src net.Conn) {

		dst, err := newStream(ctx)
		if err != nil {
			src.Close()
			log.Error("stream open fail")
			return
		}

		go concat_conn(src, dst)

	})

}
