package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var alphaNumericSpaceRegex = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)

func ValidateAlphaNumSpace(fl validator.FieldLevel) bool {
	return alphaNumericSpaceRegex.MatchString(fl.Field().String())
}

var alphaNumericSpaceDashRegex = regexp.MustCompile(`^[a-zA-Z0-9\s\-]+$`)

func ValidateAlphaNumSpaceDash(fl validator.FieldLevel) bool {
	return alphaNumericSpaceDashRegex.MatchString(fl.Field().String())
}

var bankSupportedRegex = regexp.MustCompile(`^[a-zA-Z0-9\/\-?:().,&"+\s]+$`)

func ValidateBankSupported(fl validator.FieldLevel) bool {
	return bankSupportedRegex.MatchString(fl.Field().String())
}
