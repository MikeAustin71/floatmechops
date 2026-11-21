package floatmech

import (
	"math/big"
	"testing"
)

func TestDecimal_Pow_0010(t *testing.T) {

	var vNumStr string
	var targetFractionalDigits uint
	var mode big.RoundingMode
	var exponent uint

	vNumStr = "6.9857"

	expectedInitialDecValueStr := "6.98570000"

	targetFractionalDigits = 8
	mode = big.ToNearestAway
	exponent = 9

	decimal, err := new(Decimal).NewFromFloatNumStr(vNumStr, targetFractionalDigits, mode)

	if err != nil {
		t.Errorf("NewFromFloatNumStr() failed: %s", err.Error())

		return
	}

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

	expectedResult := "39617710.92824073"

	if decResultStr != expectedResult {
		t.Errorf("\n"+
			"Expected decimal.PowExact = %s\n"+
			"Instead decimal.PowExact  = %s", expectedResult, decResultStr)

		return
	}

	return
}

// TestDecimal_NewFromFloatNumStr_0010
//
// Tests an alternative technique for calling method
// Decimal.NewFromFloatNumStr().
func TestDecimal_NewFromFloatNumStr_0010(t *testing.T) {

	var vNumStr string
	var targetFractionalDigits uint
	var mode big.RoundingMode
	var exponent uint

	vNumStr = "6.9857"

	expectedInitialDecValueStr := "6.98570000"

	targetFractionalDigits = 8
	mode = big.ToNearestAway
	exponent = 9

	decType := Decimal{}

	decimal, err := decType.NewFromFloatNumStr(vNumStr, targetFractionalDigits, mode)

	if err != nil {
		t.Errorf("NewFromFloatNumStr() failed: %s", err.Error())

		return
	}

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

	expectedResult := "39617710.92824073"

	if decResultStr != expectedResult {
		t.Errorf("\n"+
			"Expected decimal.PowExact = %s\n"+
			"Instead decimal.PowExact  = %s", expectedResult, decResultStr)

		return
	}

	return
}
