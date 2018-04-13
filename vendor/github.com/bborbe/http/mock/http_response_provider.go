package mock

import (
	"net/http"

	"bytes"
	"io/ioutil"
)

type ResponseProvider interface {
	GetResponse() *http.Response
	GetError() error
}

type responseProvider struct {
	content string
	err     error
	status  int
}

func NewResponseProvider(status int, content string, err error) *responseProvider {
	p := new(responseProvider)
	p.status = status
	p.content = content
	p.err = err
	return p
}

func (p *responseProvider) GetResponse() *http.Response {
	response := new(http.Response)
	response.StatusCode = p.status
	response.Body = ioutil.NopCloser(bytes.NewBufferString(p.content))
	return response
}

func (p *responseProvider) GetError() error {
	return p.err
}
