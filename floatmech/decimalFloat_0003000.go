package floatmech

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
)

type decimalFloatF0003000 struct {
}

func (df3 *decimalFloatF0003000) calcBits(intBits uint, reqFracDigits uint, exp uint) uint {
	fracBits := uint(math.Ceil(float64(reqFracDigits) * 3.322))
	guardBits := uint(166 + bits.Len(exp)*2)
	return intBits + fracBits + guardBits
}

func (df3 *decimalFloatF0003000) calcBitsSafe(intBits uint64, reqFracDigits uint, exp uint) (uint, error) {
	fracBits := uint64(math.Ceil(float64(reqFracDigits) * 3.322))
	guardBits := uint64(166 + bits.Len(exp)*2)
	total := intBits + fracBits + guardBits
	if total > uint64(big.MaxPrec) {
		return 0, fmt.Errorf("required precision (%d bits) exceeds big.MaxPrec", total)
	}
	return uint(total), nil
}
