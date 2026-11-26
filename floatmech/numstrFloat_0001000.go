package floatmech

import (
	"sync"

	ePref "github.com/MikeAustin71/errpref/v2"
)

type NumStrStatusDto struct {
	NumIntegerDigits    uint
	IntegerRunes        []rune
	NumFractionalDigits uint
	FractionalRunes     []rune
	ValueIsNegative     bool
	ValueIsZero         bool
	ValueIsOne          bool
}

// NumStrFloat encapsulates methods supporting number
// strings comprised of floating point numerical values.
type NumStrFloat struct {
	mu sync.Mutex
}

//		AnalyzeNumStrFloat subjects a floating point number string
//		to analysis and validation.
//
//	 If successful, the number of integer numerical digits and
//	 fractional numeric digits are returned in a NumStrStatusDto
//	 object.
//
//	 If the floating point number string is found to be invalid,
//	 an error is returned.
//
//	 Input Parameters
//	 ----------------
//
//	 floatNumStr  string
//	   This parameter must be formatted as a number string containing one or
//		  more numeric digits. An included decimal point ('.') is optional.
//
//		  If any characters other than numeric digits or a single decimal point
//		  ('.') are identified in floatNumStr, an error will be returned.
//
//	 Return Parameters
//	 -----------------
//
//		numStrStatus NumStrStatusDto
//	   If successful, this method returns a NumStrStatusDto object
//	   encapsulating detailed information on the floating point value
//	   passed as input parameter floatNumStr. This information includes
//	   number of integer and fractional digits in floatNumStr.
//
//		err error
func (nsf *NumStrFloat) AnalyzeNumStrFloat(floatNumStr string) (numStrStatus NumStrStatusDto, err error) {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"floatmech/numstrFloat_0001000.go", // sourceFile
		"NumStrFloat",                      // parentType (leave "" if none)
		"AnalyzeNumStrFloat()")             // thisFunction

	if nsf == nil {

		// GetErrorInputParameterNilPointer()
		//
		// func (egc *ErrorGenericFactory) GetErrorInputParameterNilPointer(
		//    parameter string,
		//    path string,
		//    fileName string,
		//    url string,
		//    errorContext string,
		//    errorMessage string) error
		//
		// Formats error message and error context for when nil pointers
		// are passed as Function input parameters.
		return NumStrStatusDto{}, errGenericCfg.GetErrorInputParameterNilPointer(
			"nsf == nil",
			"",
			"",
			"",
			"Invalid Pointer Receiver",
			"Error: Pointer Receiver 'nsf' is a nil pointer")
	}

	nsf.mu.Lock()

	defer nsf.mu.Unlock()

	numStrStatus, err = new(numStrFloatF0002000).analyzeNumStrFloat(floatNumStr)

	if err != nil {
		return NumStrStatusDto{}, errGenericCfg.GetErrorFromFunctionCall(
			"numStrFloatF0002000",
			"analyzeNumStrFloat",
			"",
			"",
			err)
	}

	return numStrStatus, nil
}
