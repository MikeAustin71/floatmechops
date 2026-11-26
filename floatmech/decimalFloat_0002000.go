package floatmech

import (
	"fmt"
	"strconv"
	"strings"

	ePref "github.com/MikeAustin71/errpref/v2"
)

type decimalFloatF0002000 struct {
}

// validateBaseFormat verifies that s conforms to the strict base-10 decimal or
// scientific notation grammar required by [PowDecimal].
//
// Enforced Rules:
//   - s must be non-empty and free of leading or trailing whitespace.
//   - Rational fraction slashes ('/') and digit underscores ('_') are rejected so
//     [big.Rat] and [big.Float] tiers accept identical input syntax.
//   - A leading sign ('+' or '-') on the mantissa may appear only at index 0.
//   - If 'e' or 'E' is present, the exponent must be a valid base-10 integer within
//     [-maxSciExp, maxSciExp] to prevent excessive memory allocation during
//     [big.Rat.SetString] parsing.
func (df2 *decimalFloatF0002000) validateBaseFormat(s string) error {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"floatmech/decimalFloat_0002000.go", // sourceFile
		"decimalFloatF0002000",              // parentType (leave "" if none)
		"validateBaseFormat()")              // thisFunction

	if len(s) == 0 {

		return errGenericCfg.GetErrorNativeToFunction(
			"Input parameter 's' is an empty string",
			"",
			"",
			"",
			fmt.Sprintf("s = '%s'", s),
			"Input Parameter 's' number string cannot be empty")
	}

	if strings.TrimSpace(s) != s {

		return errGenericCfg.GetErrorNativeToFunction(
			"'s' has leading or trailing whitespace",
			"",
			"",
			"",
			fmt.Sprintf(" s = '%q'", s),
			"Input Parameter 's' number string cannot have leading or trailing whitespace.")
	}

	if strings.ContainsAny(s, "/_") {

		return errGenericCfg.GetErrorNativeToFunction(
			"",
			"",
			"",
			"",
			fmt.Sprintf("Input Parameter s = '%q'", s),
			"Input Parameter 's' has invalid characters in number string.")
	}

	// Split at 'e' or 'E' if scientific notation is present
	mantissa := s
	if idx := strings.IndexAny(s, "eE"); idx != -1 {
		mantissa = s[:idx]
		exponent := s[idx+1:]

		if len(mantissa) == 0 || len(exponent) == 0 {

			return errGenericCfg.GetErrorNativeToFunction(
				"if len(mantissa) == 0 || len(exponent) == 0",
				"",
				"",
				"",
				fmt.Sprintf("Input Parameter s = '%q'", s),
				"Input Parameter 's' is a malformed scientific notation.")
		}

		// In the exponent, '-' is only allowed at index 0 of the exponent part
		if expMinus := strings.LastIndexByte(exponent, '-'); expMinus > 0 {

			return errGenericCfg.GetErrorNativeToFunction(
				"expMinus := strings.LastIndexByte(exponent, '-'); expMinus > 0",
				"",
				"",
				"",
				fmt.Sprintf("Input Parameter s = '%q' is malformed scientific notation", s),
				"Scientific notation exponent minus sign must immediately follow 'e'/'E'.")
		}

		// Parse the exponent as an int64 to catch both malformed digits and huge magnitudes
		sciExp, err := strconv.ParseInt(exponent, 10, 64)
		if err != nil {

			return errGenericCfg.GetErrorFromFunctionCall(
				"",
				"strconv.ParseInt(exponent, 10, 64)",
				fmt.Sprintf("Input Parameter s = '%q'", s),
				fmt.Sprintf("invalid or out-of-range scientific notation exponent in %q", s),
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

	// In the mantissa (base), '-' is only allowed at index 0
	if idx := strings.LastIndexByte(mantissa, '-'); idx > 0 {

		return errGenericCfg.GetErrorNativeToFunction(
			"if idx := strings.LastIndexByte(mantissa, '-'); idx > 0",
			"",
			"",
			"",
			"",
			fmt.Sprintf("minus sign must be in the first string position: %q", s))
	}

	return nil
}
