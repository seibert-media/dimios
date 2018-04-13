package requestbuilder

type HTTPRequestBuilderProvider interface {
	NewHTTPRequestBuilder(url string) HttpRequestBuilder
}

type httpRequestBuilderProvider struct {
}

func NewHTTPRequestBuilderProvider() *httpRequestBuilderProvider {
	p := new(httpRequestBuilderProvider)
	return p
}

func (p *httpRequestBuilderProvider) NewHTTPRequestBuilder(url string) HttpRequestBuilder {
	return NewHTTPRequestBuilder(url)
}
