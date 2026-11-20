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

		return
	}

	decResult := decimal.PowExact(exponent)

	decResultStr := decResult.String()

	expectedResult := "39617710.92824071"

	if decResultStr != expectedResult {
		t.Errorf("\n"+
			"Expected decimal.PowExact = %s\n"+
			"Instead decimal.PowExact  = %s", expectedResult, decResultStr)
	}

	return
}

// TestDecimal_NewFromFloat64_0010
//
// Tests an alternative technique for calling method
// Decimal.NewFromFloat64().
func TestDecimal_NewFromFloat64_0010(t *testing.T) {

	var v float64
	var targetFractionalDigits uint
	var mode big.RoundingMode
	var exponent uint

	v = 6.9857

	expectedInitialDecValueStr := "6.98570000"

	targetFractionalDigits = 8
	mode = big.ToNearestAway
	exponent = 9

	decType := Decimal{}

	decimal := decType.NewFromFloat64(v, targetFractionalDigits, mode)

	actualInitialDecValueStr := decimal.String()

	if actualInitialDecValueStr != expectedInitialDecValueStr {
		t.Errorf("\n"+
			"Expected Initial Decimal Value: %s\n"+
			"  Actual Initial Decimal Value: %s\n",
			expectedInitialDecValueStr, actualInitialDecValueStr)

		return
	}

	decResult := decimal.PowExact(exponent)

	decResultStr := decResult.String()

	expectedResult := "39617710.92824071"

	if decResultStr != expectedResult {
		t.Errorf("\n"+
			"Expected decimal.PowExact = %s\n"+
			"Instead decimal.PowExact  = %s", expectedResult, decResultStr)

		return
	}

	return
}
