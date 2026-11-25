package main

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	ePref "github.com/MikeAustin71/errpref/v2"
)

type numStrStatusDto struct {
	NumIntegerDigits    uint
	IntegerRunes        []rune
	NumFractionalDigits uint
	FractionalRunes     []rune
	ValueIsNegative     bool
	ValueIsZero         bool
	ValueIsOne          bool
}

func main() {

	lineBreakStr := strings.Repeat("-", 50)

	fileName := "floatmechops/AppTest/Test003/bigFloatPow.go"
	funcName := "main()"
	targetFuncName := "TestRunner02()"

	fmt.Printf("\n%s\n"+
		"File Name: %s\n"+
		" Function: %s\n"+
		"%s\n"+
		"Starting Run\n"+
		"Running Function: %s\n"+
		"%s\n", lineBreakStr, fileName, funcName, lineBreakStr, targetFuncName, lineBreakStr)

	TestRunner02(fileName)

	fmt.Printf("%s\n"+
		"File Name: %s\n"+
		" Function: %s\n"+
		"%s\n"+
		"Function Completed!\n"+
		"%s\n\n", lineBreakStr, fileName, funcName, lineBreakStr, lineBreakStr)
}

func TestRunner01(homeFileName string) {

	lineBreakStr := strings.Repeat("-", 50)
	funcName := "TestRunner01()"
	targetFuncName := "testPrecision01()"

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		homeFileName, // sourceFile
		"",           // parentType (leave "" if none)
		funcName)     // thisFunction

	floatNumStr := "10.81233"

	exponentValue := uint(8)

	requiredFractionalDigits := uint(14)

	//                    Decimal            1         2         3
	//                    Digits    12345678901234567890123456789012345
	expectedResultStr := "186790307.75412594147937"

	fmt.Printf("\n%v\n"+
		"Function: %v\n"+
		"File: %v\n"+
		"%v\n"+
		"floatNumStr: %v\n"+
		"exponentValue: %v\n"+
		"requiredFractionalDigits: %v\n"+
		"%v\n"+
		"Running Function: %v\n"+
		"Starting Calculation Run\n"+
		"%v\n",
		lineBreakStr,
		funcName,
		homeFileName,
		lineBreakStr,
		floatNumStr,
		exponentValue,
		requiredFractionalDigits,
		lineBreakStr,
		targetFuncName,
		lineBreakStr)

	calcValue, err := testPrecision01(floatNumStr, exponentValue, requiredFractionalDigits)

	if err != nil {

		err2 := errGenericCfg.GetErrorFromFunctionCall(
			"",
			"testPrecision02()",
			"",
			"",
			err)

		fmt.Printf("Function: %s\n"+
			"Status: Error\n"+
			"%v", funcName, err2)

		return
	}

	statusStr := "SUCCESS!"

	if calcValue != expectedResultStr {
		statusStr = "FAILURE!"
	}

	fmt.Printf("%v\n"+
		"Function Completed: %v\n"+
		"File: %v\n"+
		"%v\n"+
		"Status: %v!\n"+
		"%v\n"+
		"  Expected Value: %v\n"+
		"Calculated Value: %v\n"+
		"%v\n",
		lineBreakStr,
		funcName,
		homeFileName,
		lineBreakStr,
		statusStr,
		lineBreakStr,
		expectedResultStr,
		calcValue,
		lineBreakStr)

	return
}

