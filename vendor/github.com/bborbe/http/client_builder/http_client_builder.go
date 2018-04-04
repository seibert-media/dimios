package client_builder

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"net/url"

	"errors"

	"github.com/golang/glog"
)

const (
	DEFAULT_TIMEOUT     = 30 * time.Second
	KEEPALIVE           = 30 * time.Second
	TLSHANDSHAKETIMEOUT = 10 * time.Second
)

type DialFunc func(network, address string) (net.Conn, error)

type HttpClientBuilder interface {
	Build() *http.Client
	BuildRoundTripper() http.RoundTripper
	WithProxy() HttpClientBuilder
	WithoutProxy() HttpClientBuilder
	WithRedirects() HttpClientBuilder
	WithoutRedirects() HttpClientBuilder
	WithTimeout(timeout time.Duration) HttpClientBuilder
	WithDialFunc(dialFunc DialFunc) HttpClientBuilder
}

type httpClientBuilder struct {
	proxy         Proxy
	checkRedirect CheckRedirect
	timeout       time.Duration
	dialFunc      DialFunc
}

type Proxy func(req *http.Request) (*url.URL, error)

type CheckRedirect func(req *http.Request, via []*http.Request) error

func New() *httpClientBuilder {
	b := new(httpClientBuilder)
	b.WithoutProxy()
	b.WithRedirects()
	b.timeout = DEFAULT_TIMEOUT
	return b
}

func (h *httpClientBuilder) WithTimeout(timeout time.Duration) HttpClientBuilder {
	h.timeout = timeout
	return h
}

func (h *httpClientBuilder) WithDialFunc(dialFunc DialFunc) HttpClientBuilder {
	h.dialFunc = dialFunc
	return h
}

func (b *httpClientBuilder) BuildDialFunc() DialFunc {
	if b.dialFunc != nil {
		return b.dialFunc
	}
	return (&net.Dialer{
		Timeout: b.timeout,
		//		KeepAlive: KEEPALIVE,
	}).Dial
}

func (b *httpClientBuilder) BuildRoundTripper() http.RoundTripper {
	if glog.V(5) {
		glog.Infof("build http transport")
	}
	return &http.Transport{
		Proxy:           b.proxy,
		Dial:            b.BuildDialFunc(),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//		TLSHandshakeTimeout: TLSHANDSHAKETIMEOUT,
	}
}

func (b *httpClientBuilder) Build() *http.Client {
	if glog.V(5) {
		glog.Infof("build http client")
	}
	return &http.Client{
		Transport:     b.BuildRoundTripper(),
		CheckRedirect: b.checkRedirect,
	}
}

func (b *httpClientBuilder) WithProxy() HttpClientBuilder {
	b.proxy = http.ProxyFromEnvironment
	return b
}

func (b *httpClientBuilder) WithoutProxy() HttpClientBuilder {
	b.proxy = nil
	return b
}

func (b *httpClientBuilder) WithRedirects() HttpClientBuilder {
	b.checkRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}
	return b
}

func (b *httpClientBuilder) WithoutRedirects() HttpClientBuilder {
	b.checkRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 1 {
			return errors.New("redirects")
		}
		return nil
	}
	return b
}
