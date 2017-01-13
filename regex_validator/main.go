package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello world!")

	// get the file name sent in via the command line flag

	// use the file name to find the file and open it

	// OPTIONAL: untar the file if there exists a library to do so

	// batch read the file line by line (Read 10,000 at a time, and fill up 10k structs)

	// seperate the lines based on thier delimiter

	// put seperated values into a struct that stores the POSTCODE & ROW ID

	// iterate over collection of postcode structs an check validity of each using regex

	// put them into seperate collections depending on validity

	// write each collection to a CSV file

	// do this for another 10k items until completion
}
