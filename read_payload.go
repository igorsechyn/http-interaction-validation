package validation

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

func copyRequestPayload(request *http.Request) ([]byte, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(request.Body, &buf)
	bodyValue, err := ioutil.ReadAll(tee)
	request.Body = ioutil.NopCloser(&buf)
	return bodyValue, err
}

func readPayload(request *http.Request, preserveBodyOnRequest bool) ([]byte, error) {
	if preserveBodyOnRequest {
		return copyRequestPayload(request)
	}

	return ioutil.ReadAll(request.Body)
}
