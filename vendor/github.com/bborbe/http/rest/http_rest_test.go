package rest

import (
	"testing"

	"os"

	"net/http"
	"net/url"

	. "github.com/bborbe/assert"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestImplementsRest(t *testing.T) {
	c := New(nil)
	var i *Rest
	if err := AssertThat(c, Implements(i)); err != nil {
		t.Fatal(err)
	}
}

func TestUrl(t *testing.T) {
	counter := 0
	rest := New(func(req *http.Request) (resp *http.Response, err error) {
		counter++
		if err := AssertThat(req.URL.String(), Is("http://www.example.com/action")); err != nil {
			t.Fatal(err)
		}
		return new(http.Response), nil
	})
	rest.Call("http://www.example.com/action", nil, http.MethodGet, nil, nil, nil)
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestUrlWithParams(t *testing.T) {
	counter := 0
	rest := New(func(req *http.Request) (resp *http.Response, err error) {
		counter++
		if err := AssertThat(req.URL.String(), Is("http://www.example.com/action?a=b&a=c")); err != nil {
			t.Fatal(err)
		}
		return new(http.Response), nil
	})
	rest.Call("http://www.example.com/action", url.Values{"a": []string{"b", "c"}}, http.MethodGet, nil, nil, nil)
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}
