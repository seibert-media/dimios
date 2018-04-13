package mock

import (
	"net/http"
	"testing"

	. "github.com/bborbe/assert"
)

func TestNewHttpResponseWriterMock(t *testing.T) {
	response := NewHttpResponseWriterMock()
	var expected *http.ResponseWriter
	err := AssertThat(response, Implements(expected).Message("check type"))
	if err != nil {
		t.Error(err)
	}
}

func TestHttpResponseWriter(t *testing.T) {
	var err error
	response := NewHttpResponseWriterMock()
	err = AssertThat(response.String(), Is(""))
	if err != nil {
		t.Error(err)
	}
	response.Write([]byte("hello"))
	err = AssertThat(response.String(), Is("hello"))
	if err != nil {
		t.Error(err)
	}
	response.Write([]byte(" world"))
	err = AssertThat(response.String(), Is("hello world"))
	if err != nil {
		t.Error(err)
	}
}
