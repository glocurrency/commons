package response_test

import (
	"net/http"
	"testing"

	"github.com/glocurrency/commons/response"
	"github.com/stretchr/testify/assert"
)

type single struct{ Name string }

func TestSingleResponse(t *testing.T) {
	code, resp := response.NewSingleResponse(single{Name: "test"})
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, response.SingleResponse[single]{single{Name: "test"}}, resp)

	code, resp = response.NewSingleResponseCreated(single{Name: "test"})
	assert.Equal(t, http.StatusCreated, code)
	assert.Equal(t, response.SingleResponse[single]{single{Name: "test"}}, resp)
}
