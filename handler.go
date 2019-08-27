package validation

import (
	"context"
	"net/http"
)

type payloadKeyType struct{}

var payloadKey payloadKeyType

func PayloadFromContext(ctx context.Context) ([]byte, bool) {
	val, ok := ctx.Value(payloadKey).([]byte)
	return val, ok
}

type ValidationHandler struct {
	next          http.Handler
	config        *config
	bodyValidator *bodyValidator
}

func (handler *ValidationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result := handler.bodyValidator.validate(r)
	if !result.isValid {
		w.WriteHeader(400)
		return
	}

	ctx := context.WithValue(r.Context(), payloadKey, result.validatedValue)
	handler.next.ServeHTTP(w, r.WithContext(ctx))
}

func NewWrapper(options ...Option) func(handler http.Handler) http.Handler {
	config := getConfig(options...)
	validationHandler := &ValidationHandler{
		config: config,
		bodyValidator: &bodyValidator{
			config: config,
		},
	}
	return func(handler http.Handler) http.Handler {
		validationHandler.next = handler
		return validationHandler
	}
}
