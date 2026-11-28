package floatmech

import (
	"math/big"
	"sync"
)

type bigFloatHelper struct {
	mu sync.Mutex
}

// GetRoundingModeStr receives a big.RoundMode value and converts
// it to a string.
//
// big.RoundMode values are used to control numerical rounding in
// big.Float floating point values.
//
// Input Parameters
// ----------------
//
// mode big.RoundingMode
//
// Return Parameters
// ----------------
//
// string
// The string representation of the big.RoundingMode value.
func (bigFloatHelpr *bigFloatHelper) GetRoundingModeStr(mode big.RoundingMode) string {

	bigFloatHelpr.mu.Lock()

	defer bigFloatHelpr.mu.Unlock()

	return new(bigFloatHelper0020).getRoundingModeToStr(mode)
}
