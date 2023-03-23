package binding_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	ginbinging "github.com/gin-gonic/gin/binding"
	"github.com/glocurrency/commons/binding"
	"github.com/glocurrency/commons/logger"
	"github.com/glocurrency/commons/translator"
	locale "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
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
			router.ServeHTTP(httptest.NewRecorder(), req)
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

			router := gin.New()
			router.GET("/:id", func(ctx *gin.Context) {
				got, ok := binding.MustParseParamUUID(ctx, test.paramName)
				if test.wantOk {
					assert.True(t, ok)
					assert.Equal(t, test.paramValue, got.String())
				} else {
					assert.False(t, ok)
					return
				}
				ctx.Status(http.StatusOK)
			})
			router.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}

func TestMustDecodeBody_CanDecode(t *testing.T) {
	type testStruct struct {
		Name string `json:"name" binding:"required"`
	}

	tests := []struct {
		name   string
		body   string
		wantOk bool
	}{
		{"valid json", `{"name": "John"}`, true},
		{"missing property", `{}`, false},
		{"invalid json", `not a json`, false},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			w := httptest.NewRecorder()

			router := gin.New()
			router.POST("/", func(ctx *gin.Context) {
				var got testStruct

				ok := binding.MustDecodeBody(ctx, &got)
				assert.Equal(t, test.wantOk, ok)
			})
			router.ServeHTTP(w, req)
		})
	}
}

func TestMustDecodeBody_CanTranslate(t *testing.T) {
	type testStruct struct {
		Name string `json:"name" binding:"required"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name          string
		hasTranslator bool
		wantMessage   string
	}{
		{"no translator", false, `{"code":400,"message":"Request data invalid"}`},
		{"has translator", true, `{"code":400,"message":"Request data invalid","errors":{"Name":"Name is a required field"}}`},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"age": 20}`))
			w := httptest.NewRecorder()

			router := gin.New()

			if test.hasTranslator {
				if v, ok := ginbinging.Validator.Engine().(*validator.Validate); ok {
					uni := ut.New(locale.New())
					fallback := uni.GetFallback()
					en.RegisterDefaultTranslations(v, fallback)

					router.Use(translator.SetTranslatorMiddleware(fallback))
				}
			}

			router.POST("/", func(ctx *gin.Context) {
				var got testStruct

				ok := binding.MustDecodeBody(ctx, &got)
				assert.Equal(t, false, ok)
			})
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Equal(t, test.wantMessage, w.Body.String())
		})
	}
}
