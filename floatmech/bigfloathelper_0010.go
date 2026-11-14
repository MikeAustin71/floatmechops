package floatmech

import "sync"

type bigFloatHelper struct {
	mu sync.Mutex
}

// ComputeBigFloatPrecisionBits
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
func (bigFloatHelpr *bigFloatHelper) ComputeBigFloatPrecisionBits(
	resultDecDigits uint, multiplyCount uint) uint {

	bigFloatHelpr.mu.Lock()

	defer bigFloatHelpr.mu.Unlock()

	return new(bigFloatHelper0020).
		computeBigFloatPrecisionBits(resultDecDigits, multiplyCount)
}
