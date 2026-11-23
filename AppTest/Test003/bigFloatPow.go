package main

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	ePref "github.com/MikeAustin71/errpref/v2"
)

func main() {

	fileName := "floatmechops/AppTest/Test003/bigFloatPow.go"

	lineBreakStr := strings.Repeat("-", 50)

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		fileName, // sourceFile
		"",       // parentType (leave "" if none)
		"main()") // thisFunction

	floatNumStr := "10.81233"

	exponentValue := uint(8)

	requiredFractionalDigits := uint(14)

	//                    Decimal            1         2         3
	//                    Digits    12345678901234567890123456789012345
	expectedResultStr := "186790307.75412594147937"

	fmt.Printf("\n%v\n"+
		"Function: main\n"+
		"File: %v\n"+
		"%v\n"+
		"floatNumStr: %v\n"+
		"exponentValue: %v\n"+
		"requiredFractionalDigits: %v\n"+
		"%v\n"+
		"Starting Calculation Run\n"+
		"%v\n",
		lineBreakStr,
		fileName,
		lineBreakStr,
		floatNumStr,
		exponentValue,
		requiredFractionalDigits,
		lineBreakStr,
		lineBreakStr)

	calcValue, err := testPrecision01(floatNumStr, exponentValue, requiredFractionalDigits)

	if err != nil {

		err2 := errGenericCfg.GetErrorFromFunctionCall(
			"",
			"testPrecision01()",
			"",
			"",
			err)

		fmt.Printf("Function: main\n"+
			"Status: Error\n"+
			"%v", err2)

		return
	}

	statusStr := "SUCCESS!"

	if calcValue != expectedResultStr {
		statusStr = "FAILURE!"
	}

	fmt.Printf("%v\n"+
		"Function: main\n"+
		"File: %v\n"+
		"%v\n"+
		"Status: %v!\n"+
		"%v\n"+
		"  Expected Value: %v\n"+
		"Calculated Value: %v\n"+
		"%v\n",
		lineBreakStr,
		fileName,
		lineBreakStr,
		statusStr,
		lineBreakStr,
		expectedResultStr,
		calcValue,
		lineBreakStr)

	return
}

func testPrecision01(floatNumStr string, exponent uint, requiredFractionalDigits uint) (calcValue string, err error) {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"floatmechops/AppTest/Test003/bigFloatPow.go", // sourceFile
		"",                  // parentType (leave "" if none)
		"testPrecision01()") // thisFunction

	integerDigits, fractionalDigits, valueIsZero, valueIsOne, err := analyzeFloatNumStr(floatNumStr)

	if err != nil {

		return calcValue, errGenericCfg.GetErrorFromFunctionCall(
			"",
			"analyzeFloatNumStr()",
			"",
			"",
			err)
	}

	lineBreakStr := strings.Repeat("-", 50)

	if valueIsZero {

		numOfZeros := strings.Repeat("-", int(requiredFractionalDigits))

		calcValue = "1." + numOfZeros

		fmt.Printf("\n%v\nFunction: testPrecision01\n"+
			"floatNumStr = %v\n"+
			"exponent = %v\n"+
			"requiredFractionalDigits = %v\n"+
			"Computed Value = %v\n"+
			"%v\n\n",
			lineBreakStr, floatNumStr, exponent, requiredFractionalDigits, calcValue, lineBreakStr)

		return calcValue, nil

	}

	if valueIsOne {

		numOfZeros := ""
		calcValue = floatNumStr

		if requiredFractionalDigits > fractionalDigits {
			numOfZeros = strings.Repeat("-", int(requiredFractionalDigits-fractionalDigits))
		}

		calcValue += numOfZeros

		fmt.Printf("\n%v\nFunction: testPrecision01\n"+
			"floatNumStr = %v\n"+
			"exponent = %v\n"+
			"requiredFractionalDigits = %v\n"+
			"Computed Value = %v\n"+
			"%v\n\n",
			lineBreakStr, floatNumStr, exponent, requiredFractionalDigits, calcValue, lineBreakStr)

		return calcValue, nil

	}

	requiredFloatPrecision, err := calculateFloatPrecisionExponentMultiplication(
		integerDigits,
		fractionalDigits,
		requiredFractionalDigits,
		exponent)

	if err != nil {

		return calcValue, errGenericCfg.GetErrorFromFunctionCall(
			"",
			"calculateFloatPrecisionExponentMultiplication()",
			"",
			"",
			err)
	}

	baseMultiplier, ok := new(big.Float).SetPrec(requiredFloatPrecision).SetMode(big.ToNearestAway).SetString(floatNumStr)

	resultMultiplier, ok := new(big.Float).SetPrec(requiredFloatPrecision).SetMode(big.ToNearestAway).SetString(floatNumStr)

	if ok == false {

		return calcValue, errGenericCfg.GetErrorFromFunctionCall(
			"",
			fmt.Sprintf("ok=new(big.Float).SetString(%s)", floatNumStr),
			"(big.Float).SetString(floatNumStr) failed",
			"Error: 'ok' returned 'false'",
			err)

	}

	adjustedExponent := exponent - 1

	for range adjustedExponent {

		resultMultiplier = resultMultiplier.Mul(resultMultiplier, baseMultiplier)

	}

	calcValue = resultMultiplier.Text('f', int(requiredFractionalDigits))

	fmt.Printf("\n%v\nFunction: testPrecision01\n"+
		"Status: Calculation Completed Successfully\n"+
		"%v\n"+
		"floatNumStr = %v\n"+
		"exponent = %v\n"+
		"requiredFractionalDigits = %v\n"+
		"Computed Value = %v\n"+
		"%v\n\n",
		lineBreakStr,
		lineBreakStr,
		floatNumStr,
		exponent,
		requiredFractionalDigits,
		calcValue,
		lineBreakStr)

	return calcValue, nil
}

