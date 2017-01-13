package main

import "testing"
import "regexp"
import "fmt"

// Perform the setup for the tests ----------------------------------------------------------------------

// the main regex that will be used to validate postcodes (postcodes that match it are valid)
var mainRegex_test = regexp.MustCompile(`(GIR\s0AA)|(((^[A-PR-UWYZ][0-9][0-9]?)|(([A-PR-UWYZ][A-HK-Y][0-9][0-9])|([A-PR-UWYZ][A-HK-Y][0-9])|(WC[0-9][A-Z])|((^[A-PR-UWYZ][0-9][A-HJKPSTUW])|([A-PR-UWYZ][A-HK-Y][0-9][ABEHMNPRVWXY]))))\s[0-9][ABD-HJLNP-UW-Z]{2})`)
var mainRegexValidator_test = NewRegexValidator(mainRegex_test, MATCH_MEANS_VALID)

// regex that will match postcodes that should be excluded from the AA99 prefix (matches postcodes are invalid)
var AA99_exclusionRegex_test = regexp.MustCompile(`((BR|FY|HA|HD|HG|HR|HS|HX|JE|LD|SM|SR|WC|WN|ZE)[0-9][0-9]\s[0-9][ABD-HJLNP-UW-Z]{2})`)
var AA99_exclusionRegexValidator_test = NewRegexValidator(AA99_exclusionRegex_test, MATCH_MEANS_NOT_VALID)

// regex that will match postcode that shoulbe be excluded from the AA9 prefix (matches postcodes are invalid)
var AA9_exclusionRegex_test = regexp.MustCompile(`((AB|LL|SO)[0-9]\s[0-9][ABD-HJLNP-UW-Z]{2})`)
var AA9_exclusionRegexValidator_test = NewRegexValidator(AA9_exclusionRegex_test, MATCH_MEANS_NOT_VALID)

// ---------------------------- UNIT TESTS , regex_validator package ------------------------------------

// -------------- RegexValidator tests
// expected: true
// run the NewRegexValidator construction function and check struct was constructed properly
func Test_NewRegexValidator(t *testing.T) {

	expected := true
	result := false

	re := regexp.MustCompile("abcdefg")
	matSem := MATCH_MEANS_VALID
	reValid := NewRegexValidator(re, matSem)

	result = re == reValid.regexObj && matSem == reValid.symantics

	if result != expected {
		error := fmt.Sprintf("Given regex: %s & MatchSemantics: %s, Expected: %t   got: %t", "abcdefg", "MATCH_MEANS_VALID", expected, result)
		t.Error(error)
	}
}

