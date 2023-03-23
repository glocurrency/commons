package binding_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/binding"
	"github.com/glocurrency/commons/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logger.Log().SetOutput(io.Discard)
	gin.SetMode(gin.ReleaseMode)
	os.Exit(m.Run())
}

func TestParseParamUUID(t *testing.T) {
	tests := []struct {
		name       string
		paramName  string
		paramValue string
		wantErr    bool
	}{
		{"valid uuid", "id", uuid.NewString(), false},
		{"invalid uuid", "id", "invalid", true},
		{"invalid param name", "xyz", uuid.NewString(), true},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/"+test.paramValue, nil)

			w := httptest.NewRecorder()

			router := gin.New()
			router.GET("/:id", func(ctx *gin.Context) {
				got, err := binding.ParseParamUUID(ctx, test.paramName)
				if test.wantErr {
					assert.Error(t, err)
					return
				} else {
					assert.NoError(t, err)
					assert.Equal(t, test.paramValue, got.String())
				}
			})
			router.ServeHTTP(w, req)
		})
	}
}

func TestMustParseParamUUID(t *testing.T) {
	tests := []struct {
		name       string
		paramName  string
		paramValue string
		wantOk     bool
	}{
		{"valid uuid", "id", uuid.NewString(), true},
		{"invalid uuid", "id", "invalid", false},
		{"invalid param name", "xyz", uuid.NewString(), false},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/"+test.paramValue, nil)

			w := httptest.NewRecorder()

			router := gin.New()
			router.GET("/:id", func(ctx *gin.Context) {
				got, ok := binding.MustParseParamUUID(ctx, test.paramName)
				if test.wantOk {
					assert.True(t, ok)
					assert.Equal(t, test.paramValue, got.String())
				} else {
					assert.False(t, ok)
				}
			})
			router.ServeHTTP(w, req)
		})
	}
}
