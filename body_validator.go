package validation

import (
	"net/http"

	"github.com/alecthomas/jsonschema"
	"github.com/xeipuuv/gojsonschema"
)

type bodyValidator struct {
	config *config
}

type validationResult struct {
	validatedValue []byte
	isValid        bool
	outcome        *gojsonschema.Result
}

func (validator *bodyValidator) validate(r *http.Request) validationResult {
	if validator.shouldValidateBody() {
		bodyValue, _ := readPayload(r, validator.config.preservePayload)
		result, _ := validator.validatePayload(bodyValue)
		return validationResult{
			validatedValue: bodyValue,
			isValid:        result.Valid(),
			outcome:        result,
		}
	}

	return validationResult{
		isValid:        true,
		validatedValue: nil,
		outcome:        nil,
	}
}

func (validator *bodyValidator) validatePayload(payload []byte) (*gojsonschema.Result, error) {
	jsonSchema := jsonschema.Reflect(validator.config.requestValidationConfig.payloadValue)
	schemaLoader := gojsonschema.NewGoLoader(jsonSchema)
	dataLoader := gojsonschema.NewStringLoader((string(payload)))
	return gojsonschema.Validate(schemaLoader, dataLoader)
}

func (validator *bodyValidator) shouldValidateBody() bool {
	return validator.config.requestValidationConfig.payloadValue != nil
}
