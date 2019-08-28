package validation

import (
	"net/http"

	"github.com/alecthomas/jsonschema"
	"github.com/xeipuuv/gojsonschema"
)

type bodyValidator struct {
	config *config
}

type outcome struct {
	errors []string
}

type bodyValidationResult struct {
	validatedValue []byte
	isValid        bool
	outcome        outcome
}

func (validator *bodyValidator) validate(r *http.Request) bodyValidationResult {
	if validator.shouldValidateBody() {
		bodyValue, _ := readPayload(r, validator.config.preservePayload)
		return validator.validatePayload(bodyValue)
	}

	return bodyValidationResult{
		isValid:        true,
		validatedValue: nil,
		outcome:        outcome{},
	}
}

func (validator *bodyValidator) validatePayload(payload []byte) bodyValidationResult {
	if payload == nil {
		return validator.validateNilBody()
	}

	return validator.validateAgainstJsonSchema(payload)
}

func (validator *bodyValidator) validateNilBody() bodyValidationResult {
	if !validator.config.requestValidationConfig.bodyRequired {
		return bodyValidationResult{
			validatedValue: nil,
			isValid:        true,
			outcome:        outcome{},
		}
	}
	return bodyValidationResult{
		validatedValue: nil,
		isValid:        false,
		outcome: outcome{
			errors: []string{"body is missing, but is required"},
		},
	}
}

func (validator *bodyValidator) validateAgainstJsonSchema(payload []byte) bodyValidationResult {
	jsonSchema := jsonschema.Reflect(validator.config.requestValidationConfig.payloadValue)
	schemaLoader := gojsonschema.NewGoLoader(jsonSchema)
	dataLoader := gojsonschema.NewStringLoader((string(payload)))
	schemaValidationResult, err := gojsonschema.Validate(schemaLoader, dataLoader)

	if err != nil {
		return bodyValidationResult{
			validatedValue: payload,
			isValid:        false,
			outcome: outcome{
				errors: []string{err.Error()},
			},
		}
	}

	return bodyValidationResult{
		validatedValue: payload,
		isValid:        schemaValidationResult.Valid(),
		outcome:        toOutcome(schemaValidationResult),
	}
}

func (validator *bodyValidator) shouldValidateBody() bool {
	return validator.config.requestValidationConfig.payloadValue != nil
}

func toOutcome(result *gojsonschema.Result) outcome {
	errors := make([]string, len(result.Errors()))
	for index, error := range result.Errors() {
		errors[index] = error.String()
	}
	return outcome{
		errors: errors,
	}
}
