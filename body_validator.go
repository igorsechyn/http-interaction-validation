package validation

import (
	"net/http"

	"github.com/alecthomas/jsonschema"
	"github.com/xeipuuv/gojsonschema"
)

type bodyValidator struct {
	config                  *config
	bodyPayloadSchemaLoader gojsonschema.JSONLoader
}

type outcome struct {
	errors []string
}

type bodyValidationResult struct {
	validatedValue []byte
	isValid        bool
	outcome        outcome
}

func initBodyPayloadSchemaLoader(config *config) gojsonschema.JSONLoader {
	var schemaLoader gojsonschema.JSONLoader
	if config.requestValidationConfig.payloadValue != nil {
		reflector := &jsonschema.Reflector{
			AllowAdditionalProperties: config.requestValidationConfig.additionalProperties,
		}
		jsonSchema := reflector.Reflect(config.requestValidationConfig.payloadValue)
		schemaLoader = gojsonschema.NewGoLoader(jsonSchema)
	}

	return schemaLoader
}

func newBodyValidator(config *config) *bodyValidator {
	return &bodyValidator{
		config:                  config,
		bodyPayloadSchemaLoader: initBodyPayloadSchemaLoader(config),
	}
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
	dataLoader := gojsonschema.NewStringLoader((string(payload)))
	schemaValidationResult, err := gojsonschema.Validate(validator.bodyPayloadSchemaLoader, dataLoader)

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
	return (validator.config.requestValidationConfig.payloadValue != nil &&
		validator.config.requestValidationConfig.enabled)
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
