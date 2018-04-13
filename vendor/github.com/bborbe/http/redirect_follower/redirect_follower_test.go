package redirect_follower

import (
	"testing"

	"fmt"
	"net/http"

	"net/url"

	. "github.com/bborbe/assert"
	"github.com/bborbe/http/client_builder"
	"github.com/bborbe/http/requestbuilder"
)

func TestImplementsRedirectFollower(t *testing.T) {
	r := New(nil)
	var i *RedirectFollower
	err := AssertThat(r, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSuccess(t *testing.T) {
	expectedResponse := &http.Response{}
	expectedRequest := &http.Request{}
	var parameterRequest *http.Request
	var counter int
	r := New(func(req *http.Request) (resp *http.Response, err error) {
		counter++
		parameterRequest = req
		return expectedResponse, nil
	})

	resultResponse, resultErr := r.ExecuteRequestAndFollow(expectedRequest)
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(resultErr, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(resultResponse, Is(expectedResponse)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(parameterRequest, Is(expectedRequest)); err != nil {
		t.Fatal(err)
	}
}

func TestFailure(t *testing.T) {
	expectedError := fmt.Errorf("foo")
	expectedRequest := &http.Request{}
	var parameterRequest *http.Request
	var counter int
	r := New(func(req *http.Request) (resp *http.Response, err error) {
		counter++
		parameterRequest = req
		return nil, expectedError
	})

	resultResponse, resultErr := r.ExecuteRequestAndFollow(expectedRequest)
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(resultErr, Is(expectedError)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(resultResponse, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(parameterRequest, Is(expectedRequest)); err != nil {
		t.Fatal(err)
	}
}

func TestIntegrated(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test due to -short")
	}
	httpClientBuilder := client_builder.New()
	b := New(httpClientBuilder.BuildRoundTripper().RoundTrip)
	rb := requestbuilder.NewHTTPRequestBuilder("http://www.benjamin-borbe.de")
	request, err := rb.Build()
	if err != nil {
		t.Fatal(err)
	}
	response, err := b.ExecuteRequestAndFollow(request)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(response, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(response.StatusCode, Is(200)); err != nil {
		t.Fatal(err)
	}
}

func TestLocationToUrlRelativLocation(t *testing.T) {
	u, _ := url.Parse("http://www.benjamin-borbe.de/test/login")
	result, err := locationToUrl(u, "/denied")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(result.String(), Is("http://www.benjamin-borbe.de/denied")); err != nil {
		t.Fatal(err)
	}
}

func TestLocationToUrlEmptyLocation(t *testing.T) {
	u, _ := url.Parse("http://www.benjamin-borbe.de/test/login")
	_, err := locationToUrl(u, "")
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestLocationToUrlAbsolutLocation(t *testing.T) {
	u, _ := url.Parse("http://www.benjamin-borbe.de/test/login")
	result, err := locationToUrl(u, "http://www.example.com/news")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(result.String(), Is("http://www.example.com/news")); err != nil {
		t.Fatal(err)
	}
}
