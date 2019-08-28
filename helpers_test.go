package validation_test

import (
	validation "http-interaction-validation"
	"net/http"
	"net/http/httptest"
)

func whenWrappedHandlerIsCalled(request *http.Request, options ...validation.Option) (*httptest.ResponseRecorder, *http.Request) {
	var wrappedHandlerRequest *http.Request
	responseRecorder := httptest.NewRecorder()
	validationWrapper := validation.NewWrapper(options...)
	handler := func(w http.ResponseWriter, r *http.Request) {
		wrappedHandlerRequest = r
		w.WriteHeader(200)
	}
	validationWrapper(http.HandlerFunc(handler)).ServeHTTP(responseRecorder, request)
	return responseRecorder, wrappedHandlerRequest
}
