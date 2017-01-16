package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	BATCH_SIZE_DEFAULT int = 10000
	FIELDS_PER_RECORD  int = 2
	CHAN_DEFAULT_SIZE  int = 2000
)

func main() {

	// start timer
	startTime := time.Now()

	// get the file name sent in via the command line flag ------------------------------------------------
	path, showReport := getCommandLineArgs()

	// use the file name to find the file and open it -----------------------------------------------------
	csvFile, err := os.Open(path)
	check(err)

	// defer the closing of the file until the end of main()
	defer func() {
		e := csvFile.Close()
		check(e)
	}()

	// batch read the file line by line (Read 10,000 at a time, and fill up 10k structs) ------------------

	// create a reader to read csv files
	bufReader := bufio.NewReader(csvFile)

	// first record in the csv file will be titles so read it and keep the result for output titles
	tempLine, err := bufReader.ReadString('\n')
	check(err)
	res := strings.Split(tempLine, ",")
	columnNames := []string{strings.TrimSpace(res[0]), strings.TrimSpace(res[1])}

	// create the regex validator group we will use to validate the postcodes
	validator := createMainRegexValidatorGroup()

	// we need a wait group to sync our go routines that are running in parallel
	var readRecordWG sync.WaitGroup

	// run a number of go routines in a parallel pipelines pattern [readFromInputFile_go -> createInputRecords_go -> validateInputRecords_go]
	// each function does its job concurrently until there is no more work to do, the WaitGroup readRecordWG sycncronises them with main()
	readLines_chan := readFromInputFile_go(&readRecordWG, bufReader)
	createdInputRecords_chan := createInputRecords_go(&readRecordWG, readLines_chan)
	validImportRecs, invalidImportRecs := validateInputRecords(createdInputRecords_chan, validator)

	// make the main function wait until all functions in the "readRecordWG" have completed - we need all records to be validated before sorting
	readRecordWG.Wait()

	// run the sorting of each group in its own routine & sync with a wait group
	var sortRecordWG sync.WaitGroup
	sortRecordWG.Add(2)

	// sort the ImportRecords by their rowId
	go func() { sort.Sort(validImportRecs); sortRecordWG.Done() }()
	go func() { sort.Sort(invalidImportRecs); sortRecordWG.Done() }()

	// wait for sorting to complete
	sortRecordWG.Wait()

	// write each collection to a CSV file ----------------------------------------------------------------
	writeOutputFiles(columnNames, validImportRecs, invalidImportRecs)

	if showReport {
		printCompletionReport(startTime, len(validImportRecs), len(invalidImportRecs))
	}
}

// _go to signal that func is concurrent & uses go keyword
func validateInputRecords(in <-chan *ImportRecord, val *RegexValidatorGroup) (validGrp, invalidGrp ImportRecordGroup) {

	// make output groups of import records
	valid := NewImportRecordGroup()
	invalid := NewImportRecordGroup()

	var validateWg sync.WaitGroup
	var appendWg sync.WaitGroup
	validateWg.Add(3)

	// channels made to store the ImportRecords as they are sorted
	validChan := make(chan *ImportRecord, CHAN_DEFAULT_SIZE)
	invalidChan := make(chan *ImportRecord, CHAN_DEFAULT_SIZE)

	// create 3 go routines to validate the records
	for i := 0; i < 3; i++ {
		go func() {
			// keep working as long as the input channel is open
			for rec := range in {
				rec.isValid = val.GroupIsStringValid(rec.postcode)

				// sort the ImportRecords based on their validity
				if rec.isValid {
					validChan <- rec
				} else {
					invalidChan <- rec
				}
			}
			validateWg.Done()
		}()
	}

	// create two go routines that collect the ImportRecords from each of the channels
	appendWg.Add(2)
	go func() {
		for rec := range validChan {
			valid = append(valid, rec)
		}
		appendWg.Done()
	}()

	go func() {
		for rec := range invalidChan {
			invalid = append(invalid, rec)
		}
		appendWg.Done()
	}()

	// we must wait for the validate WaitGroup to complete before we close its channels
	validateWg.Wait()
	close(validChan)
	close(invalidChan)
	appendWg.Wait()

	return valid, invalid
}

