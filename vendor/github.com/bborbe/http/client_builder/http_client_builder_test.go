package client_builder

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsHttpClientBuilder(t *testing.T) {
	b := New()
	var i *HttpClientBuilder
	if err := AssertThat(b, Implements(i).Message("check type")); err != nil {
		t.Fatal(err)
	}
}

func TestBuildReturnNotNilValue(t *testing.T) {
	b := New()
	if err := AssertThat(b.Build(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
