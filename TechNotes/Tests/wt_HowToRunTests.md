# Running Tests mathopsgo

Open a command prompt in this directory (*strmech*) and run the
following commands.

## Test Execution WITHOUT Code Coverage
Run this in *strmech* directory:

### Windows Command Without Coverage
  `go test -v > zzzzz_tests.txt`

### Linux Command Without Coverage
`go test -v | tee zzzzz_tests.txt`

This will generate test results in the *strmech* 
directory which are stored in the text file, `zzzzz_tests.txt`. 

## Running Tests with code coverage

First pull down and install the `cover` package.
 
`go get golang.org/x/tools/cmd/cover`
  
Next, follow the test execution protocol.  
  
## Test Execution With Code Coverage
Run this in *strmech* directory:

### Windows Command With Coverage
`go test -cover -v > zzzzz_tests.txt`
`go test -cover -v --race > zzzzz_tests.txt`

### Linux Command With Coverage
`go test -cover -v | tee zzzzz_tests.txt`


## Cover Profile

Generate the code coverage detail. Run this command
in the *strmech* directory:

`go test -coverprofile=zzzzz_coverage.out`


The following provides for code coverage display in your
browser. Run this on the terminal command line and run it
in the *strmech* directory:

`go tool cover -html=zzzzz_coverage.out`