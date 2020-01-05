# http-interaction-validation
> Provides a wrapper around standard http handler to perform validation

## Motivation

There are a number of libraries that can be used for request validation, like https://github.com/thedevsaddam/govalidator and https://github.com/go-playground/validator. Both of these packages use a custom set of validation rules, either defined in a map separately from the struct describing the payload or in the field tags.

This library uses https://github.com/alecthomas/jsonschema package, which allows the user to define the validation rules as field tags in json schema format. The json schema can then be used to validate request payload, path parameters and headers. 

Another benefit of having a json schema describing requests and responses is the ability to autogenerate open api specification, similar to what https://github.com/go-chi/docgen is doing.

## Requirements
- go 1.12 or higher

## Installation

```
go get github.com/igorsechyn/http-interaction-validation
```

## Usage

Package provides a factory function to create a wrapper for a standard `http.Handler`, which will perform configured validations on request and call the original wrapper, if validation was successful. 

Validation is configured with [functional options](https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html). If no options are provided, the wrapper will just call through to the handler.

```
import (
	validation "http-interaction-validation"
	"net/http"
)

handler := func(w http.ResponseWriter, r *http.Request) {
    bytes, _ := ioutil.ReadAll(r.Body)
    w.Write(bytes)
}

withValidation := validation.NewWrapper()

wrappedHandler := withValidation(handler)
```

### Configuration options

#### `validation.PreservePayload(value bool)`

In order to validate request body, the payload has to be read from the request. This option can be used to control whether the `request.Body` should be copied for validation or not. Per default the body is copied and the wrapped handler can reread the payload. 

#### `validation.RequestValidation(options ...validation.RequestValidationConfigOption)`

Is used for setting up the request validation:

* `validation.Payload(payloadValue interface{})`: request payload validation is performed using json schema. In order to define a schema for handlers payload `jsonschema` field tags are used on the struct defining the payload. [github.com/alecthomas/jsonschema](https://github.com/alecthomas/jsonschema) package is used to create a valid json schema from the tags and `github.com/xeipuuv/gojsonschema` package to validate the payload against it.
```
type Payload struct {
    Name    string `json:"name" jsonschema:"required,minLength=1,maxLength=20,description=this is a property,title=the name"`
}
withValidation := validation.NewWrapper(
    validation.RequestValidation(validation.Payload(&Payload{})),
)
```

* `validation.BodyRequired(value bool)`: defines whether the body is required or not. Default is true and the validation will return an error response, if the body is missing

* `validation.AdditionalProperties(value bool)`: defines whether additionalProperties are allowed or not. Default is true

* `validation.Enabled(value bool)`: defines whether validation should be performed or skipped. Default is true, but can be used for example to skip validation in production environment

#### Error handling

If request validation fails, the wrapper will not call wrapped handler and send back a 400 response with the following payload

```
type ValidationResponse struct {
	Code   string   `json:"code"`
	Errors []string `json:"errors"`
}
```

see `test/e2e/e2e_test.go` for an example

## Performance

### Benchmark results

Baseline is a simple handler writing a 200 response without reading the body. It is compared to a wrapped handler which performs a body payload validation

```
go test -run=XXX -bench=.
goos: darwin
goarch: amd64
pkg: http-interaction-validation
BenchmarkWithoutValidation-4    326973954                3.69 ns/op            0 B/op          0 allocs/op
BenchmarkWithValidation-4          55587             20786 ns/op           11171 B/op        142 allocs/op
```

## Contributing
See [CONTRIBUTING.md](CONTRIBUTING.md)
