package validation

import (
	"regexp"
	"time"

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

func Validate18YearsOld(fl validator.FieldLevel) bool {
	dob, err := time.Parse(time.DateOnly, fl.Field().String())
	if err != nil {
		return false
	}

	now := time.Now()
	age := now.Year() - dob.Year()

	// If the current month is before the birth month, OR
	// it's the birth month but the current day is before the birth day...
	// ...they haven't had their birthday yet this year.
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}

	return age >= 18
}

func Validate100YearsOld(fl validator.FieldLevel) bool {
	dob, err := time.Parse(time.DateOnly, fl.Field().String())
	if err != nil {
		return false
	}

	now := time.Now()
	age := now.Year() - dob.Year()

	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}

	return age <= 100
}
