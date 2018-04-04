package manager

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestValidateReturnError(t *testing.T) {
	m := &Manager{}
	if err := AssertThat(m.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