func TestRunner02(homeFileName string) {

	lineBreakStr := strings.Repeat("-", 50)
	funcName := "TestRunner02()"
	targetFuncName := "testPrecision02()"

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		homeFileName, // sourceFile
		"",           // parentType (leave "" if none)
		funcName)     // thisFunction

	floatNumStr := "10.81233"

	exponentValue := uint(8)

	requiredFractionalDigits := uint(14)

	//                    Decimal            1         2         3
	//                    Digits    12345678901234567890123456789012345
	expectedResultStr := "186790307.75412594147937"

	fmt.Printf("\n%v\n"+
		"Function: %v\n"+
		"File: %v\n"+
		"%v\n"+
		"floatNumStr: %v\n"+
		"exponentValue: %v\n"+
		"requiredFractionalDigits: %v\n"+
		"%v\n"+
		"Running Function: %v\n"+
		"Starting Calculation Run\n"+
		"%v\n",
		lineBreakStr,
		funcName,
		homeFileName,
		lineBreakStr,
		floatNumStr,
		exponentValue,
		requiredFractionalDigits,
		lineBreakStr,
		targetFuncName,
		lineBreakStr)

	calcValue, err := testPrecision02(floatNumStr, exponentValue, requiredFractionalDigits)

	if err != nil {

		err2 := errGenericCfg.GetErrorFromFunctionCall(
			"",
			"testPrecision02()",
			"",
			"",
			err)

		fmt.Printf("Function: %s\n"+
			"Status: Error\n"+
			"%v", funcName, err2)

		return
	}

	statusStr := "SUCCESS!"

	if calcValue != expectedResultStr {
		statusStr = "FAILURE!"
	}

	fmt.Printf("%v\n"+
		"Function Completed: %v\n"+
		"File: %v\n"+
		"%v\n"+
		"Status: %v!\n"+
		"%v\n"+
		"  Expected Value: %v\n"+
		"Calculated Value: %v\n"+
		"%v\n",
		lineBreakStr,
		funcName,
		homeFileName,
		lineBreakStr,
		statusStr,
		lineBreakStr,
		expectedResultStr,
		calcValue,
		lineBreakStr)

	return
}

