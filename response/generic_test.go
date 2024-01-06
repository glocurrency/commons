package response_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/glocurrency/commons/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type single struct{ Name string }
type many []single

func TestResponse(t *testing.T) {
	code, resp := response.NewResponse(single{Name: "test"})
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, response.Response[single]{single{Name: "test"}}, resp)

	got1, err := json.Marshal(resp)
	require.NoError(t, err)
	require.Equal(t, `{"data":{"Name":"test"}}`, string(got1))

	code, resp = response.NewResponseCreated(single{Name: "test"})
	assert.Equal(t, http.StatusCreated, code)
	assert.Equal(t, response.Response[single]{single{Name: "test"}}, resp)

	got2, err := json.Marshal(resp)
	require.NoError(t, err)
	require.Equal(t, `{"data":{"Name":"test"}}`, string(got2))
}

func TestManyPaginated(t *testing.T) {
	code, resp := response.NewManyResponsePaginated(many{single{Name: "a"}, single{Name: "b"}}, 2)
	require.Equal(t, http.StatusOK, code)

	got, err := json.Marshal(resp)
	require.NoError(t, err)
	require.Equal(t, `{"data":[{"Name":"a"},{"Name":"b"}],"pagination":{"total":2}}`, string(got))
}