// HELPER FUNCTIONS

// analyzeFloatNumStr
//
// Subjects a floating point number string to analysis and validation.
//
// If successful, the number of integer numerical digits and fractional
// numeric digits are returned.
//
// If the floating point number string is found to be invalid, an error is
// returned.
//
// Input Parameters
// ----------------
//
//	floatNumStr  string
//	  This parameter must be formatted as a number string containing one or
//	  more numeric digits. An included decimal point ('.') is optional.
//
//	  If any characters other than numeric digits or a single decimal point
//	  ('.') are identified in floatNumStr, an error will be returned.
//
// Return Parameters
// -----------------
//
//	numIntegerDigits uint
//	  The number of integer digits contained in the input parameter 'floatNumStr'.
//
//	numFractionalDigits uint
//	  The number of fractional digits to the right of the decimal point contained in
//	   the input parameter 'floatNumStr'.
//
//	err error
func analyzeFloatNumStr(floatNumStr string) (
	numIntegerDigits uint,
	numFractionalDigits uint,
	valueIsZero bool,
	valueIsOne bool,
	err error) {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"floatmechops/AppTest/Test003/bigFloatPow.go", // sourceFile
		"",                     // parentType (leave "" if none)
		"analyzeFloatNumStr()") // thisFunction

	if len(floatNumStr) == 0 {

		return 0, 0, false, false,
			errGenericCfg.GetErrorInputParameterOutOfRange(
				"floatNumStr",
				"",
				"",
				"",
				"Input parameter 'floatNumStr' must contain at least one numeric digit",
				"ERROR: Input parameter 'floatNumStr' is an empty string")
	}

	foundDecimal := false

	var integerStr, fractionalStr string

	for _, char := range floatNumStr {

		if char == '.' {

			if foundDecimal == true {

				return 0, 0, false, false,
					errGenericCfg.GetErrorNativeToFunction(
						fmt.Sprintf("floatNumStr= %s", floatNumStr),
						"",
						"",
						"",
						"Input parameter 'floatNumStr' contains more than one decimal point ('.')",
						"ERROR: Input parameter 'floatNumStr' contains invalid characters")

			}

			foundDecimal = true
			continue
		}

		if char >= '0' && char <= '9' {

			if foundDecimal {

				fractionalStr += string(char)

				numFractionalDigits++

			} else {

				numIntegerDigits++

				integerStr += string(char)
			}

			continue
		}

		return 0, 0, false, false,
			errGenericCfg.GetErrorNativeToFunction(
				fmt.Sprintf("floatNumStr= %s", floatNumStr),
				"",
				"",
				"",
				"Input parameter 'floatNumStr' contains non-numeric characters",
				"ERROR: Input parameter 'floatNumStr' contains invalid characters")
	}

	if integerStr == "0" && fractionalStr == "0" {
		valueIsZero = true
	}

	if integerStr == "1" && fractionalStr == "0" {
		valueIsOne = true
	}

	return numIntegerDigits, numFractionalDigits, valueIsZero, valueIsOne, nil
}

// calculateFloatPrecisionExponentMultiplication
//
// Calculates the big.Float precision required to accurately
// represent the value of a number, comprised of specified integer
// and fractional digits, to the power of a specified exponent.
//
// Input Parameters
// ----------------
//
//	baseIntegerDigits uint
//	  The number integer digits in the base number being raised to
//	  the power of 'exponent'
//
//	baseFractionalDigits uint
//	  The number fractional digits in the base number being raised
//	  to the power of 'exponent'
//
//	requiredFractionalDigits uint
//	  The number for fractional digits of accuracy in the number
//	  generated by raising the base number to the power of
//	  'exponent'
//
//	exponent uint
//
//	  The numeric value of the exponent to which the base number
//	  will be raised by exponentiation.
//
// Return Parameters
// -----------------
//
//	bigFloatPrecision uint
//	  The precision value required for the big.Float number which
//	  will accurately represent the result of base number raised
//	  to the power of exponent.
//
//	err error
func calculateFloatPrecisionExponentMultiplication(
	baseIntegerDigits uint,
	baseFractionalDigits uint,
	requiredFractionalDigits uint,
	exponent uint) (bigFloatPrecision uint, err error) {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"floatmechops/AppTest/Test003/bigFloatPow.go", // sourceFile
		"", // parentType (leave "" if none)
		"calculateFloatPrecisionExponentMultiplication()") // thisFunction

	totalFinalDigits := uint(0)

	totalStartingDigits := baseIntegerDigits + baseFractionalDigits

	if totalStartingDigits == 0 {
		return 0, errGenericCfg.GetErrorNativeToFunction(
			"Input parameter baseIntegerDigits = 0 and baseFractionalDigits = 0",
			"",
			"",
			"",
			"There are zero baseIntegerDigits and zero baseFractionalDigits",
			"ERROR: The total number of Integer and Fractional digits is zero")
	}

	additionalDigitsForMultiplication := exponent * uint(6)

	totalFinalDigits = totalStartingDigits + additionalDigitsForMultiplication

	deltaRequiredActualFracDigits := uint(0)

	if requiredFractionalDigits > totalStartingDigits {

		deltaRequiredActualFracDigits = requiredFractionalDigits - totalStartingDigits

		totalFinalDigits += deltaRequiredActualFracDigits
	}

	baseBits := float64(totalFinalDigits) * 3.32192809488736

	// Add safety margin for exponentiation
	safety := 64.0

	// Round up to nearest whole bit
	bigFloatPrecision = uint(math.Ceil(baseBits + safety))

	return bigFloatPrecision, nil
}
