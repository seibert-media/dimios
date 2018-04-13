package main

import (
	"testing"

	"sync"

	"bytes"

	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	writer := bytes.NewBufferString("")
	input := bytes.NewBufferString("")
	wg := new(sync.WaitGroup)
	err := do(writer, input, 2, wg, nil, "/tmp")
	err = AssertThat(err, NilValue())
	if err != nil {
		t.Fatal(err)
	}
}
