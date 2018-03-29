package reader_shadow_copy

import (
	"bytes"
	"io"
)

type reader struct {
	readcloser io.ReadCloser
	bytes      bytes.Buffer
}

type ShadowCopyReader interface {
	io.ReadCloser
	Bytes() []byte
}

func New(readcloser io.ReadCloser) *reader {
	r := new(reader)
	r.readcloser = readcloser
	return r
}

func (r *reader) Close() error {
	return r.readcloser.Close()
}

func (r *reader) Read(p []byte) (int, error) {
	n, err := r.readcloser.Read(p)
	r.bytes.Write(p[:n])
	return n, err
}

func (r *reader) Bytes() []byte {
	return r.bytes.Bytes()
}
