package translator_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	ginbinging "github.com/gin-gonic/gin/binding"
	"github.com/glocurrency/commons/binding"
	"github.com/glocurrency/commons/router"
	"github.com/glocurrency/commons/translator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name    string `json:"name" binding:"omitempty,alphanumspace"`
	Company string `json:"company" binding:"omitempty,alphanumspacedash"`
	Bank    string `json:"bank" binding:"omitempty,banksupported"`
	BIC     string `json:"bic" binding:"omitempty,bic"`
	Country string `json:"country" binding:"omitempty,iso3166_1_alpha2"`
	Age     string `json:"age" binding:"omitempty,18yo,notold"`
}

// Added struct to test tag extraction and default validations
type EdgeCaseStruct struct {
	Email  string `json:"email_address,omitempty" binding:"omitempty,email"` // Tests comma splitting + default translation
	Hidden string `json:"-" binding:"required"`                              // Tests json:"-" omission
	NoTag  string `binding:"required"`                                       // Tests fallback when no json tag exists
}

func Test_RegisterTranslatorFor(t *testing.T) {
	now := time.Now()
	seventeenYearsAgo := now.AddDate(-17, 0, 0)
	eighteenYearsAgo := now.AddDate(-18, 0, 0)
	oneHundredYearsAgo := now.AddDate(-100, 0, 0)
	twoHundredYearsAgo := now.AddDate(-200, 0, 0)

	tests := []struct {
		name         string
		payload      string
		wantStatus   int
		wantResponse string
	}{
		{"alphanumspace", `{"name": "%"}`, 400, "name can only contain alphanumeric characters and spaces"},
		{"alphanumspacedash", `{"company": "%"}`, 400, "company can only contain alphanumeric characters, spaces and dashes"},
		{"banksupported", `{"bank": "%"}`, 400, "bank can only contain bank supported characters"},
		{"bic", `{"bic": "123"}`, 400, "bic must comply with BIC format"},
		{"country", `{"country": "XX"}`, 400, "country must be a valid country code"},
		{"17 years old", fmt.Sprintf(`{"age": "%s"}`, seventeenYearsAgo.Format(time.DateOnly)), 400, "age should be at least 18 years old"},
		{"18 years old", fmt.Sprintf(`{"age": "%s"}`, eighteenYearsAgo.Format(time.DateOnly)), 200, ""},
		{"100 years old", fmt.Sprintf(`{"age": "%s"}`, oneHundredYearsAgo.Format(time.DateOnly)), 200, ""},
		{"200 years old", fmt.Sprintf(`{"age": "%s"}`, twoHundredYearsAgo.Format(time.DateOnly)), 400, "age should be a valid age"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.payload))

			r := router.NewRouterWithValidation()

			if v, ok := ginbinging.Validator.Engine().(*validator.Validate); ok {
				tr := translator.RegisterTranslatorFor(v)
				r.Use(translator.SetTranslatorMiddleware(tr))
			}

			r.POST("/", func(ctx *gin.Context) {
				var testStruct TestStruct
				if !binding.MustDecodeBody(ctx, &testStruct) {
					return
				}
				ctx.Status(http.StatusOK)
			})
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantResponse != "" {
				assert.Contains(t, w.Body.String(), tt.wantResponse)
			}
		})
	}
}

func Test_RegisterTranslatorFor_EdgeCases(t *testing.T) {
	w := httptest.NewRecorder()
	// Pass an invalid email to trigger the standard English 'email' validator
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"email_address": "not-an-email"}`))

	r := router.NewRouterWithValidation()

	if v, ok := ginbinging.Validator.Engine().(*validator.Validate); ok {
		tr := translator.RegisterTranslatorFor(v)
		r.Use(translator.SetTranslatorMiddleware(tr))
	}

	r.POST("/", func(ctx *gin.Context) {
		var edgeStruct EdgeCaseStruct
		if !binding.MustDecodeBody(ctx, &edgeStruct) {
			// binding will fail here due to validation errors
			return
		}
		ctx.Status(http.StatusOK)
	})

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	body := w.Body.String()

	// 1. Check default english translation and comma splitting
	// (Should be "email_address" not "email_address,omitempty")
	assert.Contains(t, body, "email_address must be a valid email address")

	// 2. Check json:"-" behavior
	// Because name returned is "", the default english translation falls back to the struct field name "Hidden"
	assert.Contains(t, body, "Hidden is a required field")

	// 3. Check fallback when no json tag exists
	assert.Contains(t, body, "NoTag is a required field")
}