func testPrecision01(floatNumStr string, exponent uint, requiredFractionalDigits uint) (calcValue string, err error) {

	functionName := "testPrecision01"
	lineBreakStr := strings.Repeat("-", 50)

	fmt.Printf("\n\n%v\n"+
		"Running Function: %s\n"+
		"%v\n\n", lineBreakStr, functionName, lineBreakStr)

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"floatmechops/AppTest/Test003/bigFloatPow.go", // sourceFile
		"",                  // parentType (leave "" if none)
		"testPrecision01()") // thisFunction

	numStrStatus, err := analyzeFloatNumStr(floatNumStr)

	if err != nil {

		return calcValue, errGenericCfg.GetErrorFromFunctionCall(
			"",
			"analyzeFloatNumStr()",
			"",
			"",
			err)
	}

	if numStrStatus.ValueIsZero {

		numOfZeros := strings.Repeat("0", int(requiredFractionalDigits))

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

	if numStrStatus.ValueIsOne {

		numOfZeros := ""
		calcValue = floatNumStr

		if requiredFractionalDigits > numStrStatus.NumFractionalDigits {
			numOfZeros = strings.Repeat("-", int(requiredFractionalDigits-numStrStatus.NumFractionalDigits))
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

	requiredFloatPrecision, err := calculateBigFloatPrecisionExponentMultiplication(
		numStrStatus.NumIntegerDigits,
		numStrStatus.NumFractionalDigits,
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

func testPrecision02(floatNumStr string, exponent uint, requiredFractionalDigits uint) (calcValueStr string, err error) {

	functionName := "testPrecision02"
	lineBreakStr := strings.Repeat("-", 50)

	fmt.Printf("\n%v\n"+
		"Running Function: %s\n"+
		"%v\n", lineBreakStr, functionName, lineBreakStr)

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"floatmechops/AppTest/Test003/bigFloatPow.go", // sourceFile
		"",                // parentType (leave "" if none)
		"testPrecision02") // thisFunction

	calcValue, err := multiplyFloatNumStrForExponentiation(
		floatNumStr, exponent, requiredFractionalDigits)

	if err != nil {

		return calcValueStr, errGenericCfg.GetErrorFromFunctionCall(
			"",
			"multiplyFloatNumStrForExponentiation(floatNumStr, exponent, requiredFractionalDigits)",
			fmt.Sprintf("floatNumStr= %s", floatNumStr),
			"",
			err)
	}

	calcValueStr = calcValue.Text('f', int(requiredFractionalDigits))

	return calcValueStr, nil
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
	numStrStatus numStrStatusDto, err error) {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"floatmechops/AppTest/Test003/bigFloatPow.go", // sourceFile
		"",                     // parentType (leave "" if none)
		"analyzeFloatNumStr()") // thisFunction

	if len(floatNumStr) == 0 {

		return numStrStatus,
			errGenericCfg.GetErrorInputParameterOutOfRange(
				"floatNumStr",
				"",
				"",
				"",
				"Input parameter 'floatNumStr' must contain at least one numeric digit",
				"ERROR: Input parameter 'floatNumStr' is an empty string")
	}

	foundDecimal := false

	for i, char := range floatNumStr {

		if char == '-' {

			if i == 0 {

				numStrStatus.ValueIsNegative = true

				continue

			} else {

				return numStrStatus,
					errGenericCfg.GetErrorNativeToFunction(
						fmt.Sprintf("floatNumStr= %s", floatNumStr),
						"",
						"",
						"",
						"'floatNumStr' contains a minus sign (-), but NOT at the beginning of the number string.",
						"ERROR: Input parameter 'floatNumStr' contains an invalid minus sigh (-).")
			}
		}

		if char == '.' {

			if foundDecimal == true {

				return numStrStatus,
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

				numStrStatus.FractionalRunes = append(numStrStatus.FractionalRunes, char)

				numStrStatus.NumFractionalDigits++

			} else {

				// Strip off leading zeros
				if char == '0' && len(numStrStatus.IntegerRunes) == 0 {
					continue
				}

				numStrStatus.IntegerRunes = append(numStrStatus.IntegerRunes, char)

				numStrStatus.NumIntegerDigits++
			}

			continue
		}

		return numStrStatusDto{},
			errGenericCfg.GetErrorNativeToFunction(
				fmt.Sprintf("floatNumStr= %s", floatNumStr),
				"",
				"",
				"",
				"Input parameter 'floatNumStr' contains non-numeric characters",
				"ERROR: Input parameter 'floatNumStr' contains invalid characters")
	}

	if len(numStrStatus.IntegerRunes) == 0 &&
		len(numStrStatus.FractionalRunes) == 0 {

		return numStrStatusDto{},
			errGenericCfg.GetErrorNativeToFunction(
				fmt.Sprintf("floatNumStr= %s", floatNumStr),
				"",
				"",
				"",
				"No valid numerical digits were found in 'floatNumStr'",
				"ERROR: Input parameter 'floatNumStr' contains invalid characters")
	}

	if string(numStrStatus.IntegerRunes) == "0" && string(numStrStatus.FractionalRunes) == "0" {
		numStrStatus.ValueIsZero = true
	}

	if string(numStrStatus.IntegerRunes) == "1" && string(numStrStatus.FractionalRunes) == "0" {
		numStrStatus.ValueIsOne = true
	}

	return numStrStatus, nil
}

// calculateBigFloatPrecisionExponentMultiplication
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
func calculateBigFloatPrecisionExponentMultiplication(
	baseIntegerDigits uint,
	baseFractionalDigits uint,
	requiredFractionalDigits uint,
	exponent uint) (bigFloatPrecision uint, err error) {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"floatmechops/AppTest/Test003/bigFloatPow.go", // sourceFile
		"", // parentType (leave "" if none)
		"calculateFloatPrecisionExponentMultiplication()") // thisFunction

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

	projectedIntDigits := calculateUintToMultipleOfThree(baseIntegerDigits)

	projectedFractionalDigits := calculateUintToMultipleOfThree(baseFractionalDigits)

	projectedIntDigits = projectedIntDigits * exponent

	/*
		projectedFractionalDigits = projectedFractionalDigits * exponent

		if requiredFractionalDigits > projectedFractionalDigits {

			projectedFractionalDigits = requiredFractionalDigits

		}
	*/

	projectedFractionalDigits = max(requiredFractionalDigits, projectedFractionalDigits*exponent)

	totalFinalDigits := projectedIntDigits + projectedFractionalDigits

	totalFinalDigits += uint(10) // Safety margin in extra digits

	bigFloatPrecision = calculateBigFloatPrecision(totalFinalDigits)

	return bigFloatPrecision, nil
}

// calculateUintToMultipleOfThree
//
// Receives an 'uint' value and modifies that value such that the value returned value
// is a multiple of 3. This is useful in computing maximum precision requires for big.Float
// values.
//
// Input Parameters
// ================
//
//	baseUint uint
//	  A base unsigned integer value.
//
// Return Parameters
// =================
//
//	baseUnitToMultipleOf3 uint
//	  The value of inputParameter 'baseUint' returned as a multiple of 3.
//	  If the baseUint value is not initially a multiple of 3, that value will
//	  be increased by one (1) until is a multiple of 3 before it is returned to
//	  the calling function.
func calculateUintToMultipleOfThree(baseUint uint) (baseUnitToMultipleOf3 uint) {

	var modulo uint

	baseUnitToMultipleOf3 = baseUint

	for {

		modulo = baseUint % 3

		if modulo == 0 {
			break
		}

		baseUint++

	}

	return baseUnitToMultipleOf3
}

// calculateBigFloatPrecisionOnFloatNumStr
//
// Receives a floating point number string and computes the
// number of big.Float precision bits required to represent
// that numeric value precisely.
//
// Input Parameters
// ----------------
//
//	floatNumStr string
//	  A floating point number string. Eligible text characters
//	  are:
//	   (1) A minus sign in the first string position (byte).
//	   (2) Any numeric digit character (0-9)
//	   (3) One single decimal point character is optional
//	       (decimal point = period = '.')
//	  Any characters int the string which fail to qualify will
//	  trigger an error.
func calculateBigFloatPrecisionOnFloatNumStr(floatNumStr string, requiredFractionalDigits uint) (bigFloatPrecisionBits uint, err error) {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"floatmechops/AppTest/Test003/bigFloatPow.go", // sourceFile
		"",                                     // parentType (leave "" if none)
		"calculateBigFloatPrecisionOnNumStr()") // thisFunction

	if len(floatNumStr) == 0 {

		return bigFloatPrecisionBits, errGenericCfg.GetErrorInputParameterNilPointer(
			"floatNumStr == ''",
			"",
			"",
			"",
			"'floatNumStr' = Empty String",
			"ERROR: Input parameter 'floatNumStr' is an empty string")
	}

	numStrStatus, err := analyzeFloatNumStr(floatNumStr)

	if err != nil {

		return bigFloatPrecisionBits, errGenericCfg.GetErrorFromFunctionCall(
			"",
			"analyzeFloatNumStr()",
			fmt.Sprintf("analyzeFloatNumStr(floatNumStr) - floatNumStr = %s", floatNumStr),
			"",
			err)
	}

	if numStrStatus.NumFractionalDigits < requiredFractionalDigits {

		deltaFractionalDigits := requiredFractionalDigits - numStrStatus.NumFractionalDigits

		floatNumStr = floatNumStr + strings.Repeat("0", int(deltaFractionalDigits))

		numStrStatus.NumFractionalDigits += deltaFractionalDigits

	} else if numStrStatus.NumFractionalDigits > requiredFractionalDigits {

		numStrStatus.NumFractionalDigits = requiredFractionalDigits
	}

	if numStrStatus.ValueIsZero {

		bigFloatPrecisionBits = 0

		return bigFloatPrecisionBits, nil

	}

	if numStrStatus.ValueIsOne {
		bigFloatPrecisionBits = 8
		return bigFloatPrecisionBits, nil
	}

	totalFinalDigits := numStrStatus.NumIntegerDigits + numStrStatus.NumFractionalDigits

	bigFloatPrecisionBits = calculateBigFloatPrecision(totalFinalDigits)

	return bigFloatPrecisionBits, nil
}

