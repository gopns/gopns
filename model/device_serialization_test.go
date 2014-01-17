package model

import (
	"testing"
)

func TestLocaleValidation(t *testing.T) {

	if ValidateLocale("en_US") != nil {
		t.Error("en_US is a valid locale but code is saying it isn't")
	}

	if ValidateLocale("en_US1") == nil {
		t.Error("en_US1 is not a valid locale but code is saying it is")
	}

	if ValidateLocale("enUS") == nil {
		t.Error("enUS is not a valid locale but code is saying it is")
	}

	if ValidateLocale("en_USSS") == nil {
		t.Error("en_USSS is not a valid locale but code is saying it is")
	}

}
