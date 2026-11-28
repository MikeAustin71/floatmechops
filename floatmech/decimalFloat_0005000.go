package floatmech

import (
	"fmt"
	"math/big"

	ePref "github.com/MikeAustin71/errpref/v2"
)

type decimalFloatUtil050 struct{}

// getDecFloatFromNumStrDto constructs and returns a
// new, fully populated, DecimalFloat object.
func (df5 *decimalFloatUtil050) getDecFloatFromNumStrDto(
	numStrFloatStatus NumStrFloatDto,
	targetFractionalDigits uint,
	mode big.RoundingMode) (*DecimalFloat, error) {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"decimalFloat_0005000.go",    // sourceFile
		"decimalFloatUtil050",        // parentType (leave "" if none)
		"getDecFloatFromNumStrDto()") // thisFunction

	if !numStrFloatStatus.NumType.IsFloatingPoint() {

		return &DecimalFloat{}, errGenericCfg.GetErrorNativeToFunction(
			"Input numStrFloatStatus.NumType.IsFloatingPoint() == false",
			"",
			"",
			"",
			"numStrFloatStatus is NOT a floating point number.",
			"Error: Input Parameter 'numStrFloatStatus' MUST BE a floating point number.")
	}

	numStrFloatStatus.NumFractionalDigits = max(targetFractionalDigits, numStrFloatStatus.NumFractionalDigits)

	totalPrecisionBits, err := new(decimalFloatUtil030).calcPrecisionSafe(
		numStrFloatStatus.NumIntegerDigits, numStrFloatStatus.NumFractionalDigits, uint(5))

	if err != nil {
		return &DecimalFloat{}, errGenericCfg.GetErrorFromFunctionCall(
			"decimalFloatUtil030",
			"calcPrecisionSafe",
			"",
			"",
			err)
	}

	numFloatStr := string(numStrFloatStatus.IntegerRunes) + "." + string(numStrFloatStatus.FractionalRunes)

	decFloat := new(DecimalFloat)
	var ok bool

	decFloat.val, ok = new(big.Float).SetPrec(totalPrecisionBits).SetString(numFloatStr)

	if !ok {
		return &DecimalFloat{}, errGenericCfg.GetErrorNativeToFunction(
			"bigDec, ok := new(big.Float).SetPrec(totalPrecisionBits).SetString(numFloatStr)",
			"",
			"",
			"",
			fmt.Sprintf("totalPrecisionBits= '%d' numFloatStr= '%s'", totalPrecisionBits, numFloatStr),
			"Error: ok=false; SetString() Failed!")
	}

	decFloat.targetFractionalDigits = targetFractionalDigits

	// Returns Big Float Rounding Mode
	decFloat.mode = decFloat.val.Mode()

	if decFloat.mode != mode {
		bFloatHlpr := new(bigFloatHelper0020)

		decFloatMode := bFloatHlpr.getRoundingModeToStr(decFloat.mode)
		inputMode := bFloatHlpr.getRoundingModeToStr(mode)

		return &DecimalFloat{}, errGenericCfg.GetErrorNativeToFunction(
			"if decFloat.mode != Input mode",
			"",
			"",
			"",
			fmt.Sprintf("decFloat.mode = '%s' Input mode = '%s'", decFloatMode, inputMode),
			"Error: Decimal Float rounding mode conversion failed!")
	}

	// Set negative if required
	if numStrFloatStatus.NumericSign.IsNegative() {
		decFloat.val.Neg(decFloat.val)
		numFloatStr = "-" + numFloatStr
	}

	// Returns Big Float Mantissa Precision
	decFloat.prec = decFloat.val.Prec()

	if decFloat.prec != totalPrecisionBits {

		return &DecimalFloat{}, errGenericCfg.GetErrorNativeToFunction(
			"if decFloat.prec != totalPrecisionBits",
			"",
			"",
			"",
			fmt.Sprintf("decFloat.prec = '%d' Calculated Precision Bits = '%d'", decFloat.prec, totalPrecisionBits),
			"Error: Decimal Float precision conversion failed!")
	}

	actualFloatStr := decFloat.val.Text('f', int(targetFractionalDigits))

	if actualFloatStr != numFloatStr {

		return &DecimalFloat{}, errGenericCfg.GetErrorNativeToFunction(
			"if actualFloatStr != numFloatStr",
			"",
			"",
			"",
			fmt.Sprintf("actualFloatStr = '%s' Input FloatStr = '%s'", actualFloatStr, numFloatStr),
			"Error: Decimal Float value to string conversion failed!")

	}

	return decFloat, nil
}