// calculateBigFloatPrecision
//
// Computes the number of big.Float precision bit required to accurately
// represent the numerical value described by the total number of integer
// and fractional digits contained in that numerical value.
//
// Input Parameters
// ----------------
//
//	totalIntAndFracDigits uint
//	  The total number of integer and fractional digits in the numerical
//	  value for which the number big.Float precision bits is being calculated.
//
// Return Parameters
// -----------------
//
//	bigFloatPrecisionBits uint
//	  The number of big.Float precision bits required to accurately represent
//	  a numerical value containing the total integer and fractional digits
//	  specified in input parameter totalIntAndFracDigits.
//
//	  Be advised: If input parameter 'bigFloatPrecisionBits' is equal to zero,
//	  'bigFloatPrecisionBits' will be set equal to zero.
func calculateBigFloatPrecision(totalIntAndFracDigits uint) (bigFloatPrecisionBits uint) {

	if totalIntAndFracDigits == 0 {
		return 0
	}

	baseBits := float64(totalIntAndFracDigits) * 3.32192809488736

	// Add safety margin for exponentiation
	safety := 64.0

	// Round up to nearest whole bit
	bigFloatPrecisionBits = uint(math.Ceil(baseBits + safety))

	return bigFloatPrecisionBits
}

// multiplyFloatNumStrForExponentiation
//
// Receives a floating point number string, converts it to
// a big.Float value and then raises that value to the power
// of a positive, integer exponent.
//
// Input Parameters
// ================
//
// baseFloatNumStr string
//
//	A valid floating point number string. Eligible string
//	characters are:
//	  (1) A minus sign ('-') in the first character position of
//	      the string.
//	  (2) A single decimal point (a.k.a. period ('.') ) located
//	      anywhere in the string.
//	  (3) Numeric characters '0' - '9'
//
//	 If any string characters fail to meet these criteria, an
//	 error will be returned
func multiplyFloatNumStrForExponentiation(
	baseFloatNumStr string,
	raiseToPowerExponent uint,
	requiredFractionalDigits uint) (raiseToPowerResult *big.Float, err error) {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"floatmechops/AppTest/Test003/bigFloatPow.go", // sourceFile
		"",                                       // parentType (leave "" if none)
		"multiplyFloatNumStrForExponentiation()") // thisFunction

	raiseToPowerResult = big.NewFloat(0).SetPrec(0)

	if len(baseFloatNumStr) == 0 {

		return raiseToPowerResult, errGenericCfg.GetErrorInputParameterNilPointer(
			"baseFloatNumStr == ''",
			"",
			"",
			"",
			"'baseFloatNumStr' = Empty String",
			"ERROR: Input parameter 'baseFloatNumStr' is an empty string")
	}

	if raiseToPowerExponent == 0 {

		raiseToPowerResult.SetString("1")

		return raiseToPowerResult, nil
	}

	numStrStatus, err := analyzeFloatNumStr(baseFloatNumStr)

	if err != nil {

		return raiseToPowerResult, errGenericCfg.GetErrorFromFunctionCall(
			"",
			"analyzeFloatNumStr()",
			"",
			"",
			err)
	}

	// If baseFloatNumStr is a negative number,
	// convert it to positive number for calculation
	// purposes.
	if numStrStatus.ValueIsNegative {
		baseFloatNumStr =
			string(numStrStatus.IntegerRunes) + "." + string(numStrStatus.FractionalRunes)
	}

	if numStrStatus.ValueIsZero {

		raiseToPowerResult = big.NewFloat(0).SetPrec(0)

		return raiseToPowerResult, nil
	}

	if numStrStatus.ValueIsOne {

		raiseToPowerResult = big.NewFloat(1)

		if numStrStatus.ValueIsNegative {
			raiseToPowerResult.Neg(raiseToPowerResult)
		}

		return raiseToPowerResult, nil
	}

	if numStrStatus.NumIntegerDigits == 0 {

		numStrStatus.IntegerRunes = []rune("0")

		baseFloatNumStr =
			string(numStrStatus.IntegerRunes) + "." + string(numStrStatus.FractionalRunes)

		numStrStatus.NumIntegerDigits = 1
	}

	if numStrStatus.NumFractionalDigits == 0 {

		numStrStatus.FractionalRunes = []rune("0")

		baseFloatNumStr = string(numStrStatus.IntegerRunes) + ".0"

		numStrStatus.NumFractionalDigits = 1
	}

	baseFloatNumStrPrecisionBits, err := calculateBigFloatPrecisionOnFloatNumStr(baseFloatNumStr, requiredFractionalDigits)

	if err != nil {

		return raiseToPowerResult, errGenericCfg.GetErrorFromFunctionCall(
			"",
			"calculateBigFloatPrecisionOnFloatNumStr(baseFloatNumStr, requiredFractionalDigits)",
			fmt.Sprintf("calculateBigFloatPrecisionOnFloatNumStr(%s)", baseFloatNumStr),
			"",
			err)
	}

	if raiseToPowerExponent == 1 {

		raiseToPowerResult.SetPrec(baseFloatNumStrPrecisionBits).SetString(baseFloatNumStr)

		if numStrStatus.ValueIsNegative {
			raiseToPowerResult.Neg(raiseToPowerResult)
		}

		return raiseToPowerResult, nil
	}

	baseMultiplier, ok := big.NewFloat(0).SetPrec(baseFloatNumStrPrecisionBits).SetString(baseFloatNumStr)

	if ok == false {

		return raiseToPowerResult, errGenericCfg.GetErrorFromFunctionCall(
			"",
			fmt.Sprintf("baseMultiplier,ok=new(big.Float).SetString(%s)", baseFloatNumStr),
			"(big.Float).SetString(baseFloatNumStr) failed on 'baseFloatNumStr'",
			"Error: 'ok' returned 'false'",
			err)

	}

	if numStrStatus.ValueIsNegative {
		baseMultiplier.Neg(baseMultiplier)
	}

	requiredResultFloatPrecisionBits, err := calculateBigFloatPrecisionExponentMultiplication(
		numStrStatus.NumIntegerDigits,
		numStrStatus.NumFractionalDigits,
		requiredFractionalDigits,
		raiseToPowerExponent)

	if err != nil {

		return raiseToPowerResult, errGenericCfg.GetErrorFromFunctionCall(
			"",
			"calculateFloatPrecisionExponentMultiplication()",
			"",
			"",
			err)
	}

	resultMultiplier, ok := new(big.Float).SetPrec(requiredResultFloatPrecisionBits).SetMode(big.ToNearestAway).SetString(baseFloatNumStr)

	if ok == false {

		return raiseToPowerResult, errGenericCfg.GetErrorFromFunctionCall(
			"",
			fmt.Sprintf("resultMultiplier, ok=new(big.Float).SetString(%s)", baseFloatNumStr),
			"(big.Float).SetString(floatNumStr) failed on 'resultMultiplier'",
			"Error: 'ok' returned 'false'",
			err)

	}

	adjustedExponent := raiseToPowerExponent - 1

	for range adjustedExponent {

		resultMultiplier = resultMultiplier.Mul(resultMultiplier, baseMultiplier)

	}

	return resultMultiplier, nil
}

