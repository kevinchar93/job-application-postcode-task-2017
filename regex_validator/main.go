package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
func NewImportItem(record []string) *ImportRecord {
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

const (
	BATCH_SIZE_DEFAULT int = 10000
	FIELDS_PER_RECORD  int = 2
)

func main() {

	// get the file name sent in via the command line flag ------------------------------------------------
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

	// use the file name to find the file and open it -----------------------------------------------------
	csvFile, err := os.Open(path)
	check(err)

	// defer the closing of the file until the end of main()
	defer func() {
		err := csvFile.Close()
		check(err)
	}()

	// batch read the file line by line (Read 10,000 at a time, and fill up 10k structs) ------------------

	// create a reader that can read csv files , init the readers values
	csvReader := csv.NewReader(csvFile)
	csvReader.FieldsPerRecord = FIELDS_PER_RECORD

	// make slices to hold the valid and invalid ImportRecords
	//validImportRecs := make([]*ImportRecord, 0)
	//invalidImportRecs := make([]*ImportRecord, 0)

	// keep processing in batches
	completed := false
	for !completed {

		for i := 0; i < batchSize; i++ {
			currRecord, err := csvReader.Read()

			// reading is completed once we reach the end of the file
			if err == io.EOF {
				completed = true
				break
			}

			fmt.Println(currRecord)

		}
		//completed = true
	}

	// put seperated values into a struct that stores the POSTCODE & ROW ID -------------------------------

	// iterate over collection of postcode structs an check validity of each using regex ------------------

	// put them into seperate collections depending on validity -------------------------------------------

	// write each collection to a CSV file ----------------------------------------------------------------

	// do this for another 10k items until completion -----------------------------------------------------

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
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
