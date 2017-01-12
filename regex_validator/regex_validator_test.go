package regex_validator

import "testing"
import "regexp"

func TestIsStringValid(t *testing.T) {
	var re = regexp.MustCompile(`(GIR\s0AA)|(((^[A-PR-UWYZ][0-9][0-9]?)|(([A-PR-UWYZ][A-HK-Y][0-9][0-9])|([A-PR-UWYZ][A-HK-Y][0-9])|(WC[0-9][A-Z])|((^[A-PR-UWYZ][0-9][A-HJKPSTUW])|([A-PR-UWYZ][A-HK-Y][0-9][ABEHMNPRVWXY]))))\s[0-9][ABD-HJLNP-UW-Z]{2})`)

	var matchValidator = NewRegexValidator(re, MATCH_MEANS_VALID)

	expected := true
	result := matchValidator.IsStringValid("LZ10 3QP")

	if result != expected {
		t.Error("Expected true got", result)
	}
}
