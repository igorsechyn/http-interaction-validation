package validation

import (
	"context"
	"encoding/json"
	"log"
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

type ValidationResponse struct {
	Code   string   `json:"code"`
	Errors []string `json:"errors"`
}

func (handler *ValidationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result := handler.validateRequest(r)
	if !result.isValid {
		w.WriteHeader(400)
		writeErrorResponse(w, result)
		return
	}

	ctx := context.WithValue(r.Context(), payloadKey, result.validatedValue)
	handler.next.ServeHTTP(w, r.WithContext(ctx))
}

func (handler *ValidationHandler) validateRequest(r *http.Request) bodyValidationResult {
	return handler.bodyValidator.validate(r)
}

func writeErrorResponse(w http.ResponseWriter, validationResult bodyValidationResult) {
	response := ValidationResponse{
		Code:   "body.validation.failure",
		Errors: validationResult.outcome.errors,
	}

	bytes, _ := json.Marshal(response)
	_, err := w.Write(bytes)
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}

func NewWrapper(options ...Option) func(handler http.Handler) http.Handler {
	config := getConfig(options...)
	validationHandler := &ValidationHandler{
		config:        config,
		bodyValidator: newBodyValidator(config),
	}
	return func(handler http.Handler) http.Handler {
		validationHandler.next = handler
		return validationHandler
	}
}
