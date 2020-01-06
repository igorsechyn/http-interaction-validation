package validation

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alecthomas/jsonschema"
	js "github.com/qri-io/jsonschema"
)

type bodyValidator struct {
	config            *config
	bodyPayloadSchema *js.RootSchema
}

type outcome struct {
	errors []string
}

type bodyValidationResult struct {
	validatedValue []byte
	isValid        bool
	outcome        outcome
}

func initBodyPayloadSchema(config *config) *js.RootSchema {
	schema := new(js.RootSchema)
	if config.requestValidationConfig.payloadValue != nil {
		reflector := &jsonschema.Reflector{
			AllowAdditionalProperties: config.requestValidationConfig.additionalProperties,
		}
		jsonSchema := reflector.Reflect(config.requestValidationConfig.payloadValue)
		jsonSchemaBytes, err := json.Marshal(jsonSchema)
		if err != nil {
			fmt.Printf("Could not marshall json schema %v. Err: %v", string(jsonSchemaBytes), err)
		}
		err = json.Unmarshal(jsonSchemaBytes, schema)
		if err != nil {
			fmt.Printf("Could not unmarshall json schema %v. Err: %v", string(jsonSchemaBytes), err)
		}
	}

	return schema
}

func newBodyValidator(config *config) *bodyValidator {
	return &bodyValidator{
		config:            config,
		bodyPayloadSchema: initBodyPayloadSchema(config),
	}
}

func (validator *bodyValidator) validate(r *http.Request) bodyValidationResult {
	if validator.shouldValidateBody() {
		bodyValue, _ := readPayload(r, validator.config.requestValidationConfig.preservePayload)
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
	schemaValidationResult, err := validator.bodyPayloadSchema.ValidateBytes(payload)
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
		isValid:        len(schemaValidationResult) == 0,
		outcome:        toOutcome(schemaValidationResult),
	}
}

func (validator *bodyValidator) shouldValidateBody() bool {
	return (validator.config.requestValidationConfig.payloadValue != nil &&
		validator.config.requestValidationConfig.enabled)
}

func toOutcome(validationErrors []js.ValError) outcome {
	errors := make([]string, len(validationErrors))
	for index, error := range validationErrors {
		errors[index] = error.Error()
	}
	return outcome{
		errors: errors,
	}
}
