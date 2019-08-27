package validation

type config struct {
	preservePayload         bool
	requestValidationConfig *requestValidationConfig
}

type requestValidationConfig struct {
	payloadValue interface{}
}

func getConfig(options ...Option) *config {
	handlerConfig := &config{
		preservePayload: true,
		requestValidationConfig: &requestValidationConfig{
			payloadValue: nil,
		},
	}

	for _, option := range options {
		option(handlerConfig)
	}

	return handlerConfig
}

type Option func(config *config)

func PreservePayload(value bool) Option {
	return func(config *config) {
		config.preservePayload = value
	}
}

type RequestValidationConfigOption func(requestConfig *requestValidationConfig)

func Payload(payloadValue interface{}) RequestValidationConfigOption {
	return func(config *requestValidationConfig) {
		config.payloadValue = payloadValue
	}
}

func RequestValidation(options ...RequestValidationConfigOption) Option {
	requestValidationConfig := &requestValidationConfig{
		payloadValue: nil,
	}

	for _, option := range options {
		option(requestValidationConfig)
	}

	return func(config *config) {
		config.requestValidationConfig = requestValidationConfig
	}
}
