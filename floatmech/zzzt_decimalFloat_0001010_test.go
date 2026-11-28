package floatmech

import (
	"math/big"
	"testing"
)

func TestDecimalFloat_NewFromFloatNumStr_00010(t *testing.T) {

	var floatNumStr string
	var targetFractionalDigits uint
	var mode big.RoundingMode

	floatNumStr = "62.7957"

	expectedDecValueStr := "62.79570000"

	targetFractionalDigits = 8
	mode = big.ToNearestAway

	df := &DecimalFloat{}

	decimal, err := df.NewFromFloatNumStr(floatNumStr, targetFractionalDigits, mode)

	if err != nil {
		t.Errorf("NewFromFloatNumStr() failed: %s", err.Error())

		return
	}

	if decimal == nil {
		t.Errorf("NewFromFloatNumStr() returned nil Decimal")

		return
	}

	if !decimal.IsValid() {
		t.Errorf("NewFromFloatNumStr() returned an invalid Decimal (val is nil)")

		return
	}

	actualDecValueStr := decimal.String()

	if actualDecValueStr != expectedDecValueStr {
		t.Errorf("\n"+
			"Expected Decimal Value: %s\n"+
			"  Actual Decimal Value: %s\n",
			expectedDecValueStr, actualDecValueStr)

		return
	}

	return
}
