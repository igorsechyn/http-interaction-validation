package builders

import (
	"io"
	"net/http"
	"strings"
)

type RequestBuilder struct {
	method string
	body   io.Reader
	path   string
}

func (builder RequestBuilder) WithPath(value string) RequestBuilder {
	builder.path = value
	return builder
}

func (builder RequestBuilder) WithMethod(value string) RequestBuilder {
	builder.method = value
	return builder
}

func (builder RequestBuilder) WithBody(value io.Reader) RequestBuilder {
	builder.body = value
	return builder
}

func (builder RequestBuilder) Build() *http.Request {
	request, _ := http.NewRequest(builder.method, builder.path, builder.body)
	return request
}

func NewRequestBuilder() RequestBuilder {
	return RequestBuilder{
		method: "GET",
		path:   "/default/path",
		body:   strings.NewReader("default body"),
	}
}
