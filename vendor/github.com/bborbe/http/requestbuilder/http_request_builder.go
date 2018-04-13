package requestbuilder

import (
	"io"
	"net/http"
	"net/url"
)

type HttpRequestBuilder interface {
	AddParameter(key string, values ...string) HttpRequestBuilder
	AddHeader(key string, values ...string) HttpRequestBuilder
	SetMethod(key string) HttpRequestBuilder
	SetBody(reader io.Reader) HttpRequestBuilder
	AddBasicAuth(username, password string) HttpRequestBuilder
	AddContentType(contentType string) HttpRequestBuilder
	SetContentLength(contentLength int64) HttpRequestBuilder
	Build() (*http.Request, error)
}

type httpRequestBuilder struct {
	url           string
	parameter     map[string][]string
	header        http.Header
	method        string
	body          io.Reader
	username      string
	password      string
	contentLength int64
}

func NewHTTPRequestBuilder(url string) *httpRequestBuilder {
	r := new(httpRequestBuilder)
	r.method = "GET"
	r.url = url
	r.parameter = make(map[string][]string)
	r.header = make(http.Header)
	return r
}

func (r *httpRequestBuilder) AddContentType(contentType string) HttpRequestBuilder {
	r.AddHeader("Content-Type", contentType)
	return r
}

func (r *httpRequestBuilder) AddBasicAuth(username, password string) HttpRequestBuilder {
	r.username = username
	r.password = password
	return r
}

func (r *httpRequestBuilder) SetBody(body io.Reader) HttpRequestBuilder {
	r.body = body
	return r
}

func (r *httpRequestBuilder) SetMethod(method string) HttpRequestBuilder {
	r.method = method
	return r
}

func (r *httpRequestBuilder) SetContentLength(contentLength int64) HttpRequestBuilder {
	r.contentLength = contentLength
	return r
}

func (r *httpRequestBuilder) AddHeader(key string, values ...string) HttpRequestBuilder {
	for _, v := range values {
		r.header.Add(key, v)
	}
	return r
}

func (r *httpRequestBuilder) AddParameter(key string, values ...string) HttpRequestBuilder {
	r.parameter[key] = values
	return r
}

func (r *httpRequestBuilder) Build() (*http.Request, error) {
	req, err := http.NewRequest(r.method, r.getUrlWithParameter(), r.body)
	if err != nil {
		return nil, err
	}
	req.Header = r.header
	if len(r.username) > 0 || len(r.password) > 0 {
		req.SetBasicAuth(r.username, r.password)
	}
	req.ContentLength = r.contentLength
	return req, nil
}

func (r *httpRequestBuilder) getUrlWithParameter() string {
	result := r.url
	first := true
	for key, values := range r.parameter {
		for _, value := range values {
			if first {
				first = false
				result += "?"
			} else {
				result += "&"
			}
			result += url.QueryEscape(key)
			result += "="
			result += url.QueryEscape(value)
		}
	}
	return result
}
