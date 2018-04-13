package requestbuilder

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsNewHTTPRequestBuilderProvider(t *testing.T) {
	p := NewHTTPRequestBuilderProvider()
	var i *HTTPRequestBuilderProvider
	err := AssertThat(p, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewHTTPRequestBuilder(t *testing.T) {
	var err error
	p := NewHTTPRequestBuilderProvider()
	rb := p.NewHTTPRequestBuilder("http://example.com")
	err = AssertThat(rb, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	var i *HttpRequestBuilder
	err = AssertThat(rb, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}
