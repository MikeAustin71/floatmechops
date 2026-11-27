package floatmech

import (
	"fmt"
	"math/big"
	"math/bits"
	"sync"

	ePref "github.com/MikeAustin71/errpref/v2"
)

const (
	// maxSciExp caps the absolute value of a scientific notation exponent (e.g., 1e10000)
	// to prevent big.Rat.SetString from allocating huge buffers.
	maxSciExp = 10000

	// maxFracDigits caps requested fractional digits to prevent int(reqFracDigits)
	// overflow and excessive memory allocation.
	maxFracDigits = 1_000_000
)

type DecimalFloat struct {
	mu                     sync.Mutex
	val                    *big.Float
	targetFractionalDigits uint
	prec                   uint
	mode                   big.RoundingMode
}

// NewFromFloatNumStr constructs an immutable Decimal sized to
// preserve at least targetFractionalDigits decimal places for
// the value's magnitude.
//
// The input numeric string value is padded with a fixed number
// of trailing '0' characters prior to parsing, raising the
// parsed value's magnitude so that the derived working
// precision retains the requested number of fractional digits.
//
// Input Parameters
// ----------------
//
// floatNumStr string
//
//	Numeric string value used to construct the Decimal. If
//	'floatNumStr' is an empty string, an error is returned.
//
// targetFractionalDigits uint
//
//	Number of base-10 fractional digits to be preserved by the
//	constructed Decimal.
//
// mode big.RoundingMode
//
//	Rounding mode applied when the parsed value is rounded to
//	the working precision.
//
// Return Parameters
// -----------------
//
// *DecimalFloat
//
//	Pointer to the newly constructed Decimal object.
//
// error
func (df *DecimalFloat) NewFromFloatNumStr(floatNumStr string, targetFractionalDigits uint, mode big.RoundingMode) (*Decimal, error) {

	// NewErrorGenericConfig()
	//
	// func (egc *ErrorGenericFactory) NewErrorGenericConfig(
	//    sourceFile string,
	//    parentType string,
	//    thisFunction string) *ErrorGenericFactory
	//
	// Add identification information for this function
	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"decimalFloat_0001000.go",
		"DecimalFloat",
		"NewFromFloatNumStr()")

	if df == nil {

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
		return &Decimal{}, errGenericCfg.GetErrorInputParameterNilPointer(
			"d == nil",
			"",
			"",
			"",
			"Invalid Pointer Receiver",
			"Error: Pointer Receiver 'd' is a nil pointer")
	}

	df.mu.Lock()

	defer df.mu.Unlock()

	constructedDecimal, err := new(numStrConstructorMechanics).constructFromNumStr(
		floatNumStr,
		targetFractionalDigits,
		mode)

	if err != nil {

		// GetErrorFromFunctionCall()
		//
		// func (egc *ErrorGenericFactory) GetErrorFromFunctionCall(
		//    calledFunctionType string,
		//    calledFunction string,
		//    errorContext string,
		//    errorMessage string,
		//    returnedCalledFunctionError error) error
		//
		// Adds error context to errors returned by other
		// Functions.
		err2 := errGenericCfg.GetErrorFromFunctionCall(
			"numStrConstructorMechanics",
			"constructFromNumStr()",
			"",
			"",
			err)

		return &Decimal{}, err2
	}

	return constructedDecimal, nil
}

