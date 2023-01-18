
## NOTE: this program was written with Go 1.7 in early 2017 - it may note work with the latest version of Go: https://go.dev/dl/

# postcode-task, Kevin Charles

My solution for the task is written in Go 1.7 and can be built and run on any operating system. You will need an installation of the Go runtime to run the program, simple installation instructions can be found at https://golang.org/doc/install.

Please read the rest of the readme for usage instructions and other notes

- [postcode-task, Kevin Charles](#postcode-task-kevin-charles)
  - [Folder Structure](#folder-structure)
  - [Usage of Git \& tags](#usage-of-git--tags)
  - [Building the program](#building-the-program)
  - [Operating the program](#operating-the-program)
  - [Choice of The Go Programming Language](#choice-of-the-go-programming-language)
  - [About Task 1](#about-task-1)
    - [Notes](#notes)
  - [About Task 2](#about-task-2)
    - [Notes](#notes-1)
  - [About Task 3](#about-task-3)
    - [Profiling tools](#profiling-tools)
    - [Notes](#notes-2)

----------

## Folder Structure
```
job-application-postcode-task-2017/
├── regex_validator
└── profiling
    ├── 1__TASK_2_REL
    ├── 2__TASK_3_SEQ
    ├── 3__TASK_CSV_READER_REPLACED
    ├── 4__TASK_3_PIPLINE_IMPLEMENTED
    └── 5__TASK_3_PARA
```

| Folder | Description |
|------------------------------------|-----------------------------------------------------------------------------|
| `job-application-postcode-task-2017` | root folder |
| `../regex_validator` | by default holds final version of the program |
| `../profiling` | holds folders with images of profiler results |
| `../../1_TASK_2_REL` | holds profiler results of program at TASK_2_REL tag in Git |
| `../../2_TASK_3_SEQ` | holds profiler results of program at TASK_3_SEQ tag in Git |
| `../../3_TASK_CSV_READER_REPLACED` | holds profiler results of program at TASK_CSV_READER_REPLACED tag in Git |
| `../../4_TASK_3_PIPLINE_IMPLEMENTED` | holds profiler results of program at TASK_3_PIPLINE_IMPLEMENTED tag in Git |
| `../../5_TASK_3_PARA` | holds profiler results of program at TASK_3_PARA (final version) tag in Git |


---

## Usage of Git & tags

In making this program I have used the Git version control system and made use of tags to make switching between different versions of the program simple. The key tag names are as follows:

| Tag Name    | Program State                                                                                                        |
|-------------|----------------------------------------------------------------------------------------------------------------------|
| TASK_3_PARA | *Final version of the program that uses concurrency for optimisation.  It meets the requirements of tasks 1, 2, & 3.* |
| TASK_3_SEQ  | *Version of the program that meets the requirements of tasks 1, 2 & 3 but it does so without any optimisations.*       |
| TASK_2_REL  | *Version of the program that meets the requirements of tasks 1 & 2.*                                                   |

You can use Git (in the commandline) to switch between these versions of the program when in the folder **submission/code/FinalCode**, the command is as follows:

*git checkout tags/**TAG_NAME***

example : ```git checkout tags/TASK_2_REL```

This will change the source code in the current folder to match the sate of the code at that tag in the history.

----

## Building the program

Once you have the Go runtime installed and have either selected a tag (or folder) with the version you wish to build, you can then build the program. In Windows a cmd window must be used, in unix / unix-like the terminal must be used.

Make sure you are inside the `regex_validator` folder and run the command `go build`. This will produce an executable called `regex_validator` (`regex_validator.exe` on Windows) which can them be run using the instructions in the next section.

`go clean` will clean the folder of the executable

---

## Operating the program

The different versions of the program operate largely in the same way with some minor output differences. In Windows a cmd window must be used, in unix / unix-like the terminal must be used.

To run each version of the program call the executable name with the flag `-file` and a path to a .csv file (*.gz file must be extracted*).  The programs Task 3 versions of the program will have different output than the Task 2 version , details below.

Example on macOS :

    ./regex_validator -file ~/Desktop/nhs_digital_job/import_data.csv

*note on Windows the .exe extension would be used and the full path must be specified*

**TASK_3_PARA**
The program will produce two files `failed_validation.csv` and `succeeded_validation.csv` in the same folder as the executable, that store in ascending order the invalid and valid records respectively.

**TASK_3_SEQ**
The program will produce two files `failed_validation.csv` and `succeeded_validation.csv`in the same folder as the executable, that store in ascending order the invalid and valid records respectively.

**TASK_2_REL**
The program will produce a single file `failed_validation.csv` in the same folder as the executable, the records will be stored in the order in which they were found to be invalid (not sorted)

Run the program with the `-h` flag to get the full list of flags each version supports

    ./regex_validator -h

*note on Windows the .exe extension would be used*

---

## Choice of The Go Programming Language

Go is a free open source programming language created by Google, I decided to use it for a number of reasons

- **Compiled:**
Go is a compiled language and it can be said that in general a good compiler is faster than an interpreter. Compiled Go code is very fast and has execution times similar to C programs, for a program of this size it might be possible to complete the task before some interpreters have loaded.

- **Garbage collection**
Go is a garbage collected language so I do not have the mental overhead of manual memory management or have to implement a memory management pattern, GC'd languages tend to be interpreted or use a VM  which add overhead. Go avoids this giving you the ease of GC with speed
- **Statically typed**
With Go being statically typed the compiler can point out many errors before any code is run, also Go's syntax has a number of constructs that can make it feel like programming in a dynamically typed language.

- **Concurrency built in**
Go has built in support for a number of concurrency constructs most notably the Goroutine which can communicate with each other using another construct called a Channel. Goroutines are extremely light weight compared to traditional threads - a program could spawn 1 million of them with ease.

- **Large standard library**
Having the backing of Google, Go's standard library is very large given its age, this makes writing cross platform idiomatic Go easier

- **Many language tools**
Go ships with a number of language tools that make the development process easier, there are tools for testing, benching and profiling programs which I have used and many more.

---

## About Task 1

The unit tests for Task 1 are defined in `/regex_validator/task1_unit_test.go`

The tests for the project can be run with (inside the `regex_validator` folder):
`go test` or `go test -v` for verbose output

### Notes
I begun this task by breaking down the regex into sections based on the capture groups and tried it out using some online tools. In writing and running the unit tests I discovered an issue.

The sections of the regex that were meant to capture the shorter A9/A99 and A9A postcode prefix, `([A-PR-UWYZ][0-9][0-9]?)` and `([A-PR-UWYZ][0-9][A-HJKPSTUW])` respectively, were finding postcodes that match them inside of invalid AA99 or AA9A prefixes.

For example `SO1 4QQ` is an invalid AA9 prefix (`SO1`) but a valid A9 prefix (`O1`) can be found inside it. To solve this I modified the regex slightly adding the caret (`^`)  to the beginning of each regex to make them match from the beginning of the string.

**Regex support issue in Go**

In looking up Go's standard regex package I noticed it lacked support for negative look behind which was used in the regex provided - I still wanted to use Go due to the reasons stated earlier so I went about creating a work around.

I implemented the types ***RegexValidator*** and ***RegexValidatorGroup*** (file `regex_validator.go`), a ***RegexValidator*** is a struct that stores a **regex object** and **match semantics**.

The regex object is an object from Go's standard library that can be used to see if a given string matches its regex. The match semantics is a flag that determines what a string matching the regex object means - a match could mean the string is valid or invalid, this is set upon construction. The `IsStringValid` function of a ***RegexValidator*** uses both of these to evaluate and return whether or not a given string is valid.

The type ***RegexValidatorGroup*** is simply a collection of ***RegexValidators***, its function `GroupIsStringValid` will evaluate a given string an return valid if all the ***RegexValidators*** in the group evaluate it as being valid.

Using this simple but powerful type I could use the original regex with the negative look behinds removed and use regex's that match what is defined in said negative look behinds (using match-means-invalid semantics) in a ***RegexValidatorGroup*** to come to the same outcome as the original regex.

**Modified regexs used for validation**

Main regex *- postcodes that match it are valid, note added carets(^)*


       (GIR\s0AA)|(((^[A-PR-UWYZ][0-9][0-9]?)|(([A-PR-UWYZ][A-HK-Y][0-9][0-9])|([A-PR-UWYZ][A-HK-Y][0-9])|(WC[0-9][A-Z])|((^[A-PR-UWYZ][0-9][A-HJKPSTUW])|([A-PR-UWYZ][A-HK-Y][0-9][ABEHMNPRVWXY]))))\s[0-9][ABD-HJLNP-UW-Z]{2})

AA99 exclusion regex - *postcodes that match it are invalid*

    ((BR|FY|HA|HD|HG|HR|HS|HX|JE|LD|SM|SR|WC|WN|ZE)[0-9][0-9]\s[0-9][ABD-HJLNP-UW-Z]{2})

AA9 exclusion regex - *postcodes that match it are invalid*

    ((AB|LL|SO)[0-9]\s[0-9][ABD-HJLNP-UW-Z]{2})

---

## About Task 2

Running the program for task two is covered in the section **"Operating the program"**

### Notes

The overall algorithm that achieves Task 2 is very simple

 1. Get location of .csv file name via command line & use to open file
 2. Open a reader to begin reading the file,  read the file line by line
 3. Turn each line into a string array with 2 values , 1 for row id another for postcode
 4. Create ***ImportItem*** struct from each string array
 5. Validate each ImportItem using ***RegexValidatorGroup***
 6. Sort into 2 arrays based on Validity (valid/invalid)
 7. Write each array to its own output file using writer

The ***ImportItem*** is a struct that stores the **row id**, **postcode** and **validity** of a record read from the .csv file, upon construction it converts the string representation of its row id into an integer so they can be sorted later using this value.

Profiling results of the program when it was at tag `TASK_2_REL` can be found in `submission/profiling/1__TASK_2_REL` , a cpu profile & time stats from the unix `time` command were created.

---

## About Task 3

Running the program for task three is covered in the section **“Operating the program”**

### Profiling tools
**go pprof**
I made use of this package https://github.com/pkg/profile and the pprof tool in Go to get detailed CPU profiling results, the linked package outputs a .pprof file on completion. The Go pprof tool takes the .pprof file and can create a visual representation of the profiling results.

These can be seen in the folder  `submission/profiling/TASK_TAG_HERE/` with the name `cpu.png`

I used these results to determine what parts of the program were using the most CPU time.

**time command unix**
time is a command on unix and unix-like (macOS in this case) systems that can be used to determine the duration of execution of another command, the command in this situation being my program. I used it mainly to measure wall time  taken to complete my program.

The wiki page provides more information https://en.wikipedia.org/wiki/Time_(Unix)

### Notes

To begin with I went about extending the program from tag `TASK_2_REL` to achieve the requirements of task 3 without looking into making any optimisations, I simply needed to add sorting of the ***ImportRecords*** and the writing of `succeeded_validation.csv`

After doing this I profiled & timed the program(`submission/profiling/2__TASK_3_SEQ/`).  In doing this I noticed the function `csv.Reader.Read` was using the most CPU time so I went optimising this by using a slightly faster reader, a buffered reader. This improved average wall execution time from 19.72s to 17.19s

I continued doing more profiling (`submission/profiling/3_TASK_CSV_READER_REPLACED/`). Despite the read time of the buffered reader being faster than the csv reader it was still taking up the majority of the CPU time.

I devised a simple concurrent solution based on the pipeline pattern. I used Go's Goroutines and typed Channels to create a series of functions that could run concurrently with one another:

**Pipeline solution**

 1. Function `readFromInputFile_go` begins reading the csv file's lines putting them into the buffered channel - `readLines_chan`, it lets the main Goroutine continue as it does this

 2. Function `createInputRecords_go` takes as input, channel `readLines_chan`. It reads lines placed into it & creates ***InputRecords*** from them. The new records are placed into the buffered channel - `createdInputRecords_chan`, it lets the main Goroutine continue as it does this

 3. Function `validateInputRecords` takes as input, channel `createdInputRecords_chan`. It validates the ***InputRecords*** from it and returns 2 arrays (valid & invalid ***InputRecords***)

 4.  Function `validateInputRecords` returns when the two earlier functions (`readFromInputFile_go` & `createInputRecords_go`)have completed their work and close their channels

 5. Sorting of each array of validated input records is done concurrently in their own Goroutines

 6. Writing of each output file is done concurrently in their own Goroutines

This solution improved average wall execution time from 17.9s to 15.1s, a significant improvement. This was mostly due to the program being allowed to continue without having to wait for the slow io, the CPU time spent on the buffered readers `ReadString` function fell from 33.92% to 2.05%.

After this I did some more profiling(`submission/profiling/4__TASK_3_PIPLINE_IMPLEMENTED/`) and saw that the function `validateInputRecords` was now using the majority of the CPU time so I decided to see what Improvements could be made there.

The changes I made were fairly simple, I still used the Pipeline solution listed above but I modified the `validateInputRecords` to spawn 3 Goroutines that validate ***InputRecords*** from the `createdInputRecords_chan` and put their results into two local channels `validChan` & `invalidChan`.

Two other Goroutines are spawned by `validateInputRecords` which read from the `validChan` & `invalidChan` channels and create an ***InputRecord*** array for each, both of these arrays are return when all Goroutines finish their work.

This final change improved average wall execution time from 15.1s to 12.1s another significant improvement. Another point to make that can be seen in the profiling time results is that the `user` time grows as the program becomes more parallelised (as more CPU time is used across multiple CPUs) but the `real` time falls as less real-world time passes when running the program.

The final profiling results are in `submission/profiling/5__TASK_3_PARA/`

> Written with [StackEdit](https://stackedit.io/).
