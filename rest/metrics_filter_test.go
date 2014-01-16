package rest

import (
	"strings"
	"testing"
)

func TestKeySanitization(t *testing.T) {

	sanitizedKey := sanitizeKey("()\\/?<>;:~`@#$%^&*+=")
	if strings.ContainsAny(sanitizedKey, " ,()\\/?<>;:~@#$%^&*+=") {
		t.Error("Invalid sanitization. Returned string :", sanitizedKey)
	}
}
