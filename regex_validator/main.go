package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func errorExit(str string, code int) {
	fmt.Fprintf(os.Stderr, "%s\n", str)
	os.Exit(code)
}

func main() {

	// get the file name sent in via the command line flag ------------------------------------------------
	var path string
	flag.StringVar(&path, "file", "", "the location of the .csv file")
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

	// OPTIONAL: untar the file if there exists a library to do so ----------------------------------------

	// batch read the file line by line (Read 10,000 at a time, and fill up 10k structs) ------------------

	// seperate the lines based on thier delimiter --------------------------------------------------------

	// put seperated values into a struct that stores the POSTCODE & ROW ID -------------------------------

	// iterate over collection of postcode structs an check validity of each using regex ------------------

	// put them into seperate collections depending on validity -------------------------------------------

	// write each collection to a CSV file ----------------------------------------------------------------

	// do this for another 10k items until completion -----------------------------------------------------

	fmt.Println("reached end successfully")
}
