package validation_test

import (
	validation "http-interaction-validation"
	"http-interaction-validation/test_support/builders"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPayload struct {
	Name string `json:"name" jsonschema:"required,title=name"`
}

var requestBuilder = builders.NewRequestBuilder()
var defaultPayload = `{"name":"value"}`

func TestHandler_PayloadValidation(t *testing.T) {

	t.Run("it should copy request payload for validation, if not configured otherwise", func(t *testing.T) {
		request := requestBuilder.
			WithBody(strings.NewReader(defaultPayload)).
			Build()
		_, nextHandlerRequest := whenHandlerIsCalled(
			request,
			validation.RequestValidation(validation.Payload(&TestPayload{})),
		)

		body, err := ioutil.ReadAll(nextHandlerRequest.Body)
		assert.NoError(t, err)
		assert.Equal(t, defaultPayload, string(body), "request should preserve body reader")
	})

	t.Run("it should add body payload into request context", func(t *testing.T) {
		request := requestBuilder.WithBody(strings.NewReader(defaultPayload)).Build()
		_, nextHandlerRequest := whenHandlerIsCalled(
			request,
			validation.RequestValidation(validation.Payload(&TestPayload{})),
		)

		bytes, ok := validation.PayloadFromContext(nextHandlerRequest.Context())
		assert.True(t, ok)
		assert.Equal(t, defaultPayload, string(bytes))
	})

	t.Run("it should not keep body payload on request, when PreservePayload option is false", func(t *testing.T) {
		request := requestBuilder.WithBody(strings.NewReader(defaultPayload)).Build()
		_, nextHandlerRequest := whenHandlerIsCalled(
			request,
			validation.RequestValidation(validation.Payload(&TestPayload{})),
			validation.PreservePayload(false),
		)

		body, err := ioutil.ReadAll(nextHandlerRequest.Body)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(body), "request payload should be empty")
	})

	t.Run("it should not read the body from request, if request payload validation is not configured", func(t *testing.T) {
		request := requestBuilder.WithBody(strings.NewReader(defaultPayload)).Build()
		_, nextHandlerRequest := whenHandlerIsCalled(
			request,
			validation.PreservePayload(false),
		)

		bytesFromContext, ok := validation.PayloadFromContext(nextHandlerRequest.Context())
		assert.True(t, ok)
		assert.Nil(t, bytesFromContext, "body value in context should be nil")
		body, err := ioutil.ReadAll(nextHandlerRequest.Body)
		assert.NoError(t, err)
		assert.Greater(t, len(body), 0, "request payload should not be empty")
	})

	t.Run("it should return a 400 response, if payload validation fails", func(t *testing.T) {
		request := requestBuilder.WithBody(strings.NewReader("{}")).Build()
		responseRecorder, _ := whenHandlerIsCalled(
			request,
			validation.RequestValidation(validation.Payload(&TestPayload{})),
		)

		response := responseRecorder.Result()
		assert.Equal(t, 400, response.StatusCode)
	})
}
