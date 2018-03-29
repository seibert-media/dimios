package sync

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestNewSyncer(t *testing.T) {
	a := NewSyncer(nil, nil)
	if err := AssertThat(a, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
