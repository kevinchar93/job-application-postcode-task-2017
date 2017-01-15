package main

import (
	"strconv"
)

// type to represet a record from an imported .csv file , rowId & postcodes or the record is stored in
// their native types rather than both being stored as strings
type ImportRecord struct {
	rowId         uint32
	postcode      string
	beenValidated bool
	isValid       bool
}

// takes a record read by a cvs reader (which creates a string slice of each record) and creates a properly
// typed ImportRecord from this slice
func NewImportRecord(record []string) *ImportRecord {
	if len(record) > FIELDS_PER_RECORD || len(record) < FIELDS_PER_RECORD {
		panic("invalid record received")
	}

	const ROW_ID_IDX = 0
	const POSTCODE_IDX = 1

	rowIdInt, err := strconv.ParseInt(record[ROW_ID_IDX], 10, 32)
	check(err)

	return &ImportRecord{rowId: uint32(rowIdInt),
		postcode:      record[POSTCODE_IDX],
		beenValidated: false,
		isValid:       false}
}
