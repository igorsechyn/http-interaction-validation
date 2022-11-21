package validation

type config struct {
	requestValidationConfig *requestValidationConfig
}

type requestValidationConfig struct {
	preservePayload      bool
	payloadValue         interface{}
	bodyRequired         bool
	enabled              bool
	additionalProperties bool
}

func getConfig(options ...Option) *config {
	handlerConfig := &config{
		requestValidationConfig: &requestValidationConfig{
			preservePayload:      true,
			payloadValue:         nil,
			bodyRequired:         true,
			enabled:              true,
			additionalProperties: true,
		},
	}

	for _, option := range options {
		option(handlerConfig)
	}

	return handlerConfig
}

type Option func(config *config)

type RequestValidationConfigOption func(requestConfig *requestValidationConfig)

func PreservePayload(value bool) RequestValidationConfigOption {
	return func(config *requestValidationConfig) {
		config.preservePayload = value
	}
}

func Payload(payloadValue interface{}) RequestValidationConfigOption {
	return func(config *requestValidationConfig) {
		config.payloadValue = payloadValue
	}
}

func BodyRequired(value bool) RequestValidationConfigOption {
	return func(config *requestValidationConfig) {
		config.bodyRequired = value
	}
}

func AdditionalProperties(value bool) RequestValidationConfigOption {
	return func(config *requestValidationConfig) {
		config.additionalProperties = value
	}
}

func Enabled(value bool) RequestValidationConfigOption {
	return func(config *requestValidationConfig) {
		config.enabled = value
	}
}

func RequestValidation(options ...RequestValidationConfigOption) Option {
	return func(config *config) {
		for _, option := range options {
			option(config.requestValidationConfig)
		}
	}
}