// expected: true
// call IsStringValid on a validator with "MATCH_MEANS_VALID" symantic and give it a matching string
func Test_IsStringValid__MatchMeansValid_WithMatchingString(t *testing.T) {

	postcode := "EC1A 1BB"
	expected := true
	result := mainRegexValidator_test.IsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given string: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// call IsStringValid on a validator with "MATCH_MEANS_VALID" symantic and give it a string that does not match
func Test_IsStringValid__MatchMeansValid_WithNonMatchingString(t *testing.T) {
	postcode := "LS44PL"
	expected := false
	result := mainRegexValidator_test.IsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given string: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// call IsStringValid on a validator with "MATCH_MEANS_NOT_VALID" symantic and give it a matching string
//
func Test_IsStringValid__MatchMeansNotValid_WithMatchingString(t *testing.T) {
	postcode := "FY10 4PL"
	expected := false
	result := AA99_exclusionRegexValidator_test.IsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given string: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: true
// call IsStringValid on a validator with "MATCH_MEANS_NOT_VALID" symantic and give it a string that does not match
//
func Test_IsStringValid__MatchMeansNotValid_WithNonMatchingString(t *testing.T) {
	postcode := "Q1A 9AA"
	expected := true
	result := AA99_exclusionRegexValidator_test.IsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given string: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: true
// call IsStringValid, the validator will use the "mainRegex_test" and be given a string that matches
func Test_IsStringValid__MainRegex_WithMatchingString(t *testing.T) {
	postcode := "EC1A 1BB"
	expected := true
	result := mainRegexValidator_test.IsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given string: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// call IsStringValid, the validator will use the "mainRegex_test" and be given a string that does not match
//
func Test_IsStringValid__MainRegex_WithNonMatchingString(t *testing.T) {
	postcode := "$%± ()()"
	expected := false
	result := mainRegexValidator_test.IsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given string: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// call IsStringValid, the validator will use the "AA99_exclusionRegex_test" and be given a string that matches
func Test_IsStringValid__AA99ExclusionRegex_WithMatchingString(t *testing.T) {
	postcode := "FY10 4PL"
	expected := false
	result := AA99_exclusionRegexValidator_test.IsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given string: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: true
// call IsStringValid, the validator will use the "AA99_exclusionRegex_test" and be given a string that does not match
func Test_IsStringValid__AA99ExclusionRegex_WithNonMatchingString(t *testing.T) {
	postcode := "EC1A 1BB"
	expected := true
	result := AA99_exclusionRegexValidator_test.IsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given string: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: false
// call IsStringValid, the validator will use the "AA9_exclusionRegex_test" and be given a string that matches
func Test_IsStringValid__AA9ExclusionRegex_WithMatchingString(t *testing.T) {
	postcode := "SO1 4QQ"
	expected := false
	result := AA9_exclusionRegexValidator_test.IsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given string: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// expected: true
// call IsStringValid, the validator will use the "AA9_exclusionRegex_test" and be given a string that does not match
func Test_IsStringValid__AA9ExclusionRegex_WithNonMatchingString(t *testing.T) {
	postcode := "EC1A 1BB"
	expected := true
	result := AA9_exclusionRegexValidator_test.IsStringValid(postcode)

	if result != expected {
		error := fmt.Sprintf("Given string: %s, Expected: %t   got: %t", postcode, expected, result)
		t.Error(error)
	}
}

// -------------- RegexValidatorGroup tests

// expected: true
// run the NewRegexValidatorGroup constructor function and check that struct was constructed properly
func Test_NewRegexValidatorGroup(t *testing.T) {
	expected := true
	result := false

	reValidGroup := NewRegexValidatorGroup()

	result = reValidGroup.validators != nil

	if result != expected {
		error := fmt.Sprintf("Given RegexValidatorGroup to create, Expected: %t, got: %t", expected, result)
		t.Error(error)
	}
}

// expected: true
// create a RegexValidatorGroup, call AddRegexValidator to add 3 RegexValidators to the group check that
// each has been added properly
func Test_AddRegexValidator(t *testing.T) {
	expected := true
	result := false

	reValidator1 := mainRegexValidator_test
	reValidator2 := AA99_exclusionRegexValidator_test
	reValidator3 := AA9_exclusionRegexValidator_test

	reValidGroup := NewRegexValidatorGroup()

	reValidGroup.AddRegexValidator(reValidator1)
	result = len(reValidGroup.validators) == 1 && reValidGroup.validators[0] == reValidator1

	if result != expected {
		error := fmt.Sprintf("Given %s, Expected: %t, got: %t", "reValidator1", expected, result)
		t.Error(error)
	}

	reValidGroup.AddRegexValidator(reValidator2)
	result = len(reValidGroup.validators) == 2 && reValidGroup.validators[0] == reValidator1 && reValidGroup.validators[1] == reValidator2

	if result != expected {
		error := fmt.Sprintf("Given %s, Expected: %t, got: %t", "reValidator1 & 2", expected, result)
		t.Error(error)
	}

	reValidGroup.AddRegexValidator(reValidator3)
	result = len(reValidGroup.validators) == 3 && reValidGroup.validators[0] == reValidator1 && reValidGroup.validators[1] == reValidator2 && reValidGroup.validators[2] == reValidator3

	if result != expected {
		error := fmt.Sprintf("Given %s, Expected: %t, got: %t", "reValidator1 , 2 & 3", expected, result)
		t.Error(error)
	}
}

// expected: true
// call GroupIsStringValid on a series of valid postcodes, the RegexValidatorGroup contains the 3 RegexValidators
// defined at the top of the testing area
func Test_GroupIsStringValid_GivenValidStrings(t *testing.T) {
	expected := true
	result := false

	reValidator1 := mainRegexValidator_test
	reValidator2 := AA99_exclusionRegexValidator_test
	reValidator3 := AA9_exclusionRegexValidator_test

	reValidGroup := NewRegexValidatorGroup()
	reValidGroup.AddRegexValidator(reValidator1)
	reValidGroup.AddRegexValidator(reValidator2)
	reValidGroup.AddRegexValidator(reValidator3)

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
		result = reValidGroup.GroupIsStringValid(element)

		if result != expected {
			error := fmt.Sprintf("Given postcode: %s, Expected: %t, got: %t", element, expected, result)
			t.Error(error)
		}
	}
}

// expected: false
// call GroupIsStringValid on a series of invalid postcodes, the RegexValidatorGroup contains the 3
// RegexValidators defined at the top of the testing area - in particular this test should return
// postcodes "FY10 4PL" & "SO1 4QQ" as being invalid and the mainRegexValidator_test cannot do this by itself
func Test_GroupIsStringValid_GivenInvalidStrings(t *testing.T) {
	expected := false
	result := false

	reValidator1 := mainRegexValidator_test
	reValidator2 := AA99_exclusionRegexValidator_test
	reValidator3 := AA9_exclusionRegexValidator_test

	reValidGroup := NewRegexValidatorGroup()
	reValidGroup.AddRegexValidator(reValidator1)
	reValidGroup.AddRegexValidator(reValidator2)
	reValidGroup.AddRegexValidator(reValidator3)

	postcodes := []string{"$%± ()()",
		"A1 9A",
		"FY10 4PL",
		"SO1 4QQ"}

	for _, element := range postcodes {
		result = reValidGroup.GroupIsStringValid(element)

		if result != expected {
			error := fmt.Sprintf("Given postcode: %s, Expected: %t, got: %t", element, expected, result)
			t.Error(error)
		}
	}
}
