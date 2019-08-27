package validation_test

import (
	validation "http-interaction-validation"
	"net/http"
	"net/http/httptest"
)

func whenHandlerIsCalled(request *http.Request, options ...validation.Option) (*httptest.ResponseRecorder, *http.Request) {
	var nextHandlerRequest *http.Request
	responseRecorder := httptest.NewRecorder()
	validationWrapper := validation.NewWrapper(options...)
	handler := func(w http.ResponseWriter, r *http.Request) {
		nextHandlerRequest = r
		w.WriteHeader(200)
	}
	validationWrapper(http.HandlerFunc(handler)).ServeHTTP(responseRecorder, request)
	return responseRecorder, nextHandlerRequest
}
