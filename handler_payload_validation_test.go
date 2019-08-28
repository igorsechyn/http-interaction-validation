package validation_test

import (
	validation "http-interaction-validation"
	"http-interaction-validation/test_support/builders"
	"io"
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
		_, wrappedHandlerRequest := whenWrappedHandlerIsCalled(
			request,
			validation.RequestValidation(validation.Payload(&TestPayload{})),
		)

		body, err := ioutil.ReadAll(wrappedHandlerRequest.Body)
		assert.NoError(t, err)
		assert.Equal(t, defaultPayload, string(body), "request should preserve body reader")
	})

	t.Run("it should add body payload into request context", func(t *testing.T) {
		request := requestBuilder.WithBody(strings.NewReader(defaultPayload)).Build()
		_, wrappedHandlerRequest := whenWrappedHandlerIsCalled(
			request,
			validation.RequestValidation(validation.Payload(&TestPayload{})),
		)

		bytes, ok := validation.PayloadFromContext(wrappedHandlerRequest.Context())
		assert.True(t, ok)
		assert.Equal(t, defaultPayload, string(bytes))
	})

	t.Run("it should not keep body payload on request, when PreservePayload option is false", func(t *testing.T) {
		request := requestBuilder.WithBody(strings.NewReader(defaultPayload)).Build()
		_, wrappedHandlerRequest := whenWrappedHandlerIsCalled(
			request,
			validation.RequestValidation(validation.Payload(&TestPayload{})),
			validation.PreservePayload(false),
		)

		body, err := ioutil.ReadAll(wrappedHandlerRequest.Body)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(body), "request payload should be empty")
	})

	t.Run("it should not read the body from request, if request payload validation is not configured", func(t *testing.T) {
		request := requestBuilder.WithBody(strings.NewReader(defaultPayload)).Build()
		_, wrappedHandlerRequest := whenWrappedHandlerIsCalled(
			request,
			validation.PreservePayload(false),
		)

		bytesFromContext, ok := validation.PayloadFromContext(wrappedHandlerRequest.Context())
		assert.True(t, ok)
		assert.Nil(t, bytesFromContext, "body value in context should be nil")
		body, err := ioutil.ReadAll(wrappedHandlerRequest.Body)
		assert.NoError(t, err)
		assert.Greater(t, len(body), 0, "request payload should not be empty")
	})

	t.Run("it should validate body payload based on provided options", func(t *testing.T) {
		testCases := []struct {
			description        string
			options            []validation.Option
			payload            io.Reader
			expectedStatusCode int
			expectedResponse   string
		}{
			{
				description:        "body has wrong format",
				options:            []validation.Option{validation.RequestValidation(validation.Payload(&TestPayload{}))},
				payload:            strings.NewReader("{}"),
				expectedStatusCode: 400,
				expectedResponse:   `{"code":"body.validation.failure","errors":["(root): name is required"]}`,
			},
			{
				description: "body has wrong format, and request validation is disabled",
				options: []validation.Option{validation.RequestValidation(
					validation.Payload(&TestPayload{}), validation.Enabled(false)),
				},
				payload:            strings.NewReader("{}"),
				expectedStatusCode: 200,
				expectedResponse:   "",
			},
			{
				description:        "body is malformed",
				options:            []validation.Option{validation.RequestValidation(validation.Payload(&TestPayload{}))},
				payload:            strings.NewReader("{not json"),
				expectedStatusCode: 400,
				expectedResponse:   `{"code":"body.validation.failure","errors":["invalid character 'n' looking for beginning of object key string"]}`,
			},
			{
				description: "body is missing, but is required",
				options: []validation.Option{validation.RequestValidation(
					validation.Payload(&TestPayload{})),
				},
				payload:            nil,
				expectedStatusCode: 400,
				expectedResponse:   `{"code":"body.validation.failure","errors":["body is missing, but is required"]}`,
			},
			{
				description: "body is missing, but is not required",
				options: []validation.Option{validation.RequestValidation(
					validation.Payload(&TestPayload{}),
					validation.BodyRequired(false)),
				},
				payload:            nil,
				expectedStatusCode: 200,
				expectedResponse:   "",
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.description, func(t *testing.T) {
				request := requestBuilder.WithBody(testCase.payload).Build()
				responseRecorder, _ := whenWrappedHandlerIsCalled(
					request,
					testCase.options...,
				)

				response := responseRecorder.Result()
				assert.Equal(t, testCase.expectedStatusCode, response.StatusCode)
				bytes, _ := ioutil.ReadAll(response.Body)
				assert.Equal(t, string(bytes), testCase.expectedResponse)
			})
		}
	})
}