//		PowDecimal raises baseStr to a non-negative integer exponent (exp) and returns:
//	 	1. A fixed-point decimal string rounded to reqFracDigits fractional places.
//	 	2. A high-precision *big.Float representing that exact rounded result.
//
//		Special Cases:
//			- Any valid base raised to exp == 0 (including "0"^0 and "-2.5"^0) returns 1.
//	  	- exp == 1 performs direct rational rounding without an exponentiation loop.
//	  	- Moderate inputs (exp <= 250 and reqFracDigits <= 1000) use exact [big.Rat]
//	    	arithmetic; larger inputs automatically switch to guard-buffered [big.Float]
//	    	exponentiation by squaring.
//
//		Input Parameters
//		----------------
//
//		baseStr string
//			Input Format Requirements (baseStr):
//
//			baseStr must be a non-empty base-10 decimal or scientific notation string with no
//			whitespace, underscores ('_'), or fraction slashes ('/').
//
//			Allowable scientific notation strings follow the structure:
//
//			[Sign][Mantissa][ExponentMarker][ExponentSign][ExponentDigits]
//
//			Where:
//	  		- Sign (optional): '+' or '-' at index 0 only.
//	  		- Mantissa (required): Base-10 digits ('0'-'9') with at most one decimal point
//	    		('.'). Leading or trailing decimal points (e.g., ".5" or "5.") are permitted.
//	  		- ExponentMarker (required for scientific notation): 'e' or 'E'.
//	  		- ExponentSign (optional): '+' or '-' immediately following 'e' or 'E'.
//	  		- ExponentDigits (required): Base-10 integer digits representing a signed
//	    		exponent within [-maxSciExp, maxSciExp] ([-10000, 10000]).
//
//			Valid Examples:
//
//				"-2.5", "0", "-1.25e-10", "2.5E+4", ".5e-3", "1e10000"
//
//			Invalid Examples:
//
//				" -2.5"      (leading/trailing whitespace)
//				"1_000.5e2"  (underscore separators)
//				"1/3"        (rational fraction slash)
//				"12-3.4e5"   (minus sign not at index 0)
//				"1.25e1-0"   (exponent minus sign not immediately after 'e'/'E')
//				"1e10001"    (exponent exceeds [-maxSciExp, maxSciExp])
//
//		error
//			Returns an error if baseStr violates the format or exponent bounds above, or if
//			reqFracDigits exceeds maxFracDigits (1,000,000).
func (df *DecimalFloat) PowDecimal(
	baseStr string, exp uint, reqFracDigits uint) (string, *big.Float, error) {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"floatmech/decimalFloat_0001000.go", // sourceFile
		"DecimalFloat",                      // parentType (leave "" if none)
		"PowDecimal()")                      // thisFunction

	if df == nil {

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
		return "", big.NewFloat(0), errGenericCfg.GetErrorInputParameterNilPointer(
			"d == nil",
			"",
			"",
			"",
			"Invalid Pointer Receiver",
			"Error: Pointer Receiver 'd' is a nil pointer")
	}

	df.mu.Lock()

	defer df.mu.Unlock()

	if reqFracDigits > maxFracDigits {

		return "", big.NewFloat(0),
			errGenericCfg.GetErrorNativeToFunction(
				fmt.Sprintf("reqFracDigits= %v", reqFracDigits),
				"",
				"",
				"",
				"Input parameter 'reqFracDigits' exceeds maximum allowed value.",
				fmt.Sprintf("ERROR: reqFracDigits (%d) exceeds maximum allowed (%d)", reqFracDigits, maxFracDigits))
	}

	if err := new(decimalFloatUtil020).validateBaseFormat(baseStr); err != nil {

		return "", big.NewFloat(0), errGenericCfg.GetErrorFromFunctionCall(
			"",
			"validateBaseFormat(baseStr)",
			fmt.Sprintf("baseStr: %s", baseStr),
			"",
			err)
	}

	dfUtil030 := new(decimalFloatUtil030)

	// Fast path: exp == 0 (includes 0^0 = 1 and (-x)^0 = 1)
	if exp == 0 {

		if _, ok := new(big.Rat).SetString(baseStr); !ok {

			return "", big.NewFloat(0),
				errGenericCfg.GetErrorNativeToFunction(
					"'ok' == false",
					"",
					"",
					"",
					fmt.Sprintf("'baseStr' = %v", baseStr),
					fmt.Sprintf("new(big.Rat).SetString(baseStr) returned 'ok' == false"))
		}

		ratOne := big.NewRat(1, 1)

		resStr := ratOne.FloatString(int(reqFracDigits))

		resFloat, ok := new(big.Float).SetPrec(dfUtil030.calcBits(1, reqFracDigits, 1)).SetString(resStr)

		if !ok {

			return "", big.NewFloat(0), errGenericCfg.GetErrorFromFunctionCall(
				"",
				"new(big.Float).SetPrec(calcBits(1, reqFracDigits, 1)).SetString(resStr)",
				fmt.Sprintf("baseStr: %s", baseStr),
				"",
				fmt.Errorf("new(big.Rat).SetString() returned 'ok' == false"))

		}

		return resStr, resFloat, nil
	}

	// Fast path: exp == 1 (exact rounding with no multiplication loop)
	if exp == 1 {

		baseRat, ok := new(big.Rat).SetString(baseStr)

		if !ok {
			return "", nil, fmt.Errorf("invalid base number string: %q", baseStr)
		}

		resStr := baseRat.FloatString(int(reqFracDigits))

		finalPrec := dfUtil030.calcBits(uint(baseRat.Num().BitLen()), reqFracDigits, 1)

		resFloat, _ := new(big.Float).SetPrec(finalPrec).SetString(resStr)

		return resStr, resFloat, nil
	}

	// TIER 1: Exact Rational Math for moderate exponents & fractional precision
	// exp <= 250 && reqFracDigits <= 1000
	if exp <= 250 && reqFracDigits <= 1000 {

		baseRat, ok := new(big.Rat).SetString(baseStr)

		if !ok {
			return "", nil, fmt.Errorf("invalid base number string: %q", baseStr)
		}

		resRat := dfUtil030.powRat(baseRat, exp)

		resStr := resRat.FloatString(int(reqFracDigits))

		finalPrec := dfUtil030.calcBits(uint(resRat.Num().BitLen()), reqFracDigits, 1)

		resFloat, _ := new(big.Float).SetPrec(finalPrec).SetString(resStr)

		return resStr, resFloat, nil
	}

	// TIER 2: Guard-buffered big.Float for large exponents
	// (exponent value >250) or reqFracDigits > 1,000 fractional digits
	coarse, ok := new(big.Float).SetPrec(64).SetString(baseStr)

	if !ok {
		return "", big.NewFloat(0),
			errGenericCfg.GetErrorNativeToFunction(
				"'ok' == false",
				"",
				"",
				"",
				fmt.Sprintf("'baseStr' = %v", baseStr),
				fmt.Sprintf("new(big.Float).SetPrec(64).SetString(baseStr) returned 'ok' == false"))
	}

	// Estimate integer bits needed using absolute magnitude with overflow protection
	var mant big.Float

	binExp := coarse.MantExp(&mant)

	var estIntBits uint64

	if binExp > 0 {

		hi, lo := bits.Mul64(uint64(binExp), uint64(exp))

		if hi != 0 {
			return "", nil, fmt.Errorf("result magnitude overflows precision limits")
		}

		estIntBits = lo
	}

	workPrec, err := dfUtil030.calcBitsSafe(estIntBits, reqFracDigits, exp)

	if err != nil {
		return "", nil, err
	}

	baseFloat, _ := new(big.Float).SetPrec(workPrec).SetString(baseStr)

	resFloat := dfUtil030.powFloat(baseFloat, exp, workPrec)

	resStr := resFloat.Text('f', int(reqFracDigits))

	// Re-parse the formatted string so the returned big.Float matches resStr
	finalFloat, ok := new(big.Float).SetPrec(workPrec).SetString(resStr)

	if !ok {
		return "", nil, fmt.Errorf("SetString(resStr) - result magnitude overflows precision limits")
	}

	return resStr, finalFloat, nil
}
