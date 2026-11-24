package main

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
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

	calcValue, err := testPrecision02(floatNumStr, exponentValue, requiredFractionalDigits)

	if err != nil {

		err2 := errGenericCfg.GetErrorFromFunctionCall(
			"",
			"testPrecision02()",
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
		"Function Completed: main\n"+
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

	functionName := "testPrecision01"
	lineBreakStr := strings.Repeat("-", 50)

	fmt.Printf("\n\n%v\n"+
		"Running Function: %s\n"+
		"%v\n\n", lineBreakStr, functionName, lineBreakStr)

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

	if valueIsZero {

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

	requiredFloatPrecision, err := calculateBigFloatPrecisionExponentMultiplication(
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

func testPrecision02(floatNumStr string, exponent uint, requiredFractionalDigits uint) (calcValueStr string, err error) {

	functionName := "testPrecision02"
	lineBreakStr := strings.Repeat("-", 50)

	fmt.Printf("\n\n%v\n"+
		"Running Function: %s\n"+
		"%v\n\n", lineBreakStr, functionName, lineBreakStr)

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

	bigFloatPrecision = calculateBigFloatPrecision(totalFinalDigits)

	return bigFloatPrecision, nil
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

	integerDigits, fractionalDigits, valueIsZero, valueIsOne, err := analyzeFloatNumStr(floatNumStr)

	if err != nil {

		return bigFloatPrecisionBits, errGenericCfg.GetErrorFromFunctionCall(
			"",
			"analyzeFloatNumStr()",
			fmt.Sprintf("analyzeFloatNumStr(floatNumStr) - floatNumStr = %s", floatNumStr),
			"",
			err)
	}

	if fractionalDigits < requiredFractionalDigits {

		deltaFractionalDigits := requiredFractionalDigits - fractionalDigits

		floatNumStr = floatNumStr + strings.Repeat("0", int(deltaFractionalDigits))

		fractionalDigits += deltaFractionalDigits

	} else if fractionalDigits > requiredFractionalDigits {

		fractionalDigits = requiredFractionalDigits
	}

	if valueIsZero {
		bigFloatPrecisionBits = 0
		return bigFloatPrecisionBits, nil

	}

	if valueIsOne {
		bigFloatPrecisionBits = 8
		return bigFloatPrecisionBits, nil
	}

	totalFinalDigits := integerDigits + fractionalDigits

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

	integerDigits, fractionalDigits, valueIsZero, valueIsOne, err := analyzeFloatNumStr(baseFloatNumStr)

	if err != nil {

		return raiseToPowerResult, errGenericCfg.GetErrorFromFunctionCall(
			"",
			"analyzeFloatNumStr()",
			"",
			"",
			err)
	}

	if valueIsZero {
		baseFloatNumStr = "0." + strings.Repeat("0", int(requiredFractionalDigits))
	}

	if valueIsOne {
		baseFloatNumStr = "1." + strings.Repeat("0", int(requiredFractionalDigits))
	}

	if integerDigits == 0 {

		baseFloatNumStr = "0." + strconv.FormatUint(uint64(fractionalDigits), 10)

		integerDigits = 1
	}

	if fractionalDigits == 0 {

		baseFloatNumStr = strconv.FormatUint(uint64(integerDigits), 10)

		baseFloatNumStr += ".0"

		fractionalDigits = 1
	}

	//totalFinalDigits := integerDigits + fractionalDigits

	if raiseToPowerExponent == 1 {
		raiseToPowerResult.SetString(baseFloatNumStr)
		return raiseToPowerResult, nil
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

	baseMultiplier, ok := big.NewFloat(0).SetPrec(baseFloatNumStrPrecisionBits).SetString(baseFloatNumStr)

	if ok == false {

		return raiseToPowerResult, errGenericCfg.GetErrorFromFunctionCall(
			"",
			fmt.Sprintf("baseMultiplier,ok=new(big.Float).SetString(%s)", baseFloatNumStr),
			"(big.Float).SetString(baseFloatNumStr) failed on 'baseFloatNumStr'",
			"Error: 'ok' returned 'false'",
			err)

	}

	requiredResultFloatPrecisionBits, err := calculateBigFloatPrecisionExponentMultiplication(
		integerDigits,
		fractionalDigits,
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

	integerDigits, fractionalDigits, valueIsZero, valueIsOne, err := analyzeFloatNumStr(baseFloatNumStr)

	if err != nil {

		return raiseToPowerResult, errGenericCfg.GetErrorFromFunctionCall(
			"",
			"analyzeFloatNumStr()",
			"",
			"",
			err)
	}

	if valueIsZero {
		raiseToPowerResult.SetString("0")
		return raiseToPowerResult, nil

	}

	if valueIsOne {
		raiseToPowerResult.SetString("1")
		return raiseToPowerResult, nil
	}

	totalBaseDigits := integerDigits + fractionalDigits

	baseMultiplierPrecisionBits := calculateBigFloatPrecision(totalBaseDigits)

	requiredResultFloatPrecisionBits, err := calculateBigFloatPrecisionExponentMultiplication(
		integerDigits,
		fractionalDigits,
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
