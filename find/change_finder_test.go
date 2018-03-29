package find

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestNew(t *testing.T) {
	a := NewFinder("")
	if err := AssertThat(a, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
