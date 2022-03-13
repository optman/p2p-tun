package tun

import (
	"io"
	"time"
)

type Stream interface {
	io.Reader
	io.Writer
	io.Closer
	SetDeadline(time.Time) error
}
