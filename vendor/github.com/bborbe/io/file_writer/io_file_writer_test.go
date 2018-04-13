package file_writer

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsWriter(t *testing.T) {
	b, _ := NewFileWriter("")
	var expected *FileWriter
	err := AssertThat(b, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
