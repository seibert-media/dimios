package apply

import (
	"testing"

	. "github.com/bborbe/assert"
	restclient "k8s.io/client-go/rest"
)

func TestNew(t *testing.T) {
	applier, err := New(true, &restclient.Config{})
	if err != nil {
		t.Fatal("Apply_New() failed with", err)
	}
	if err := AssertThat(applier, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
