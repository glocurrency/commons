package translator

import (
	"reflect"
	"strings"

	locale "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
)

func RegisterTranslatorFor(v *validator.Validate) ut.Translator {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	uni := ut.New(locale.New())
	fallback := uni.GetFallback()

	en.RegisterDefaultTranslations(v, fallback)

	v.RegisterTranslation("bic", fallback, func(ut ut.Translator) error {
		return ut.Add("bic", "{0} must comply with BIC format", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("bic", fe.Field())
		return t
	})

	v.RegisterTranslation("iso3166_1_alpha2", fallback, func(ut ut.Translator) error {
		return ut.Add("iso3166_1_alpha2", "{0} must be a valid country code", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("iso3166_1_alpha2", fe.Field())
		return t
	})

	v.RegisterTranslation("alphanumspace", fallback, func(ut ut.Translator) error {
		return ut.Add("alphanumspace", "{0} can only contain alphanumeric characters and spaces", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphanumspace", fe.Field())
		return t
	})

	v.RegisterTranslation("alphanumspacedash", fallback, func(ut ut.Translator) error {
		return ut.Add("alphanumspacedash", "{0} can only contain alphanumeric characters, spaces and dashes", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphanumspacedash", fe.Field())
		return t
	})

	v.RegisterTranslation("banksupported", fallback, func(ut ut.Translator) error {
		return ut.Add("banksupported", "{0} can only contain bank suported characters", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("banksupported", fe.Field())
		return t
	})

	v.RegisterTranslation("18yo", fallback, func(ut ut.Translator) error {
		return ut.Add("18yo", "{0} should be at least 18 years old", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("18yo", fe.Field())
		return t
	})

	v.RegisterTranslation("notold", fallback, func(ut ut.Translator) error {
		return ut.Add("notold", "{0} should be a valid age", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("notold", fe.Field())
		return t
	})

	return fallback
}
