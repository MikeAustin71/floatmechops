package floatmech

import (
	"fmt"
	"math/big"
)

type bigFloatHelper0020 struct{}

func (bigFloatHelp0020 *bigFloatHelper0020) getRoundingModeToStr(
	mode big.RoundingMode) string {

	switch mode {
	case big.ToNearestAway:
		return "ToNearestAway"
	case big.ToNearestEven:
		return "ToNearestEven"
	case big.AwayFromZero:
		return "AwayFromZero"
	case big.ToZero:
		return "ToZero"
	case big.ToNegativeInf:
		return "ToNegativeInf"
	case big.ToPositiveInf:
		return "ToPositiveInf"
	default:
		return fmt.Sprintf("Invalid(%d)", mode)
	}

}
