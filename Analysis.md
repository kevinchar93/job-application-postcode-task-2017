Profiling Function call % of CPU time

(95.61 %) main func
---------------------------------------------------------
(45.37 %) csv.Reader.Read
(19.77 %) main.RegexValidatorGroup.GroupIsStringValid
(8.83 %)  main.writeOutputFiles
(8.35 %) sort.Sort
(6.97 %) main.NewImportRecord
(0.72 %) runtime.growslice


What was done in regex bit
-----------------------------
analysed regex to understand it, broke it down into sections
tried it out using some online tools
noticed some sections of it needed minor changes to not match invalid postcodes that can be found in longer postcodes
used caret (^) char to solve this issue on the regex groups for the smaller postcode prefixes
decided to use go for implementation language
go is fast, siutable for small tools like this, strongly types, compiled (fast), massive standard lib, good support for concurrency
noticed that go's regex lib did not support lookahead/lookbehind in any form
created a work around by creating a new simple but powerful type RegexValidator & RegexValidatorGroup
created tests to test correctness and unit tests laid out in the brief

What was done in profiling
-----------------------------
ran profiler using profiling library : https://github.com/pkg/profile
this library abstracts profiling package that Go has makes it a single line that needs adding to the file
used the pprof tool that comes with go lang to create visulisations of profiler results .png files
create profile of cpu a few times and a memory profile to check memory usage was reasonable

checked over profiling results to find out which bits of the program were taking the most time
made a list of functions an ordered according to % of cpu time each Took
looked into making functions that take most cpu time run faster csv.Reader.Read in particular
changed csv.Reader.Read to a bufio.Reader.ReadString function call reducing time it takes to read from 47% to about 39%
changed program to use parallel pipeline for reading -> creating import records -> validating records
validation records taking up most cpu time
parallelized that section letting it run on 3 goroutines and have 2 other routines adding up the results
significant improvements in wall time made
