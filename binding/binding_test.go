package binding_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	ginbinging "github.com/gin-gonic/gin/binding"
	"github.com/glocurrency/commons/binding"
	"github.com/glocurrency/commons/router"
	"github.com/glocurrency/commons/translator"
	locale "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

			router := router.NewRouterWithValidation()
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

			router := router.NewRouterWithValidation()
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

			router := router.NewRouterWithValidation()
			router.POST("/", func(ctx *gin.Context) {
				var got testStruct

				ok := binding.MustDecodeBody(ctx, &got)
				assert.Equal(t, test.wantOk, ok)
			})
			router.ServeHTTP(w, req)
		})
	}
}

func TestMustDecodeBody_CanDecodeTwice(t *testing.T) {
	type testStruct struct {
		Name string `json:"name" binding:"required"`
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "John"}`))
	w := httptest.NewRecorder()

	router := router.NewRouterWithValidation()
	router.POST("/", func(ctx *gin.Context) {
		var got1 testStruct
		require.True(t, binding.MustDecodeBody(ctx, &got1))
		require.Equal(t, "John", got1.Name)

		var got2 testStruct
		require.True(t, binding.MustDecodeBody(ctx, &got2))
		require.Equal(t, "John", got1.Name)
	})
	router.ServeHTTP(w, req)
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
