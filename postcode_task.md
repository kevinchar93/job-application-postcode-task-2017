# Technical Task

## Part 1 - Postcode validation

Write code that will validate UK postcodes.

You are given a regular expression that validates postcodes (shown in verbose form below):

        (GIR\s0AA) |
        (
            # A9 or A99 prefix
            ( ([A-PR-UWYZ][0-9][0-9]?) |
                 # AA99 prefix with some excluded areas
                (([A-PR-UWYZ][A-HK-Y][0-9](?<!(BR|FY|HA|HD|HG|HR|HS|HX|JE|LD|SM|SR|WC|WN|ZE)[0-9])[0-9]) |
                 # AA9 prefix with some excluded areas
                 ([A-PR-UWYZ][A-HK-Y](?<!AB|LL|SO)[0-9]) |
                 # WC1A prefix
                 (WC[0-9][A-Z]) |
                 (
                    # A9A prefix
                   ([A-PR-UWYZ][0-9][A-HJKPSTUW]) |
                    # AA9A prefix
                   ([A-PR-UWYZ][A-HK-Y][0-9][ABEHMNPRVWXY])
                 )
                )
              )
              # 9AA suffix
            \s[0-9][ABD-HJLNP-UW-Z]{2}
            )

Write unit tests and implement the regular expression to check the validity of the following postcodes:

| Postcode | Expected problem |
|----------|---------|
| $%± ()()| Junk |
| XX XXX | Invalid |
| A1 9A | Incorrect inward code length |
| LS44PL | No space |
| Q1A 9AA| 'Q' in first position |
| V1A 9AA| 'V' in first position|
| X1A 9BB| 'X' in first position|
| LI10 3QP | 'I' in second position |
| LJ10 3QP | 'J' in second position |
| LZ10 3QP | 'Z' in second position |
| A9Q 9AA | 'Q' in third position with 'A9A' structure|
| AA9C 9AA | 'C' in fourth position with 'AA9A' structure	|
|FY10 4PL| Area with only single digit districts|
|SO1 4QQ|  Area with only double digit districts|
| EC1A 1BB | None |
| W1A 0AX | None |
| M1 1AE | None |
| B33 8TH | None |
| CR2 6XH | None |
| DN55 1PT | None |
| GIR 0AA | None |
|SO10 9AA    |None|
|FY9 9AA      |None|
|WC1A 9AA    |None|

Please read (https://en.wikipedia.org/wiki/Postcodes_in_the_United_Kingdom#Validation), does the regular expression validate all UK postcode cases? 

## Part 2 - Bulk import

Imagine you are migrating demographic data, which includes postcode information and you need to check the data for invalid postcodes.

Write bulk import code that will validate the postcodes in the data file of around 2 million postcodes ([download from google drive](https://drive.google.com/file/d/0BwxZ38NLOGvoTFE4X19VVGJ5NEk/view?usp=sharing)) named `import_data.csv.gz` and report on the `row_id` where validation fails, the structure of `import_data.csv.gz` is shown below:

| row_id | postcode |
|--------|----------|
| 1 | AABC 123|
| 2 | AACD 4PQ|
|...|...|

If you need to untar the file, that is acceptable.

At the end of running the bulk import you should produce a file named, `failed_validation.csv` with the same columns as above.

## Part 3 - Performance engineering

Modify the code in **Part 2** to produce two files:

    succeeded_validation.csv
    failed_validation.csv
    
The postcodes in the two files need to be ordered as per the `row_id`, in ascending numeric order.

Analyse the performance of your solution and make an attempt to optimise the performance of the operation (in terms of overall 'wall' time taken).  Describe how you improved the performance of the code, and how you measured the impact of your changes.

It is acceptable to not use the regular expression (or different regular expression(s)) for this part of the task, but the output in terms of the correctness of the validation needs to match the critieria in **Part 1**.

# Constraints / instructions

- Use a language of your choosing.
- Use only standard libraries (unless your language of choice doesn't have a regular expression library).
- Include instructions on how to run your solution in a markdown formatted file in the root of your solution named `README.md`
- Include notes on your analysis (e.g. where you have found the regular expression provided doesn't deal with all UK postcode edge cases) either as notes within the code, or in another file `ANALYSIS.md` if you prefer.
- Submit the task to *******@gmail.com
- Do not include compiled code or binaries.
- Do not include any output files or any postcode test files.
