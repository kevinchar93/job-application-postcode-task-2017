package main

import "testing"
import "regexp"
import "fmt"
import "os"

// Perform the setup for the tests ----------------------------------------------------------------------

// the main regex that will be used to validate postcodes (postcodes that match it are valid)
var mainRegex = regexp.MustCompile(`(GIR\s0AA)|(((^[A-PR-UWYZ][0-9][0-9]?)|(([A-PR-UWYZ][A-HK-Y][0-9][0-9])|([A-PR-UWYZ][A-HK-Y][0-9])|(WC[0-9][A-Z])|((^[A-PR-UWYZ][0-9][A-HJKPSTUW])|([A-PR-UWYZ][A-HK-Y][0-9][ABEHMNPRVWXY]))))\s[0-9][ABD-HJLNP-UW-Z]{2})`)
var mainRegexValidator = NewRegexValidator(mainRegex, MATCH_MEANS_VALID)

// regex that will match a subset of postcodes with the AA99 prefix (postcodes that match it are invalid)
var AA99_exclusionRegex = regexp.MustCompile(`((BR|FY|HA|HD|HG|HR|HS|HX|JE|LD|SM|SR|WC|WN|ZE)[0-9][0-9]\s[0-9][ABD-HJLNP-UW-Z]{2})`)
var AA99_exclusionRegexValidator = NewRegexValidator(AA99_exclusionRegex, MATCH_MEANS_NOT_VALID)

// regex that will match a subset of postcodes with the AA9 prefix (postcodes that match it are invalid)
var AA9_exclusionRegex = regexp.MustCompile(`((AB|LL|SO)[0-9]\s[0-9][ABD-HJLNP-UW-Z]{2})`)
var AA9_exclusionRegexValidator = NewRegexValidator(AA9_exclusionRegex, MATCH_MEANS_NOT_VALID)
var postCodeRegexValidator = NewRegexValidatorGroup()

func TestMain(m *testing.M) {
	testSetup()
	retCode := m.Run()
	os.Exit(retCode)
}

func testSetup() {
	postCodeRegexValidator.AddRegexValidator(mainRegexValidator)
	postCodeRegexValidator.AddRegexValidator(AA99_exclusionRegexValidator)
	postCodeRegexValidator.AddRegexValidator(AA9_exclusionRegexValidator)
}

// Unit tests for Task 1 are below ----------------------------------------------------------------------

// expected: false
// given invalid postcode, Junk
func Test_GroupIsStringValid__Junk(t *testing.T) {
	postcode := "$%± ()()"
	expected := false
	result := postCodeRegexValidator.GroupIsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given postcode: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// given invalid postcode, Invalid
func Test_GroupIsStringValid__Invalid(t *testing.T) {
	postcode := "XX XXX"
	expected := false
	result := postCodeRegexValidator.GroupIsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given postcode: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// given invalid postcode, 	Incorrect inward code length
func Test_GroupIsStringValid__IncorrectInwardCodeLength(t *testing.T) {
	postcode := "A1 9A"
	expected := false
	result := postCodeRegexValidator.GroupIsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given postcode: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// given invalid postcode, No space
func Test_GroupIsStringValid__NoSpace(t *testing.T) {
	postcode := "LS44PL"
	expected := false
	result := postCodeRegexValidator.GroupIsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given postcode: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// given invalid postcode, 'Q' in first position
func Test_GroupIsStringValid__QinFirstPosition(t *testing.T) {
	postcode := "Q1A 9AA"
	expected := false
	result := postCodeRegexValidator.GroupIsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given postcode: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// given invalid postcode, 	'V' in first position
func Test_GroupIsStringValid__VinFirstPosition(t *testing.T) {
	postcode := "V1A 9AA"
	expected := false
	result := postCodeRegexValidator.GroupIsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given postcode: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// given invalid postcode, 'X' in first position
func Test_GroupIsStringValid__XinFirstPosition(t *testing.T) {
	postcode := "X1A 9BB"
	expected := false
	result := postCodeRegexValidator.GroupIsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given postcode: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// given invalid postcode, 'I' in second position
func Test_GroupIsStringValid__IinSecondPosition(t *testing.T) {
	postcode := "LI10 3QP"
	expected := false
	result := postCodeRegexValidator.GroupIsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given postcode: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// given invalid postcode, 'J' in second position
func Test_GroupIsStringValid__JinSecondPosition(t *testing.T) {
	postcode := "LJ10 3QP"
	expected := false
	result := postCodeRegexValidator.GroupIsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given postcode: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// given invalid postcode, 'Z' in second position
func Test_GroupIsStringValid__ZinSecondPosition(t *testing.T) {
	postcode := "LZ10 3QP"
	expected := false
	result := postCodeRegexValidator.GroupIsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given postcode: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// given invalid postcode, 'Q' in third position with 'A9A' structure
func Test_GroupIsStringValid__QinThirdPositionWithA9Astructure(t *testing.T) {
	postcode := "A9Q 9AA"
	expected := false
	result := postCodeRegexValidator.GroupIsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given postcode: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// given invalid postcode, 'C' in fourth position with 'AA9A' structure
func Test_GroupIsStringValid__CinForthPositionWithA9Astructure(t *testing.T) {
	postcode := "AA9C 9AA"
	expected := false
	result := postCodeRegexValidator.GroupIsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given postcode: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// given invalid postcode, Area with only single digit districts
func Test_GroupIsStringValid__AreaWithOnlySingleDigitDistricts(t *testing.T) {
	postcode := "FY10 4PL"
	expected := false
	result := postCodeRegexValidator.GroupIsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given postcode: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// given invalid postcode, Area with only double digit districts
func Test_GroupIsStringValid__AreaWithOnlyDoubleDigitDistricts(t *testing.T) {
	postcode := "SO1 4QQ"
	expected := false
	result := postCodeRegexValidator.GroupIsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given postcode: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: true
// check a collection of valid postcodes there should be no problem
func Test_GroupIsStringValid__NoExpectedProblem(t *testing.T) {
	expected := true
	result := false

	postcodes := []string{"EC1A 1BB",
		"W1A 0AX",
		"M1 1AE",
		"B33 8TH",
		"CR2 6XH",
		"DN55 1PT",
		"GIR 0AA",
		"SO10 9AA",
		"FY9 9AA",
		"WC1A 9AA"}

	for _, element := range postcodes {
		result = postCodeRegexValidator.GroupIsStringValid(element)

		if result != expected {
			error := fmt.Sprintf("Given postcode: %s, Expected: %t, got: %t", element, expected, result)
			t.Error(error)
		}
	}
}
