package apply

import (
	"testing"

	. "github.com/bborbe/assert"
	restclient "k8s.io/client-go/rest"
)

func TestNew(t *testing.T) {
	config := &restclient.Config{}
	a, err := New(config)
	if err != nil {
		t.Fatal("Apply_New() failed with", err)
	}
	if err := AssertThat(a, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
