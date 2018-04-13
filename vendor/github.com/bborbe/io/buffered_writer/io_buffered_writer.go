package buffered_writer

import (
	"io"
)

const BUFFER_SIZE = 100

type bufferedWriter struct {
	done   chan bool
	buffer chan []byte
	writer io.Writer
}

func NewBufferedWriter(writer io.Writer) *bufferedWriter {
	b := new(bufferedWriter)
	b.buffer = make(chan []byte, BUFFER_SIZE)
	b.writer = writer
	b.done = make(chan bool)
	go b.work()
	return b
}

func (b *bufferedWriter) work() {
	for p := range b.buffer {
		b.writer.Write(p)
	}
	b.done <- true
}

func (b *bufferedWriter) Write(p []byte) (n int, err error) {
	b.buffer <- p
	return len(p), nil
}

func (b *bufferedWriter) Close() error {
	close(b.buffer)
	<-b.done
	return nil
}
