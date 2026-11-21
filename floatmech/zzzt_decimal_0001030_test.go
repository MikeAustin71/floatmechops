package floatmech

import (
	"math/big"
	"testing"
)

func TestDecimal_PowExact_00010(t *testing.T) {

	var vNumStr string
	var targetFractionalDigits uint
	var mode big.RoundingMode
	var exponent uint

	vNumStr = "7.8715"

	expectedInitialDecValueStr := "7.8715000000"

	targetFractionalDigits = 10
	mode = big.ToNearestAway
	exponent = 6

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

	expectedResult := "237872.9172319539"

	if decResultStr != expectedResult {
		t.Errorf("\n"+
			"Expected decimal.PowExact = %s\n"+
			"Instead decimal.PowExact  = %s", expectedResult, decResultStr)

		return
	}

	return
}

func TestDecimal_PowExact_00020(t *testing.T) {

	var vNumStr string
	var targetFractionalDigits uint
	var mode big.RoundingMode
	var exponent uint

	vNumStr = "10.81233"

	expectedInitialDecValueStr := "10.81233000000000"

	targetFractionalDigits = 14
	mode = big.ToNearestAway
	exponent = 8

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

	expectedResult := "186790307.75412594147937"

	if decResultStr != expectedResult {
		t.Errorf("\n"+
			"Expected decimal.PowExact = %s\n"+
			"Instead decimal.PowExact  = %s", expectedResult, decResultStr)

		return
	}

	return
}

func TestDecimal_PowExact_00030(t *testing.T) {

	var vNumStr string
	var targetFractionalDigits uint
	var mode big.RoundingMode
	var exponent uint

	vNumStr = "15.81233"

	expectedInitialDecValueStr := "15.81233000000000"

	targetFractionalDigits = 14
	mode = big.ToNearestAway
	exponent = 4

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

	expectedResult := "62514.89090130012128"

	if decResultStr != expectedResult {
		t.Errorf("\n"+
			"Expected decimal.PowExact = %s\n"+
			"Instead decimal.PowExact  = %s", expectedResult, decResultStr)

		return
	}

	return
}

func TestDecimal_PowExact_00040(t *testing.T) {

	var vNumStr string
	var targetFractionalDigits uint
	var mode big.RoundingMode
	var exponent uint

	vNumStr = "1.0691372836118"

	expectedInitialDecValueStr := "1.06913728361180"

	targetFractionalDigits = 14
	mode = big.ToNearestAway
	exponent = 24

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

	expectedResult := "4.97511810598396"

	if decResultStr != expectedResult {
		t.Errorf("\n"+
			"Expected decimal.PowExact = %s\n"+
			"Instead decimal.PowExact  = %s", expectedResult, decResultStr)

		return
	}

	return
}
