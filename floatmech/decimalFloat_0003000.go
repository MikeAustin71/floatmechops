package floatmech

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
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
