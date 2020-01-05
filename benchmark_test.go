package validation_test

import (
	validation "http-interaction-validation"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var response *httptest.ResponseRecorder

func BenchmarkWithoutValidation(b *testing.B) {
	b.ReportAllocs()
	request := requestBuilder.
		WithBody(strings.NewReader(defaultPayload)).
		Build()

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	for n := 0; n < b.N; n++ {
		handler.ServeHTTP(responseRecorder, request)
	}
}

func BenchmarkWithValidation(b *testing.B) {
	b.ReportAllocs()

	request := requestBuilder.
		WithBody(strings.NewReader(defaultPayload)).
		Build()
	responseRecorder := httptest.NewRecorder()
	validationWrapper := validation.NewWrapper(
		validation.RequestValidation(
			validation.Payload(&TestPayload{}),
			validation.Enabled(true),
		),
		validation.PreservePayload(false))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	wrappedHandler := validationWrapper(handler)

	for n := 0; n < b.N; n++ {
		wrappedHandler.ServeHTTP(responseRecorder, request)
	}
}
