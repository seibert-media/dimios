package buffered_writer

import (
	"bytes"
	"io"
	"testing"
	"time"

	. "github.com/bborbe/assert"
)

func TestImplementsWriter(t *testing.T) {
	writer := bytes.NewBufferString("")
	b := NewBufferedWriter(writer)
	var expected *io.WriteCloser
	err := AssertThat(b, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestWrite(t *testing.T) {
	writer := bytes.NewBufferString("")
	b := NewBufferedWriter(writer)
	l, err := b.Write([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(l, Is(5))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	err = AssertThat(writer.String(), Is("hello"))
	if err != nil {
		t.Fatal(err)
	}
	b.Write([]byte(" world"))
	time.Sleep(100 * time.Millisecond)
	err = AssertThat(writer.String(), Is("hello world"))
	if err != nil {
		t.Fatal(err)
	}
}

//func TestWriteAfterClose(t *testing.T) {
//	writer := bytes.NewBufferString("")
//	b := NewBufferedWriter(writer)
//	b.Close()
//	time.Sleep(100 * time.Millisecond)
//	b.Write([]byte("hello"))
//	err := AssertThat(writer.String(), Is(""))
//	if err != nil {
//		t.Fatal(err)
//	}
//}