func createInputRecords_go(wg *sync.WaitGroup, in <-chan []string) <-chan *ImportRecord {
	// make out output channel & increment the WaitGroup
	wg.Add(1)
	out := make(chan *ImportRecord, CHAN_DEFAULT_SIZE)

	// the creation of new InputRecord structs is done in its own go routine
	go func() {
		// keep working as long as the input channel is open
		for strRec := range in {
			// create an import records from the string slice read from the channel "in" & put it into the outchannel
			out <- NewImportRecord(strRec)
		}
		// close out output channel upon completion & signal completion to the WaitGroup
		close(out)
		wg.Done()
	}()

	// return the channel that we will be putting created ImportRecords into
	return out
}

func readFromInputFile_go(wg *sync.WaitGroup, reader *bufio.Reader) <-chan []string {
	// make our output channel & increment the WaitGroup
	wg.Add(1)
	out := make(chan []string, CHAN_DEFAULT_SIZE)
	completed := false

	// reading of the file runs is done in its own go routine
	go func() {
		for !completed {
			// read a record from each line in the csv file & reading complete when we hit EOF
			line, e := reader.ReadString('\n')
			if e == io.EOF {
				completed = true
				break
			}
			check(e)

			// split the string at the comma and turn it into a string slice, trim any space from each string
			res := strings.Split(line, ",")
			record := []string{strings.TrimSpace(res[0]), strings.TrimSpace(res[1])}

			// place each read record into the channel to send to consumer routine
			out <- record
		}
		// close the channel when we are finished reading & signal completion to the wait group
		close(out)
		wg.Done()
	}()

	// return the channel that we will be putting the lines we have read into
	return out
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

	// write to both output files in parallel & use WaitGroup to sync
	var writerWG sync.WaitGroup

	writerWG.Add(2)
	go func() {
		defer writerWG.Done()
		// create a valid record output file & buffered writer to create said file
		validOutfile, err := os.Create("succeeded_validation.csv")
		check(err)
		validRecWriter := bufio.NewWriter(validOutfile)

		// write the column names first
		_, err = fmt.Fprintf(validRecWriter, "%s,%s\n", columnNames[0], columnNames[1])
		check(err)

		// write each record using our writer
		for _, element := range validRecs {
			temp := fmt.Sprintf("%d,%s", element.rowId, element.postcode)
			_, e := fmt.Fprintln(validRecWriter, temp)
			check(e)
		}

		// flush & close file now we are finished
		validRecWriter.Flush()
		if e := validOutfile.Close(); e != nil {
			check(e)
		}
	}()

	go func() {
		defer writerWG.Done()
		// create a invalid record output file & buffered writer to create said file
		invalidOutfile, err := os.Create("failed_validation.csv")
		check(err)
		invalidRecWriter := bufio.NewWriter(invalidOutfile)

		// write the column names first
		_, err = fmt.Fprintf(invalidRecWriter, "%s,%s\n", columnNames[0], columnNames[1])
		check(err)

		// write each record using our writer
		for _, element := range invalidRecs {
			temp := fmt.Sprintf("%d,%s", element.rowId, element.postcode)
			_, e := fmt.Fprintln(invalidRecWriter, temp)
			check(e)
		}

		// flush & close file now we are finished
		invalidRecWriter.Flush()
		if e := invalidOutfile.Close(); e != nil {
			check(e)
		}
	}()

	writerWG.Wait()
}

func getCommandLineArgs() (string, bool) {
	var path string
	flag.StringVar(&path, "file", "", "the location of the .csv file")

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

	return path, showReport
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
