package translator_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	ginbinging "github.com/gin-gonic/gin/binding"
	"github.com/glocurrency/commons/binding"
	"github.com/glocurrency/commons/translator"
	"github.com/glocurrency/commons/validation"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name    string `json:"name" binding:"omitempty,alphanumspace"`
	Company string `json:"company" binding:"omitempty,alphanumspacedash"`
	Bank    string `json:"bank" binding:"omitempty,banksupported"`
	BIC     string `json:"bic" binding:"omitempty,bic"`
	Country string `json:"country" binding:"omitempty,iso3166_1_alpha2"`
}

func Test_RegisterTranslatorFor(t *testing.T) {
	tests := []struct {
		name         string
		payload      string
		wantStatus   int
		wantResponse string
	}{
		{"alphanumspace", `{"name": "%"}`, 400, "name can only contain alphanumeric characters and spaces"},
		{"alphanumspacedash", `{"company": "%"}`, 400, "company can only contain alphanumeric characters, spaces and dashes"},
		{"banksupported", `{"bank": "%"}`, 400, "bank can only contain bank suported characters"},
		{"bic", `{"bic": "123"}`, 400, "bic must comply with BIC format"},
		{"country", `{"country": "XX"}`, 400, "country must be a valid country code"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.payload))

			router := gin.New()
			if v, ok := ginbinging.Validator.Engine().(*validator.Validate); ok {
				v.RegisterValidation("alphanumspace", validation.ValidateAlphaNumSpace)
				v.RegisterValidation("alphanumspacedash", validation.ValidateAlphaNumSpaceDash)
				v.RegisterValidation("banksupported", validation.ValidateBankSupported)

				t := translator.RegisterTranslatorFor(v)
				router.Use(translator.SetTranslatorMiddleware(t))
			}

			router.POST("/", func(ctx *gin.Context) {
				var testStruct TestStruct
				if !binding.MustDecodeBody(ctx, &testStruct) {
					return
				}
				ctx.Status(http.StatusOK)
			})
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.wantResponse)
		})
	}
}
