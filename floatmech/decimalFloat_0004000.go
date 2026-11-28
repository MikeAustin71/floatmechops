package floatmech

import (
	"fmt"
	"strconv"
	"strings"

	ePref "github.com/MikeAustin71/errpref/v2"
)

type decimalFloatUtil040 struct {
}

func (df4 *decimalFloatUtil040) validateFloatNumStr(
	floatNumStr string) error {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"decimalFloat_0004000.go",
		"decimalFloatUtil040",
		"validateFloatNumStr()")

	if len(floatNumStr) == 0 {

		return errGenericCfg.GetErrorNativeToFunction(
			"Input parameter 'floatNumStr' is an empty string",
			"",
			"",
			"",
			fmt.Sprintf("s = '%s'", floatNumStr),
			"Input Parameter 'floatNumStr' number string cannot be empty")
	}

	if strings.TrimSpace(floatNumStr) != floatNumStr {

		return errGenericCfg.GetErrorNativeToFunction(
			"'floatNumStr' has leading or trailing whitespace",
			"",
			"",
			"",
			fmt.Sprintf(" floatNumStr = '%q'", floatNumStr),
			"Input Parameter 'floatNumStr' number string cannot have leading or trailing whitespace.")
	}

	// Start checking for scientific notation
	// Split at 'e' or 'E' if scientific notation is present
	mantissa := floatNumStr

	// Test for scientific notation.
	// If scientific notation discovered, apply additional
	// validation tests. Otherwise, skip scientific notation
	// test and proceed to check for minus sign.
	if idx := strings.IndexAny(floatNumStr, "eE"); idx != -1 {
		mantissa = floatNumStr[:idx]
		exponent := floatNumStr[idx+1:]

		if len(mantissa) == 0 || len(exponent) == 0 {

			return errGenericCfg.GetErrorNativeToFunction(
				"if len(mantissa) == 0 OR len(exponent) == 0",
				"",
				"",
				"",
				fmt.Sprintf("Input Parameter floatNumStr = '%q'", mantissa),
				"Input Parameter 'floatNumStr' is a malformed scientific notation.")
		}

		// In the exponent, '-' is only allowed at index 0 of the exponent part
		if expMinus := strings.LastIndexByte(exponent, '-'); expMinus > 0 {

			return errGenericCfg.GetErrorNativeToFunction(
				"expMinus := strings.LastIndexByte(exponent, '-'); expMinus > 0",
				"",
				"",
				"",
				fmt.Sprintf("Input Parameter floatNumStr = '%q' is malformed scientific notation", floatNumStr),
				"Scientific notation exponent minus sign must immediately follow 'e'/'E'.")
		}

		// Parse the exponent as an int64 to catch both malformed digits and huge magnitudes
		sciExp, err := strconv.ParseInt(exponent, 10, 64)

		if err != nil {

			return errGenericCfg.GetErrorFromFunctionCall(
				"",
				"strconv.ParseInt(exponent, 10, 64)",
				fmt.Sprintf("Input Parameter floatNumStr = '%q'", floatNumStr),
				fmt.Sprintf("invalid or out-of-range scientific notation exponent in %q", floatNumStr),
				err)
		}

		if sciExp > maxSciExp || sciExp < -maxSciExp {

			return errGenericCfg.GetErrorNativeToFunction(
				"if sciExp > maxSciExp || sciExp < -maxSciExp",
				"",
				"",
				"",
				"",
				fmt.Sprintf("scientific notation exponent %d exceeds allowed range [-%d, %d]", sciExp, maxSciExp, maxSciExp))
		}

	}

	// For scientific notation, the mantissa (base), '-' is only
	// allowed at index 0.
	//
	// Likewise for integer of floating point number strings, '-' is
	// only allowed at index 0.
	if idx := strings.LastIndexByte(mantissa, '-'); idx > 0 {

		return errGenericCfg.GetErrorNativeToFunction(
			"if idx := strings.LastIndexByte(floatNumStr, '-'); idx > 0",
			"",
			"",
			"",
			"",
			fmt.Sprintf("minus sign must be in the first string position: %q", floatNumStr))
	}

	return nil
}
