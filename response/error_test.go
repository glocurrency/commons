package response_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/glocurrency/commons/response"
	"github.com/stretchr/testify/assert"
)

func TestErrResponse(t *testing.T) {
	code, resp := response.NewErrResponse(errors.New("i am an error"))
	assert.Equal(t, http.StatusInternalServerError, code)
	assert.Equal(t, response.ErrResponse{Code: code, Message: "i am an error"}, resp)

	code, resp = response.NewErrResponseForbidden("message")
	assert.Equal(t, http.StatusForbidden, code)
	assert.Equal(t, response.ErrResponse{Code: code, Message: "message"}, resp)

	code, resp = response.NewErrResponseNotFound("message")
	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, response.ErrResponse{Code: code, Message: "message"}, resp)

	code, resp = response.NewErrResponseBadRequest("message")
	assert.Equal(t, http.StatusBadRequest, code)
	assert.Equal(t, response.ErrResponse{Code: code, Message: "message"}, resp)

	code, resp = response.NewErrResponseException("message")
	assert.Equal(t, http.StatusInternalServerError, code)
	assert.Equal(t, response.ErrResponse{Code: code, Message: "message"}, resp)

	errors := map[string]string{"key": "value"}

	code, resp = response.NewErrResponseValidationErrors("message", errors)
	assert.Equal(t, http.StatusBadRequest, code)
	assert.Equal(t, response.ErrResponse{Code: code, Message: "message", Errors: errors}, resp)

	assert.True(t, resp.HasErrorField("key"))
	assert.True(t, resp.HasErrorFields([]string{"key"}))
	assert.False(t, resp.HasErrorField("foo"))
	assert.False(t, resp.HasErrorFields([]string{"foo"}))
}
