package util

import (
	"io"
	"sync"
)

type closeWriter interface {
	CloseWrite() error
}

func ConcatStream(src, dst io.ReadWriteCloser) {
	var wg sync.WaitGroup

	cp := func(dst, src io.ReadWriteCloser) {
		defer wg.Done()

		io.Copy(dst, src)
		if c, ok := dst.(closeWriter); ok { //libp2p nework.steam provide CloseWrite
			c.CloseWrite()
		}
	}

	wg.Add(2)

	go cp(dst, src)
	go cp(src, dst)

	wg.Wait()
}
