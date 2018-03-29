package apply

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestNewApplier(t *testing.T) {
	a := NewApplier()
	if err := AssertThat(a, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
