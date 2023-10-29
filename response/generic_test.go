package response_test

import (
	"net/http"
	"testing"

	"github.com/glocurrency/commons/response"
	"github.com/stretchr/testify/assert"
)

type single struct{ Name string }

func TestResponse(t *testing.T) {
	code, resp := response.NewResponse(single{Name: "test"})
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, response.Response[single]{single{Name: "test"}}, resp)

	code, resp = response.NewResponseCreated(single{Name: "test"})
	assert.Equal(t, http.StatusCreated, code)
	assert.Equal(t, response.Response[single]{single{Name: "test"}}, resp)
}
