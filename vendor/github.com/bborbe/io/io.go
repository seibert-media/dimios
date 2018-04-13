package io

import "io"

type Flush interface {
	Flush() error
}

type Writer interface {
	io.Writer
}

type Reader interface {
	io.Reader
}

type Closer interface {
	io.Closer
}

type WriteCloseFlusher interface {
	Writer
	Closer
	Flush
}

type WriteFlusher interface {
	Writer
	Flush
}
