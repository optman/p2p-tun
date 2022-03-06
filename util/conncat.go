package util

import (
	"io"
	"sync"
)

func ConcatStream(src, dst io.ReadWriteCloser) {
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
