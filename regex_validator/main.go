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
	"sort"
	"time"
)

const (
	BATCH_SIZE_DEFAULT int = 10000
	FIELDS_PER_RECORD  int = 2
)

func main() {

	// start timer
	startTime := time.Now()

	// get the file name sent in via the command line flag ------------------------------------------------
	path, batchSize, showReport := getCommandLineArgs()

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
	validImportRecs := NewImportRecordGroup()
	invalidImportRecs := NewImportRecordGroup()

	// create the regex validator group we will use to validate the postcodes
	validator := createMainRegexValidatorGroup()

	// keep processing in batches
	completed := false
	for !completed {
		for i := 0; i < batchSize; i++ {
			// read a record from each line in the csv file
			currRecord, e := csvReader.Read()

			// reading is completed once we reach the end of the file
			if e == io.EOF {
				completed = true
				break
			}
			check(e)

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

	// sort the ImportRecords by their rowId
	sort.Sort(validImportRecs)
	sort.Sort(invalidImportRecs)

	// write each collection to a CSV file ----------------------------------------------------------------
	writeOutputFiles(columnNames, validImportRecs, invalidImportRecs)

	if showReport {
		printCompletionReport(startTime, len(validImportRecs), len(invalidImportRecs))
	}
}

func printCompletionReport(startTime time.Time, numValid, numInvalid int) {
	// get time since beginning
	elapsed := time.Since(startTime)

	// output short final report
	speed := float64((numValid + numInvalid)) / elapsed.Seconds()
	fmt.Println("-------------------------------------")
	fmt.Println("         Completion Report")
	fmt.Println("-------------------------------------")
	fmt.Printf("Total records: %d\n", numValid+numInvalid)
	fmt.Printf("Succeeded: %d\n", numValid)
	fmt.Printf("Failed: %d\n", numInvalid)
	fmt.Println("-------------------------------------")
	fmt.Printf("Took: %s\n", elapsed)
	fmt.Printf("Speed: %.2f records per second\n", speed)
	fmt.Println("-------------------------------------")
}

func writeOutputFiles(columnNames []string, validRecs, invalidRecs []*ImportRecord) {
	// create and write to an output file
	invalidOutfile, err := os.Create("failed_validation.csv")
	check(err)
	validOutfile, err := os.Create("succeeded_validation.csv")
	check(err)

	// defer closing of the output files until the end of the func
	defer func() {
		if e := invalidOutfile.Close(); e != nil {
			check(e)
		}
		if e := validOutfile.Close(); e != nil {
			check(e)
		}
	}()

	// create a two buffered writers to create the output files
	invalidRecWriter := bufio.NewWriter(invalidOutfile)
	validRecWriter := bufio.NewWriter(validOutfile)

	// write the column names first
	_, err = fmt.Fprintf(invalidRecWriter, "%s,%s\n", columnNames[0], columnNames[1])
	check(err)
	_, err = fmt.Fprintf(validRecWriter, "%s,%s\n", columnNames[0], columnNames[1])
	check(err)

	// write each record usiing our writers
	for _, element := range invalidRecs {
		temp := fmt.Sprintf("%d,%s", element.rowId, element.postcode)
		_, e := fmt.Fprintln(invalidRecWriter, temp)
		check(e)
	}

	for _, element := range validRecs {
		temp := fmt.Sprintf("%d,%s", element.rowId, element.postcode)
		_, e := fmt.Fprintln(validRecWriter, temp)
		check(e)
	}

	// write any buffered data to writer before we finish
	invalidRecWriter.Flush()
	validRecWriter.Flush()
}

func getCommandLineArgs() (string, int, bool) {
	var path string
	flag.StringVar(&path, "file", "", "the location of the .csv file")

	var batchSize int
	flag.IntVar(&batchSize, "batch-size", BATCH_SIZE_DEFAULT, "how many records the program will process at one time")

	var showReport bool
	flag.BoolVar(&showReport, "report", false, "turn on to show a short report upon completion")

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

	return path, batchSize, showReport
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
