package mock

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestNewHttpRequestMock(t *testing.T) {
	r, err := NewHttpRequestMock("http://www.example.com/asdf?foo=bar")
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(r, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}

func TestRequestUri(t *testing.T) {
	r, err := NewHttpRequestMock("http://www.example.com/asdf?foo=bar")
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(r.RequestURI, Is("/asdf?foo=bar"))
	if err != nil {
		t.Fatal(err)
	}
}
