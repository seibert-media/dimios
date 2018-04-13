package mock

import (
	"bytes"
	"net/http"
)

type responseWriterMock struct {
	status int
	writer *bytes.Buffer
	header http.Header
}

func NewHttpResponseWriterMock() *responseWriterMock {
	r := new(responseWriterMock)
	r.header = make(http.Header)
	r.writer = bytes.NewBufferString("")
	return r
}

func (r *responseWriterMock) Header() http.Header {
	return r.header
}

func (r *responseWriterMock) Write(b []byte) (int, error) {
	return r.writer.Write(b)
}

func (r *responseWriterMock) WriteHeader(status int) {
	r.status = status
}

func (r *responseWriterMock) Status() int {
	return r.status
}

func (r *responseWriterMock) String() string {
	return r.writer.String()
}

func (r *responseWriterMock) Bytes() []byte {
	return r.writer.Bytes()
}
