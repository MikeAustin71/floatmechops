package floatmech

import (
	"math/big"
	"testing"
)

func TestDecimal_Pow_0010(t *testing.T) {

	var v float64
	var targetFractionalDigits uint
	var mode big.RoundingMode
	var exponent uint

	v = 6.9857

	expectedInitialDecValueStr := "6.98570000"

	targetFractionalDigits = 8
	mode = big.ToNearestAway
	exponent = 9

	decimal := new(Decimal).NewFromFloat64(v, targetFractionalDigits, mode)

	actualInitalDecValueStr := decimal.String()

	if actualInitalDecValueStr != expectedInitialDecValueStr {
		t.Errorf("\n"+
			"Expected Initial Decimal Value: %s\n"+
			"  Actual Initial Decimal Value: %s\n",
			expectedInitialDecValueStr, actualInitalDecValueStr)
	}

	decResult := decimal.Pow(exponent)

	decResultStr := decResult.String()

	expectedResult := "39617710.92824073"

	if decResultStr != expectedResult {
		t.Errorf("\n"+
			"Expected decimal.PowExact = %s\n"+
			"Instead decimal.PowExact  = %s", expectedResult, decResultStr)
	}

	return
}
