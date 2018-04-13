package writer_nop_close

import (
	"testing"

	"io"

	. "github.com/bborbe/assert"
)

func TestImplementsWriteCloser(t *testing.T) {
	b := New(nil)
	var expected *io.WriteCloser
	if err := AssertThat(b, Implements(expected)); err != nil {
		t.Fatal(err)
	}
}
