package floatmech

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"

	ePref "github.com/MikeAustin71/errpref/v2"
)

type decimalFloatUtil030 struct {
}

func (df3 *decimalFloatUtil030) calcBits(
	intBits uint, reqFracDigits uint, exp uint) uint {

	fracBits := uint(math.Ceil(float64(reqFracDigits) * 3.32192809488736))

	guardBits := uint(166 + bits.Len(exp)*2)

	return intBits + fracBits + guardBits
}

func (df3 *decimalFloatUtil030) calcBitsSafe(
	intBits uint64, reqFracDigits uint, exp uint) (uint, error) {

	fracBits := uint64(math.Ceil(float64(reqFracDigits) * 3.32192809488736))

	guardBits := uint64(166 + bits.Len(exp)*2)

	total := intBits + fracBits + guardBits

	if total > uint64(big.MaxPrec) {
		return 0, fmt.Errorf("required precision (%d bits) exceeds big.MaxPrec", total)
	}

	return uint(total), nil
}

//		calcPrecisionSafe calculates the number of precision bits
//		required for a big.Float number to accurately represent
//	 the total number of digits specified by input parameters
//	 'intDigits', 'reqFracDigits' and 'safetyDigits'.
//
//		Input Parameters
//	 ----------------
//
//		intDigits uint
//		The number of integer digits in the floating point number.
//
//		reqFracDigits uint
//		The number of fractional digits to the right of the decimal
//		point in the floating point number.
//
//		safetyDigits	uint
//		The number of extra digit positions added for safety
//
//		Return Parameters
//		-----------------
//
//		uint
//		The total number of precision bits required to accurately
//		represent the subject floating point number.
//
//		error
func (df3 *decimalFloatUtil030) calcPrecisionSafe(
	intDigits uint, reqFracDigits uint, safetyDigits uint) (uint, error) {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"decimalFloat_0003000.go", // sourceFile
		"decimalFloatUtil030",     // parentType (leave "" if none)
		"calcPrecisionSafe()")     // thisFunction

	totalProjectedDigits := uint64(intDigits + reqFracDigits + safetyDigits)

	totalBaseBits := uint64(math.Ceil(float64(totalProjectedDigits) * 3.32192809488736))

	totalFinalBits := totalBaseBits

	if totalFinalBits > uint64(big.MaxPrec) {

		return 0, errGenericCfg.GetErrorNativeToFunction(
			fmt.Sprintf("intDigits= '%d' reqFracDigits= '%d'", intDigits, reqFracDigits),
			"",
			"",
			"",
			fmt.Sprintf("Calculated Required Precision Bits= '%d'", totalFinalBits),
			fmt.Sprintf("Total Precision Bits '%d' exceeds Maximum Precision Bits '%d'",
				totalFinalBits, big.MaxPrec))
	}

	return uint(totalProjectedDigits), nil
}

func (df3 *decimalFloatUtil030) powFloat(base *big.Float, exp uint, prec uint) *big.Float {

	result := new(big.Float).SetPrec(prec).SetFloat64(1)

	b := new(big.Float).SetPrec(prec).Set(base)

	for exp > 0 {

		if exp&1 == 1 {

			result.Mul(result, b)
		}

		exp >>= 1

		if exp > 0 {

			b.Mul(b, b)
		}
	} // End of for loop

	return result
}

func (df3 *decimalFloatUtil030) powRat(
	base *big.Rat, exp uint) *big.Rat {

	result := big.NewRat(1, 1)

	b := new(big.Rat).Set(base)

	for exp > 0 {

		if exp&1 == 1 {

			result.Mul(result, b)
		}

		exp >>= 1

		if exp > 0 {

			b.Mul(b, b)

		}

	}

	return result
}
