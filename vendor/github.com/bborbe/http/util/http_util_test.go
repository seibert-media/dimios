package util

import (
	"net/http"
	"testing"

	"net/url"

	"bytes"
	"io/ioutil"

	. "github.com/bborbe/assert"
)

func TestResponseToByteArray(t *testing.T) {
	var err error
	var content []byte

	response := new(http.Response)
	response.Body = ioutil.NopCloser(bytes.NewBufferString("test"))

	if content, err = ResponseToByteArray(response); err != nil {
		t.Fatal(err)
	}

	if err := AssertThat(string(content), Is("test")); err != nil {
		t.Fatal(err)
	}
}

func TestResponseToString(t *testing.T) {
	var err error
	var content string

	response := new(http.Response)
	response.Body = ioutil.NopCloser(bytes.NewBufferString("test"))

	if content, err = ResponseToString(response); err != nil {
		t.Fatal(err)
	}

	if err := AssertThat(content, Is("test")); err != nil {
		t.Fatal(err)
	}
}

func TestFindFileExtension(t *testing.T) {
	var err error
	response := &http.Response{}
	ext, err := FindFileExtension(response)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(ext, Is("")); err != nil {
		t.Fatal(err)
	}
}

func TestFindFileExtensionUrlWithDot(t *testing.T) {
	var err error
	var u *url.URL
	if u, err = url.ParseRequestURI("http://www.example/robots.txt"); err != nil {
		t.Fatal(err)
	}
	response := &http.Response{Request: &http.Request{URL: u}}
	ext, err := FindFileExtension(response)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(ext, Is("txt")); err != nil {
		t.Fatal(err)
	}
}

func TestFindFileExtensionUrlWithDotAtLast(t *testing.T) {
	var err error
	var u *url.URL
	if u, err = url.ParseRequestURI("http://www.example/robots."); err != nil {
		t.Fatal(err)
	}
	response := &http.Response{Request: &http.Request{URL: u}}
	ext, err := FindFileExtension(response)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(ext, Is("")); err != nil {
		t.Fatal(err)
	}
}

func TestFindFileExtensionHeader(t *testing.T) {
	var err error
	response := &http.Response{Header: http.Header{}}
	ext, err := FindFileExtension(response)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(ext, Is("")); err != nil {
		t.Fatal(err)
	}
}

func TestFindFileExtensionHeaderContentTypeKownType(t *testing.T) {
	var err error
	response := &http.Response{Header: http.Header{"Content-Type": []string{"image/jpeg"}}}
	ext, err := FindFileExtension(response)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(ext, Is("jpg")); err != nil {
		t.Fatal(err)
	}
}

func TestFindFileExtensionHeaderContentTypeUnkownType(t *testing.T) {
	var err error
	response := &http.Response{Header: http.Header{"Content-Type": []string{"text/foo"}}}
	ext, err := FindFileExtension(response)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(ext, Is("")); err != nil {
		t.Fatal(err)
	}
}
