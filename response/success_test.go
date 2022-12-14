package response_test

import (
	"net/http"
	"testing"

	"github.com/glocurrency/commons/response"
	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	code, resp := response.NewSuccessResponse("message")
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, response.SuccessResponse{Code: code, Message: "message"}, resp)

	code, resp = response.NewSuccessResponseCreated("message")
	assert.Equal(t, http.StatusCreated, code)
	assert.Equal(t, response.SuccessResponse{Code: code, Message: "message"}, resp)
}
