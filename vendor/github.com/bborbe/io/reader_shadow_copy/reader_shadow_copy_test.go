package reader_shadow_copy

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/io/reader_nop_close"
)

func TestImplementsReadCloser(t *testing.T) {
	b := New(nil)
	var expected *io.ReadCloser
	if err := AssertThat(b, Implements(expected)); err != nil {
		t.Fatal(err)
	}
}

func newReaderCloser(content string) io.ReadCloser {
	return reader_nop_close.New(bytes.NewBufferString(content))
}

func TestImplementsShadow(t *testing.T) {
	content := "hello world"
	b := New(newReaderCloser(content))
	org, err := ioutil.ReadAll(b)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(org), Is(content)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(b.Bytes()), Is(content)); err != nil {
		t.Fatal(err)
	}
}

func TestImplementsShadowLong(t *testing.T) {
	content := createString(10000)
	b := New(newReaderCloser(content))
	org, err := ioutil.ReadAll(b)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(org), Is(content)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(b.Bytes()), Is(content)); err != nil {
		t.Fatal(err)
	}
}

func createString(l int) string {
	result := bytes.NewBufferString("")
	for i := 0; i < l; i++ {
		result.WriteString("a")
	}
	return result.String()
}