// multiplyForExponentiation
//
// Receives a *big.Float floating point value (baseFloatValue) and
// raises that value to the power specified by input parameter,
// 'raiseToPowerExponent'.
//
// The calculated value of 'baseFloatValue' raised to the power
// of 'raiseToPowerExponent' (baseFloatValue^raiseToPowerExponent)
// is returned as *big.Float result value.
//
// This exponentiation operation will be performed using multiple
// multiplication operations.
//
// Input Parameters
// ----------------
//
//	 baseFloatValue *big.Float
//	   The base input value which will be raised to the power
//			specified by input parameter 'raiseToPowerExponent'.
//
//	 baseFractionalDigits uint
//	   The number of fractional digits to the right of the
//		  decimal place in input parameter 'baseFloatValue'.
//		  This will determine the degree accuracy for the
//		  'baseFloatValue' value in all future calculations.
//
//
//	 raiseToPowerExponent uint
//		  The exponent value. 'baseFloatValue' will be raised to
//		  the power specified by this exponent value.
//
// Return Parameters
// -----------------
//
// raiseToPowerResult *big.Float
//
//	The calculated result of the exponentiation operation.
//	Input parameter 'baseFloatValue' will be raised to the power
//	of 'raiseToPowerExponent' and returned here in
//	'raiseToPowerResult'.
//
// err error
func multiplyForExponentiation(
	baseFloatValue *big.Float,
	baseFractionalDigits uint,
	requiredFractionalDigits uint,
	raiseToPowerExponent uint) (raiseToPowerResult *big.Float, err error) {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"floatmechops/AppTest/Test003/bigFloatPow.go", // sourceFile
		"",                            // parentType (leave "" if none)
		"multiplyForExponentiation()") // thisFunction

	raiseToPowerResult = big.NewFloat(0).SetPrec(0)

	if baseFloatValue == nil {

		return raiseToPowerResult,
			errGenericCfg.GetErrorInputParameterNilPointer(
				"baseFloatValue == nil",
				"",
				"",
				"",
				"Input parameter 'baseFloatValue' is invalid",
				"ERROR: Input parameter 'baseFloatValue' is a nil pointer")
	}

	if raiseToPowerExponent == 0 {
		raiseToPowerResult.SetString("1")
		return raiseToPowerResult, nil
	}

	if raiseToPowerExponent == 1 {
		baseFloatValueStr := baseFloatValue.Text('f', int(requiredFractionalDigits))
		raiseToPowerResult.SetString(baseFloatValueStr)
		return raiseToPowerResult, nil
	}

	baseFloatNumStr := baseFloatValue.Text('f', int(baseFractionalDigits))

	numStrStatus, err := analyzeFloatNumStr(baseFloatNumStr)

	if err != nil {

		return raiseToPowerResult, errGenericCfg.GetErrorFromFunctionCall(
			"",
			"analyzeFloatNumStr()",
			"",
			"",
			err)
	}

	if numStrStatus.ValueIsZero {
		raiseToPowerResult.SetString("0")
		return raiseToPowerResult, nil

	}

	if numStrStatus.ValueIsOne {
		raiseToPowerResult.SetString("1")
		return raiseToPowerResult, nil
	}

	totalBaseDigits := numStrStatus.NumIntegerDigits + numStrStatus.NumFractionalDigits

	baseMultiplierPrecisionBits := calculateBigFloatPrecision(totalBaseDigits)

	requiredResultFloatPrecisionBits, err := calculateBigFloatPrecisionExponentMultiplication(
		numStrStatus.NumIntegerDigits,
		numStrStatus.NumFractionalDigits,
		requiredFractionalDigits,
		raiseToPowerExponent)

	if err != nil {

		return raiseToPowerResult, errGenericCfg.GetErrorFromFunctionCall(
			"",
			"calculateFloatPrecisionExponentMultiplication()",
			"",
			"",
			err)
	}

	baseMultiplier, ok := new(big.Float).SetPrec(baseMultiplierPrecisionBits).SetMode(big.ToNearestAway).SetString(baseFloatNumStr)

	if ok == false {

		return raiseToPowerResult, errGenericCfg.GetErrorFromFunctionCall(
			"",
			fmt.Sprintf("baseMultiplier,ok=new(big.Float).SetString(%s)", baseFloatNumStr),
			"(big.Float).SetString(baseFloatNumStr) failed on 'baseMultiplier'",
			"Error: 'ok' returned 'false'",
			err)

	}

	resultMultiplier, ok := new(big.Float).SetPrec(requiredResultFloatPrecisionBits).SetMode(big.ToNearestAway).SetString(baseFloatNumStr)

	if ok == false {

		return raiseToPowerResult, errGenericCfg.GetErrorFromFunctionCall(
			"",
			fmt.Sprintf("resultMultiplier, ok=new(big.Float).SetString(%s)", baseFloatNumStr),
			"(big.Float).SetString(floatNumStr) failed on 'resultMultiplier'",
			"Error: 'ok' returned 'false'",
			err)

	}

	adjustedExponent := raiseToPowerExponent - 1

	for range adjustedExponent {

		resultMultiplier = resultMultiplier.Mul(resultMultiplier, baseMultiplier)

	}

	return resultMultiplier, nil
}
