package validation_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/glocurrency/commons/validation"
	"github.com/stretchr/testify/assert"
)

type mockFieldLevel struct {
	val string
}

func (m mockFieldLevel) Field() reflect.Value {
	return reflect.ValueOf(m.val)
}

// The other methods are unused by your validator, so can return zero values
func (m mockFieldLevel) Top() reflect.Value      { return reflect.Value{} }
func (m mockFieldLevel) Parent() reflect.Value   { return reflect.Value{} }
func (m mockFieldLevel) FieldName() string       { return "" }
func (m mockFieldLevel) StructFieldName() string { return "" }
func (m mockFieldLevel) Param() string           { return "" }
func (m mockFieldLevel) GetTag() string          { return "" }
func (m mockFieldLevel) ExtractType(field reflect.Value) (reflect.Value, reflect.Kind, bool) {
	return field, field.Kind(), false
}
func (m mockFieldLevel) GetStructFieldOK() (reflect.Value, reflect.Kind, bool) {
	return reflect.Value{}, reflect.Invalid, false
}
func (m mockFieldLevel) GetStructFieldOKAdvanced(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool) {
	return reflect.Value{}, reflect.Invalid, false
}
func (m mockFieldLevel) GetStructFieldOK2() (reflect.Value, reflect.Kind, bool, bool) {
	return reflect.Value{}, reflect.Invalid, false, false
}
func (m mockFieldLevel) GetStructFieldOKAdvanced2(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool, bool) {
	return reflect.Value{}, reflect.Invalid, false, false
}

func TestValidate18YearsOld(t *testing.T) {
	now := time.Now()
	eighteenYearsAgo := now.AddDate(-18, 0, 0)
	seventeenYearsAgo := now.AddDate(-17, 0, 0)
	twentyYearsAgo := now.AddDate(-20, 0, 0)

	tests := []struct {
		name  string
		dob   string
		valid bool
	}{
		{"Exactly 18", eighteenYearsAgo.Format("2006-01-02"), true},
		{"Under 18", seventeenYearsAgo.Format("2006-01-02"), false},
		{"Over 18", twentyYearsAgo.Format("2006-01-02"), true},
		{"Way Over 18", "123456-01-01", false},
		{"Invalid date", "invalid-date", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := mockFieldLevel{val: tt.dob}
			assert.Equal(t, tt.valid, validation.Validate18YearsOld(field))
		})
	}
}

func TestValidate100YearsOld(t *testing.T) {
	now := time.Now()
	hundredYearsAgo := now.AddDate(-100, 0, 0)
	overHundred := now.AddDate(-101, 0, 0)
	underHundred := now.AddDate(-50, 0, 0)

	tests := []struct {
		name  string
		dob   string
		valid bool
	}{
		{"Exactly 100", hundredYearsAgo.Format("2006-01-02"), true},
		{"Over 100", overHundred.Format("2006-01-02"), false},
		{"Under 100", underHundred.Format("2006-01-02"), true},
		{"Invalid date", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := mockFieldLevel{val: tt.dob}
			assert.Equal(t, tt.valid, validation.Validate100YearsOld(field))
		})
	}
}

func TestValidateAlphaNumSpace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"Only letters", "HelloWorld", true},
		{"Letters and numbers", "Hello123", true},
		{"Letters, numbers, spaces", "Hello 123 World", true},
		{"Empty string", "", false},
		{"Symbols", "Hello@123", false},
		{"Dash not allowed", "Hello-123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fl := mockFieldLevel{val: tt.input}
			assert.Equal(t, tt.valid, validation.ValidateAlphaNumSpace(fl))
		})
	}
}

func TestValidateAlphaNumSpaceDash(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"Letters, numbers, spaces", "Hello 123", true},
		{"Includes dash", "Hello-123", true},
		{"Only dash", "-", true},
		{"Symbols not allowed", "Hello@123", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fl := mockFieldLevel{val: tt.input}
			assert.Equal(t, tt.valid, validation.ValidateAlphaNumSpaceDash(fl))
		})
	}
}

func TestValidateBankSupported(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"Letters and numbers", "Bank123", true},
		{"With spaces", "Bank 123", true},
		{"With symbols", `Bank/123-?:().,&"+`, true},
		{"Invalid symbol", "Bank@123", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fl := mockFieldLevel{val: tt.input}
			assert.Equal(t, tt.valid, validation.ValidateBankSupported(fl))
		})
	}
}
