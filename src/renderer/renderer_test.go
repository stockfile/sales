package renderer

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	ID   string
	Name string
}

func TestRenderServerError(t *testing.T) {
	response, _ := RenderServerError("some server error")

	var responseError ResponseError

	if err := json.Unmarshal([]byte(response.Body), &responseError); err != nil {
		t.Error("Can't parse response body")
	}

	assert.Equal(t, 500, response.StatusCode, "status code is 500")
	assert.Equal(t, "{\"error\":\"some server error\"}", response.Body, "error message with error key")
	assert.Equal(t, "some server error", responseError.Error, "body should be the passed body parameter")
}

func TestRenderClientError(t *testing.T) {
	response, _ := RenderClientError("not found", 404)

	var responseError ResponseError

	if err := json.Unmarshal([]byte(response.Body), &responseError); err != nil {
		t.Error("Can't parse response body")
	}

	assert.Equal(t, 404, response.StatusCode, "status code is 404")
	assert.Equal(t, "{\"error\":\"not found\"}", response.Body, "error message with error key")
	assert.Equal(t, "not found", responseError.Error, "body should be the passed body parameter")
}

func TestRenderSuccess(t *testing.T) {
	someStruct := &TestStruct{
		ID:   "1",
		Name: "Julian",
	}

	body, _ := json.Marshal(someStruct)
	response, _ := RenderSuccess(body)

	assert.Equal(t, 200, response.StatusCode, "status code is 200")
	assert.Equal(t, "{\"ID\":\"1\",\"Name\":\"Julian\"}", response.Body, "body is serialised properly")
}
