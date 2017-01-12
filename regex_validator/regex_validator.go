package regex_validator

import (
	"regexp"
)

// type used in conjunction with RegexValidator to signal if a string that matches the supplied
// regex is invalid or valid
type MatchSymantics uint8

const (
	MATCH_MEANS_VALID     = MatchSymantics(iota) // if the string matches the regex it is valid
	MATCH_MEANS_NOT_VALID                        // if the string matches the regex it is NOT valid
)

// type that stores a compiled regex object and a match signal (explained above)
type RegexValidator struct {
	regexObj  *regexp.Regexp
	symantics MatchSymantics
}

// create and return a pointer to a new RegexValidator
func NewRegexValidator(regexObj *regexp.Regexp, symantics MatchSymantics) *RegexValidator {
	return &RegexValidator{regexObj: regexObj, symantics: symantics}
}

// check if a string is valid, to do this it checks id the string matches the regex of the RegexValidator
// and the MatchSymantics of the RegexValidator
func (r *RegexValidator) IsStringValid(str string) bool {
	var isValid bool

	isMatch := r.regexObj.MatchString(str)

	// evaluate if the string is valid or not by checking if it matches the regex & its symantics
	if isMatch == true && r.symantics == MATCH_MEANS_VALID {
		isValid = true

	} else if isMatch == true && r.symantics == MATCH_MEANS_NOT_VALID {
		isValid = false

	} else if isMatch == false && r.symantics == MATCH_MEANS_VALID {
		isValid = false

	} else if isMatch == false && r.symantics == MATCH_MEANS_NOT_VALID {
		// a strange case because not being a match here does not neccessarily mean it is valid overall but
		// in this context it would be considered valid (combination of regexs will ensure only valid strings
		// are accepted, RegexValidatorGroup.IsStringValid applies multiple regexs to a single string)
		isValid = true
	}

	return isValid
}
