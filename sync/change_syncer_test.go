package sync

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestNew(t *testing.T) {
	a := New(nil, nil)
	if err := AssertThat(a, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
