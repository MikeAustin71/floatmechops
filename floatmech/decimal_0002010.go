package floatmech

import (
	"fmt"
	"math/big"
	"strings"

	ePref "github.com/MikeAustin71/errpref/v2"
)

// floatNumStrPadding
//
// Number of trailing '0' characters appended to an
// input numeric string prior to parsing. The padding
// raises the parsed value's magnitude so that the
// derived working precision retains the target
// fractional digits.
const floatNumStrPadding uint = 30

// numStrConstructorMechanics
//
// Provides private methods designed to construct a
// 'Decimal' object from a numeric input string.
type numStrConstructorMechanics struct{}

// constructFromNumStr
//
// Appends the configured padding to the input numeric
// string, parses the padded value at infinite precision,
// derives the working precision for the target number
// of fractional digits, and stores the value rounded
// to that precision using the supplied rounding mode.
//
// Input Parameters
// ----------------
//
//	vNumStr string
//	  Numeric string value used to construct the Decimal.
//
//	targetFractionalDigits uint
//	  Number of base-10 fractional digits to be preserved
//	  by the constructed Decimal.
//
//	mode big.RoundingMode
//	  Rounding mode applied when the parsed value is
//	  rounded to the working precision.
//
// Return Parameters
// -----------------
//
//	*Decimal
//	  Pointer to the newly constructed Decimal object.
//
//	error
func (nscm *numStrConstructorMechanics) constructFromNumStr(
	vNumStr string,
	targetFractionalDigits uint,
	mode big.RoundingMode) (*Decimal, error) {

	// NewErrorGenericConfig()
	//
	// func (egc *ErrorGenericFactory) NewErrorGenericConfig(
	//    sourceFile string,
	//    parentType string,
	//    thisFunction string) *ErrorGenericFactory
	//
	// Add identification information for this function
	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"decimal_0002010.go",
		"numStrConstructorMechanics",
		"constructFromNumStr()")

	if len(vNumStr) == 0 {

		// GetErrorNativeToFunction()
		//
		// func (egc *ErrorGenericFactory) GetErrorNativeToFunction(
		//    parameter string,
		//    path string,
		//    fileName string,
		//    url string,
		//    errorContext string,
		//    errorMessage string) error
		//
		// Formats error message and error context for errors
		// originating in the current function's internal logic.
		return &Decimal{}, errGenericCfg.GetErrorNativeToFunction(
			"vNumStr=\"\"",
			"",
			"",
			"",
			"Input numeric string is empty",
			"Error: Input parameter 'vNumStr' is an empty string")
	}

	var found = false

	for _, char := range vNumStr {

		if char == '.' {
			found = true
			continue
		}

		if char < '0' || char > '9' {

			// GetErrorNativeToFunction()
			//
			// func (egc *ErrorGenericFactory) GetErrorNativeToFunction(
			//    parameter string,
			//    path string,
			//    fileName string,
			//    url string,
			//    errorContext string,
			//    errorMessage string) error
			//
			// Formats error message and error context for errors
			// originating in the current function's internal logic.
			return &Decimal{}, errGenericCfg.GetErrorNativeToFunction(
				fmt.Sprintf("Input parameter 'vNumStr'= %v", vNumStr),
				"",
				"",
				"",
				"Input parameter 'vNumStr' eligible characters are numbers and a decimal point.",
				"Error: Input parameter 'vNumStr' contains non-numeric characters.")
		}
	}

	var paddedNumStr string

	// Append the configured padding to the input numeric
	// string prior to parsing.
	if !found {
		// If paddedNumStr does not contain a decimal point, add one.
		paddedNumStr = vNumStr + "." + strings.Repeat("0", int(floatNumStrPadding))

	} else {

		paddedNumStr = vNumStr + strings.Repeat("0", int(floatNumStrPadding))
	}

	// Parse the padded numeric string at infinite precision.
	// At infinite precision the parsed value is exact.
	exact, ok := new(big.Float).SetString(paddedNumStr)

	if !ok {

		// GetErrorNativeToFunction()
		//
		// func (egc *ErrorGenericFactory) GetErrorNativeToFunction(
		//    parameter string,
		//    path string,
		//    fileName string,
		//    url string,
		//    errorContext string,
		//    errorMessage string) error
		//
		// Formats error message and error context for errors
		// originating in the current function's internal logic.
		return &Decimal{}, errGenericCfg.GetErrorNativeToFunction(
			fmt.Sprintf("paddedNumStr=%s", paddedNumStr),
			"",
			"",
			"",
			"Failed to parse input numeric string",
			"Error: Input parameter 'vNumStr' is not a valid numeric string")
	}

	// Derive the exponent of the parsed value. For a zero
	// value the exponent is zero.
	var exp2 int

	if exact.Sign() != 0 {
		exp2 = exact.MantExp(nil)
	}

	// The parsed value's magnitude is elevated by the
	// padding, so the target fractional digit count is
	// elevated by the same amount when computing the
	// working precision.
	effectiveTfd := targetFractionalDigits + floatNumStrPadding

	prec := defaultConverter().toPrec(effectiveTfd, exp2)

	// Round the exact parsed value to the working precision
	// using the supplied rounding mode.
	val := new(big.Float).SetPrec(prec).SetMode(mode).Set(exact)

	return &Decimal{
		val:                    val,
		targetFractionalDigits: targetFractionalDigits,
		prec:                   prec,
		mode:                   mode,
	}, nil
}
