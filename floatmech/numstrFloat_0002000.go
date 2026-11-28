package floatmech

import (
	"fmt"
	"strings"

	ePref "github.com/MikeAustin71/errpref/v2"
)

type numStrFloatUtil020 struct {
}

//	 analyzeNumStrFloat subjects a floating point number string
//	 to analysis and validation.
//
//	 If successful, the number of integer numerical digits and
//	 fractional numeric digits are returned in a NumStrFloatDto
//	 object.
//
//	 If the floating point number string is found to be invalid,
//	 an error is returned.
//
//	 Input Parameters
//	 ----------------
//
//	 floatNumStr  string
//	   This parameter must be formatted as a number string containing one or
//		  more numeric digits. An included decimal point ('.') is optional.
//
//		  If any characters other than numeric digits or a single decimal point
//		  ('.') are identified in floatNumStr, an error will be returned.
//
//	 Return Parameters
//	 -----------------
//
//		numStrStatus NumStrFloatDto
//	   If successful, this method returns a NumStrStatusDto object
//	   encapsulating detailed information on the floating point value
//	   passed as input parameter floatNumStr. This information includes
//	   number of integer and fractional digits in floatNumStr.
//
//		err error
func (nStr *numStrFloatUtil020) analyzeNumStrFloat(floatNumStr string, targetFractionalDigits uint) (
	numStrStatus NumStrFloatDto, err error) {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"floatmech/numstrFloat_0002000.go", // sourceFile
		"numStrFloatF0002000",              // parentType (leave "" if none)
		"analyzeNumStrFloat()")             // thisFunction

	if len(floatNumStr) == 0 {

		return numStrStatus,
			errGenericCfg.GetErrorInputParameterOutOfRange(
				"floatNumStr",
				"",
				"",
				"",
				"Input parameter 'floatNumStr' must contain at least one numeric digit",
				"ERROR: Input parameter 'floatNumStr' is an empty string")
	}

	foundDecimal := false

	for i, char := range floatNumStr {

		if i == 0 {

			if char == '-' {

				numStrStatus.NumericSign.SetNegative()

				continue

			} else {

				numStrStatus.NumericSign.SetPositive()
			}

		} else {
			// i != 0
			if char == '-' {
				return numStrStatus,
					errGenericCfg.GetErrorNativeToFunction(
						fmt.Sprintf("floatNumStr= %s", floatNumStr),
						"",
						"",
						"",
						"'floatNumStr' contains a minus sign (-), but NOT at the beginning of the number string.",
						"ERROR: Input parameter 'floatNumStr' contains an invalid minus sign (-).")
			}
		}

		if char == '.' {

			if foundDecimal == true {

				return numStrStatus,
					errGenericCfg.GetErrorNativeToFunction(
						fmt.Sprintf("floatNumStr= %s", floatNumStr),
						"",
						"",
						"",
						"Input parameter 'floatNumStr' contains more than one decimal point ('.')",
						"ERROR: Input parameter 'floatNumStr' contains invalid characters")

			}

			foundDecimal = true

			continue
		}

		if char >= '0' && char <= '9' {

			if foundDecimal {

				numStrStatus.FractionalRunes = append(numStrStatus.FractionalRunes, char)

				numStrStatus.NumFractionalDigits++

			} else {

				// Strip off leading integer zeros
				if char == '0' && len(numStrStatus.IntegerRunes) == 0 {
					continue
				}

				numStrStatus.IntegerRunes = append(numStrStatus.IntegerRunes, char)

				numStrStatus.NumIntegerDigits++
			}

			continue
		}

		return NumStrFloatDto{},
			errGenericCfg.GetErrorNativeToFunction(
				fmt.Sprintf("floatNumStr= %s", floatNumStr),
				"",
				"",
				"",
				"Input parameter 'floatNumStr' contains non-numeric characters",
				"ERROR: Input parameter 'floatNumStr' contains invalid characters")
	}

	if len(numStrStatus.IntegerRunes) == 0 &&
		len(numStrStatus.FractionalRunes) == 0 {

		return NumStrFloatDto{},
			errGenericCfg.GetErrorNativeToFunction(
				fmt.Sprintf("floatNumStr= %s", floatNumStr),
				"",
				"",
				"",
				"No valid numerical digits were found in 'floatNumStr'",
				"ERROR: Input parameter 'floatNumStr' contains invalid characters")
	}

	if len(numStrStatus.IntegerRunes) == 0 && len(numStrStatus.FractionalRunes) > 0 {

		numStrStatus.IntegerRunes = []rune("0")

		numStrStatus.NumIntegerDigits = 1
	}

	if len(numStrStatus.IntegerRunes) > 0 && len(numStrStatus.FractionalRunes) == 0 {

		numStrStatus.FractionalRunes = []rune(strings.Repeat("0", int(targetFractionalDigits)))

		numStrStatus.NumFractionalDigits = targetFractionalDigits

	}

	if string(numStrStatus.IntegerRunes) == "0" && string(numStrStatus.FractionalRunes) == "0" {
		numStrStatus.ValueIsZero = true
		numStrStatus.NumericSign.SetZero()
	}

	if string(numStrStatus.IntegerRunes) == "1" && string(numStrStatus.FractionalRunes) == "0" {
		numStrStatus.ValueIsOne = true
	}

	return numStrStatus, nil
}
