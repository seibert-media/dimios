package file

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestGetObjectsNamespaceError(t *testing.T) {
	p := New("/tmp", nil)
	_, err := p.GetObjects("invalid-ns")
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(err.Error(), Is("namespace invalid-ns not found")); err != nil {
		t.Fatal(err)
	}
}
