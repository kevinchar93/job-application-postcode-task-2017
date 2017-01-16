package main

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_IsStringValid__NewImportRecord(t *testing.T) {

	expected := true
	result := false

	testRecords := [][]string{[]string{"1064397", "MK12 5EY"},
		[]string{"1995262", "W6 8EX"},
		[]string{"803671", "IP20 9DL"},
		[]string{"1779862", "SW12 0EF"},
		[]string{"1434686", "RG10 8AU"},
		[]string{"1508688", "RM20 4AP"},
		[]string{"472843", "DH4 7DU"}}

	for _, element := range testRecords {
		item := NewImportRecord(element)

		num, _ := strconv.Atoi(element[0])
		result = item.isValid == false && item.postcode == element[1] && item.rowId == uint32(num)

		if result != expected {
			error := fmt.Sprintf("Given NewImportRecord, Expected: %t   got: %t", expected, result)
			t.Error(error)
		}
	}

}
