package main

import (
	"testing"

	"bytes"

	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	writer := bytes.NewBufferString("")
	input := bytes.NewBufferString("")
	err := do(writer, input)
	err = AssertThat(err, NilValue())
	if err != nil {
		t.Fatal(err)
	}
}
