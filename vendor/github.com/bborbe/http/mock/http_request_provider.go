package mock

import (
	"net/http"
)

type RequestProvider interface {
	GetRequest() *http.Request
	GetError() error
}

type requestProvider struct {
	err error
	req *http.Request
}

func NewRequestProvider(req *http.Request, err error) *requestProvider {
	p := new(requestProvider)
	p.req = req
	p.err = err
	return p
}

func (p *requestProvider) GetRequest() *http.Request {
	return p.req
}

func (p *requestProvider) GetError() error {
	return p.err
}
