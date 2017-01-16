package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

const (
	BATCH_SIZE_DEFAULT int = 10000
	FIELDS_PER_RECORD  int = 2
)

func main() {

	// get the file name sent in via the command line flag ------------------------------------------------
	path, batchSize := getCommandLineArgs()

	// use the file name to find the file and open it -----------------------------------------------------
	csvFile, err := os.Open(path)
	check(err)

	// defer the closing of the file until the end of main()
	defer func() {
		e := csvFile.Close()
		check(e)
	}()

	// batch read the file line by line (Read 10,000 at a time, and fill up 10k structs) ------------------

	// create a reader that can read csv files , init the readers values
	csvReader := csv.NewReader(csvFile)
	csvReader.FieldsPerRecord = FIELDS_PER_RECORD

	// first record in the csv file will be titles so read it and discard the record
	columnNames, err := csvReader.Read()
	check(err)

	// make slices to hold the valid and invalid ImportRecords
	validImportRecs := make([]*ImportRecord, 0)
	invalidImportRecs := make([]*ImportRecord, 0)

	// create the regex validator group we will use to validate the postcodes
	validator := createMainRegexValidatorGroup()

	// keep processing in batches
	completed := false
	for !completed {
		for i := 0; i < batchSize; i++ {
			// read a record from each line in the csv file
			currRecord, e := csvReader.Read()
			if e == io.EOF {
				completed = true
				break
			}
			check(e)

			// reading is completed once we reach the end of the file

			// create an ImportRecord from each csv record & check to see if the record's postcode is valid
			importRec := NewImportRecord(currRecord)
			importRec.isValid = validator.GroupIsStringValid(importRec.postcode)
			importRec.beenValidated = true

			// sort the ImportRecords based on their validity
			if importRec.isValid {
				validImportRecs = append(validImportRecs, importRec)
			} else {
				invalidImportRecs = append(invalidImportRecs, importRec)
			}
		}
	}

	// write each collection to a CSV file ----------------------------------------------------------------
	// create and write to an output file
	invalidImportsOutFile, err := os.Create("failed_validation.csv")
	check(err)

	// defer closing of the output file until the end of main
	defer func() {
		if e := invalidImportsOutFile.Close(); err != nil {
			check(e)
		}
	}()

	// create a new buffered writer to create the output file
	invalidRecWriter := bufio.NewWriter(invalidImportsOutFile)

	// write the column names first
	_, err = fmt.Fprintf(invalidRecWriter, "%s,%s\n", columnNames[0], columnNames[1])
	check(err)

	// write each record usiing our writer
	for _, elemement := range invalidImportRecs {
		temp := fmt.Sprintf("%d,%s", elemement.rowId, elemement.postcode)
		_, e := fmt.Fprintln(invalidRecWriter, temp)
		check(e)
	}

	// write any buffered data to writer before we finish
	invalidRecWriter.Flush()

}

func getCommandLineArgs() (string, int) {
	var path string
	flag.StringVar(&path, "file", "", "the location of the .csv file")

	var batchSize int
	flag.IntVar(&batchSize, "batch-size", BATCH_SIZE_DEFAULT, "how many records the program will process at one time")

	flag.Parse()

	if len(path) == 0 {
		errorExit("No path to or name of a .csv file was provided", 1)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		errStr := fmt.Sprintf("The file or path provided does not exist: \"%s\"", path)
		errorExit(errStr, 1)
	}

	if filepath.Ext(path) != ".csv" {
		errorExit("File must have the extension .csv", 1)
	}

	return path, batchSize
}

func createMainRegexValidatorGroup() *RegexValidatorGroup {
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

	postCodeRegexValidator.AddRegexValidator(mainRegexValidator)
	postCodeRegexValidator.AddRegexValidator(AA99_exclusionRegexValidator)
	postCodeRegexValidator.AddRegexValidator(AA9_exclusionRegexValidator)

	return postCodeRegexValidator
}

func errorExit(str string, code int) {
	fmt.Fprintf(os.Stderr, "%s\n", str)
	os.Exit(code)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
