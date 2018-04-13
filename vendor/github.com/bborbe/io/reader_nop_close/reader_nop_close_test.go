package reader_nop_close

import (
	"testing"

	"io"

	. "github.com/bborbe/assert"
)

func TestImplementsReadCloser(t *testing.T) {
	b := New(nil)
	var expected *io.ReadCloser
	if err := AssertThat(b, Implements(expected)); err != nil {
		t.Fatal(err)
	}
}
