package main

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

// This test duplicates:
//
//	TestDecimal_PowExact_00020
//	zzzt_decimal_0001030_test.go
//
//	The difference - TestPrec002() computes the correct result.
func main() {
	TestPrec002()
}

// TestPrec002
//
// This Passes
func TestPrec002() {

	lineBreakStr := strings.Repeat("-", 50)

	fmt.Printf("\n%v\nRunning TestPrec002()\nFile: /floatmechops/AppTest/Test002/float_pow.go\n%v\n\n",
		lineBreakStr,
		lineBreakStr)

	//                    Decimal            1         2         3
	//                    Digits    12345678901234567890123456789012345
	expectedResultStr := "186790307.75412594147937"

	requiredFractionalDigits := uint(14)

	var requiredTotalDigits uint

	// Calculate total number of integer and
	// fractional digits.
	for _, s := range expectedResultStr {

		if s == '.' {
			continue
		}

		if s >= '0' && s <= '9' {
			requiredTotalDigits++
		}

	}

	exponentValue := uint(8)

	baseStr := "10.81233"

	requiredFloatPrecision := ComputeBigFloatPrecisionBits(requiredTotalDigits, exponentValue)

	baseFloat, isOk :=
		new(big.Float).
			SetPrec(requiredFloatPrecision).
			SetMode(big.AwayFromZero).
			SetString(baseStr)

	if !isOk {
		fmt.Printf("baseFloat SetString(\"3.14159\") Failed!")
		return
	}

	resultBaseToExponent, isOk :=
		new(big.Float).
			SetPrec(requiredFloatPrecision).
			SetMode(big.AwayFromZero).
			SetString(baseStr)

	if !isOk {
		fmt.Printf("resultBaseToExponent SetString(\"3.14159\") Failed!")
		return
	}

	adjustedExponentValue := exponentValue - 1

	for range adjustedExponentValue {

		resultBaseToExponent.Mul(resultBaseToExponent, baseFloat)
	}

	/*
	   for i := uint(0); i < adjustedExponentValue; i++ {

	     resultBaseToExponent.Mul(resultBaseToExponent, baseFloat)
	   }
	*/
	actualResultStr := resultBaseToExponent.Text('f', int(requiredFractionalDigits))

	fmt.Printf("Base: %v\n",
		baseStr)

	fmt.Printf("Exponent: %v\n",
		exponentValue)

	fmt.Printf("Required FractionalDigits: %v\n", requiredFractionalDigits)

	fmt.Printf("Total Integer and Fractional Digits\n"+
		"in the Result: %v\n", requiredTotalDigits)

	fmt.Printf("Expected Result: %v\n", expectedResultStr)

	fmt.Printf("Actual Result:   %v\n\n%v\n\n",
		actualResultStr,
		lineBreakStr)

}

// ComputeBigFloatPrecisionBits
//
//	This function computes the number of bits required to represent a
//	floating-point number with the specified number of decimal digits.
//
//	The function uses the following formula to compute the number of bits:
//
//	baseBits = decDigits * log2(10) + safety
//
//	where:
//
//	decDigits - The number of decimal digits to be represented.
//
//	log2(10) - The base-2 logarithm of 10.
//
//	safety - A safety margin for exponentiation.
//
//	Important:
//
//	Remember, You must compute the total decimal digits (decDigits) in
//	the floating point number to be supported with the computed and
//	returned precision value.
//
//	Input Parameters:
//
//	resultDecDigits - The number of decimal digits to be represented.
//
//	multiplyCount - The number of times the base number multiplies by
//	itself.
func ComputeBigFloatPrecisionBits(resultDecDigits uint, multiplyCount uint) uint {

	// Convert decimal digits → bits
	// baseBits := float64(decDigits) * math.Log2(10)
	//
	// math.Log2(10) =
	//  3.321 928 094 887 362 347 870 319 429 489 390 175 864 831 393 024 6
	// float64 maximum 15-digits of precision
	baseBits := float64(resultDecDigits) * 3.321928094887362

	// Add safety margin for exponentiation
	safety := 64.0 + (float64(multiplyCount) * 4.0)

	// Round up to nearest whole bit
	return uint(math.Ceil(baseBits + safety))
}
