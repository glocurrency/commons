package translator_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	ginbinging "github.com/gin-gonic/gin/binding"
	"github.com/glocurrency/commons/logger"
	"github.com/glocurrency/commons/translator"
	locale "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logger.Log().SetOutput(io.Discard)
	gin.SetMode(gin.ReleaseMode)
	os.Exit(m.Run())
}

func TestTranslator(t *testing.T) {
	t.Parallel()

	uni := ut.New(locale.New())
	fallback := uni.GetFallback()

	tests := []struct {
		name           string
		bindTranslator bool
	}{
		{"no translator", false},
		{"has translator", true},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			router := gin.New()

			if test.bindTranslator {
				if v, ok := ginbinging.Validator.Engine().(*validator.Validate); ok {
					en.RegisterDefaultTranslations(v, fallback)
					router.Use(translator.SetTranslatorMiddleware(fallback))
				}
			}

			router.GET("/", func(ctx *gin.Context) {
				got, ok := translator.GetTranslator(ctx)
				if test.bindTranslator {
					assert.True(t, ok)
					assert.Same(t, fallback, got)
				} else {
					assert.False(t, ok)
				}
			})
			router.ServeHTTP(w, req)
		})
	}
}
