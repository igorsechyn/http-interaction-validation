package e2e_test

import (
	"bytes"
	"context"
	validation "http-interaction-validation"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type E2eTestPayload struct {
	Name string `json:"name" jsonschema:"required,title=name"`
}

func TestMain(m *testing.M) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		w.Write(bytes)
	}
	wrapper := validation.NewWrapper(
		validation.RequestValidation(validation.Payload(&E2eTestPayload{})),
	)
	server := startServer(wrapper(http.HandlerFunc(handler)))
	os.Exit(m.Run())
	server.Shutdown(context.Background())
}
func TestValidation_Ok(t *testing.T) {
	payload := []byte(`{"name":"me"}`)
	response, err := http.Post("http://localhost:8080/", "application/json", bytes.NewBuffer(payload))

	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
}

func TestValidation_NotOk(t *testing.T) {
	payload := []byte(`{"wrong":"format"}`)
	response, err := http.Post("http://localhost:8080/", "application/json", bytes.NewBuffer(payload))

	assert.NoError(t, err)
	assert.Equal(t, 400, response.StatusCode)
	responseBody, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, `{"code":"body.validation.failure","errors":["(root): name is required","(root): Additional property wrong is not allowed"]}`, string(responseBody))
}
