package floatmech

import (
	"math"
)

type bigFloatHelper0020 struct{}

// computeBigFloatPrecisionBits
//
//		This function computes the number of bits required to represent a
//		floating-point number with the specified number of decimal digits.
//
//		The function uses the following formula to compute the number of bits:
//
//		baseBits = resultDecDigits * log2(10) + safety
//
//		where:
//
//		resultDecDigits - The number of decimal digits to be represented.
//
//		log2(10) - The base-2 logarithm of 10.
//
//		safety - A safety margin for exponentiation.
//
//		Important:
//
//		Remember, You must compute the total decimal digits (resultDecDigits) in
//		the floating point result number to be supported with the computed
//		and returned precision value. 'resultDecDigits' is passed to this function
//	 as an input parameter.
//
//		Input Parameters:
//
//		resultDecDigits - The number of decimal digits to be represented.
//
//		multiplyCount - The number of times the base number multiplies itself. In
//		the exponent operation, X^n, 'n' is the exponent or the number of times
//		the base number multiplies itself.
func (bigFloatHelp0020 *bigFloatHelper0020) computeBigFloatPrecisionBits(
	resultDecDigits uint, multiplyCount uint) uint {

	// Convert decimal digits → bits
	// baseBits := float64(decDigits) * math.Log2(10)
	//
	// math.Log2(10) =
	//  3.321 928 094 887 362 347 870 319 429 489 390 175 864 831 393 024 6
	// float64 max 15-digits of precision
	baseBits := float64(resultDecDigits) * 3.321928094887362

	// Add safety margin for exponentiation
	safety := 64.0 + (float64(multiplyCount) * 4.0)

	// Round up to nearest whole bit
	return uint(math.Ceil(baseBits + safety))
}
