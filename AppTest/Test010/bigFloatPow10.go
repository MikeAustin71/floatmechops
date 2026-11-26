package main

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"strconv"
	"strings"

	ePref "github.com/MikeAustin71/errpref/v2"
)

func main() {

	lineBreakStr := strings.Repeat("-", 50)

	fileName := "floatmechops/AppTest/Test003/bigFloatPow.go"
	funcName := "main()"
	targetFuncName := "testRunner01()"

	fmt.Printf("\n%s\n"+
		"File Name: %s\n"+
		" Function: %s\n"+
		"%s\n"+
		"Starting Run\n"+
		"Running Function: %s\n"+
		"%s\n", lineBreakStr, fileName, funcName, lineBreakStr, targetFuncName, lineBreakStr)

	testRunner01(fileName)

	fmt.Printf("%s\n"+
		"File Name: %s\n"+
		" Function: %s\n"+
		"%s\n"+
		"Function Completed!\n"+
		"%s\n\n", lineBreakStr, fileName, funcName, lineBreakStr, lineBreakStr)
}

func testRunner01(homeFileName string) {

	lineBreakStr := strings.Repeat("-", 50)
	funcName := "TestRunner02()"
	targetFuncName := "PowDecimal()"

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		homeFileName, // sourceFile
		"",           // parentType (leave "" if none)
		funcName)     // thisFunction

	floatNumStr := "10.81233"

	exponentValue := uint(8)

	requiredFractionalDigits := uint(14)

	//                    Decimal            1         2         3
	//                    Digits    12345678901234567890123456789012345
	expectedResultStr := "186790307.75412594147937"

	fmt.Printf("\n%v\n"+
		"Function: %v\n"+
		"File: %v\n"+
		"%v\n"+
		"floatNumStr: %v\n"+
		"exponentValue: %v\n"+
		"requiredFractionalDigits: %v\n"+
		"%v\n"+
		"Running Function: %v\n"+
		"Starting Calculation Run\n"+
		"%v\n",
		lineBreakStr,
		funcName,
		homeFileName,
		lineBreakStr,
		floatNumStr,
		exponentValue,
		requiredFractionalDigits,
		lineBreakStr,
		targetFuncName,
		lineBreakStr)

	calcValueStr, calcFloatValue, err := PowDecimal(floatNumStr, exponentValue, requiredFractionalDigits)

	if err != nil {

		err2 := errGenericCfg.GetErrorFromFunctionCall(
			"",
			"PowDecimal()",
			"",
			"",
			err)

		fmt.Printf("Function: %s\n"+
			"Status: Error\n"+
			"%v", funcName, err2)

		return
	}

	statusStr := "SUCCESS!"

	if calcValueStr != expectedResultStr {
		statusStr = "FAILURE!"
	}

	fmt.Printf("%v\n"+
		"Function Completed: %v\n"+
		"File: %v\n"+
		"%v\n"+
		"Status: %v!\n"+
		"%v\n"+
		"  Expected Value: %v\n"+
		"Calculated Value: %v\n"+
		" Big Float Value: %v\n"+
		"%v\n",
		lineBreakStr,
		funcName,
		homeFileName,
		lineBreakStr,
		statusStr,
		lineBreakStr,
		expectedResultStr,
		calcValueStr,
		calcFloatValue.Text('f', int(requiredFractionalDigits)),
		lineBreakStr)

	return
}

const (
	// maxSciExp caps the absolute value of a scientific notation exponent (e.g., 1e10000)
	// to prevent big.Rat.SetString from allocating huge buffers.
	maxSciExp = 10000

	// maxFracDigits caps requested fractional digits to prevent int(reqFracDigits)
	// overflow and excessive memory allocation.
	maxFracDigits = 1_000_000
)

// PowDecimal raises baseStr to a non-negative integer exponent (exp) and returns:
//  1. A fixed-point decimal string rounded to reqFracDigits fractional places.
//  2. A high-precision *big.Float representing that exact rounded result.
//
// Special Cases:
//   - Any valid base raised to exp == 0 (including "0"^0 and "-2.5"^0) returns 1.
//   - exp == 1 performs direct rational rounding without an exponentiation loop.
//   - Moderate inputs (exp <= 250 and reqFracDigits <= 1000) use exact [big.Rat]
//     arithmetic; larger inputs automatically switch to guard-buffered [big.Float]
//     exponentiation by squaring.
//
// Input Format Requirements (baseStr):
//
// baseStr must be a non-empty base-10 decimal or scientific notation string with no
// whitespace, underscores ('_'), or fraction slashes ('/').
//
// Allowable scientific notation strings follow the structure:
//
//	[Sign][Mantissa][ExponentMarker][ExponentSign][ExponentDigits]
//
// Where:
//   - Sign (optional): '+' or '-' at index 0 only.
//   - Mantissa (required): Base-10 digits ('0'-'9') with at most one decimal point
//     ('.'). Leading or trailing decimal points (e.g., ".5" or "5.") are permitted.
//   - ExponentMarker (required for scientific notation): 'e' or 'E'.
//   - ExponentSign (optional): '+' or '-' immediately following 'e' or 'E'.
//   - ExponentDigits (required): Base-10 integer digits representing a signed
//     exponent within [-maxSciExp, maxSciExp] ([-10000, 10000]).
//
// Valid Examples:
//
//	"-2.5", "0", "-1.25e-10", "2.5E+4", ".5e-3", "1e10000"
//
// Invalid Examples:
//
//	" -2.5"      (leading/trailing whitespace)
//	"1_000.5e2"  (underscore separators)
//	"1/3"        (rational fraction slash)
//	"12-3.4e5"   (minus sign not at index 0)
//	"1.25e1-0"   (exponent minus sign not immediately after 'e'/'E')
//	"1e10001"    (exponent exceeds [-maxSciExp, maxSciExp])
//
// Errors:
//
// Returns an error if baseStr violates the format or exponent bounds above, or if
// reqFracDigits exceeds maxFracDigits (1,000,000).
func PowDecimal(baseStr string, exp uint, reqFracDigits uint) (string, *big.Float, error) {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"floatmechops/AppTest/Test003/bigFloatPow.go", // sourceFile
		"",             // parentType (leave "" if none)
		"PowDecimal()") // thisFunction

	if reqFracDigits > maxFracDigits {
		return "", big.NewFloat(0),
			errGenericCfg.GetErrorNativeToFunction(
				fmt.Sprintf("reqFracDigits= %v", reqFracDigits),
				"",
				"",
				"",
				"Input parameter 'reqFracDigits' exceeds maximum allowed value.",
				fmt.Sprintf("ERROR: reqFracDigits (%d) exceeds maximum allowed (%d)", reqFracDigits, maxFracDigits))
	}

	if err := validateBaseFormat(baseStr); err != nil {

		return "", big.NewFloat(0), errGenericCfg.GetErrorFromFunctionCall(
			"",
			"validateBaseFormat(baseStr)",
			fmt.Sprintf("baseStr: %s", baseStr),
			"",
			err)
	}

	// Fast path: exp == 0 (includes 0^0 = 1 and (-x)^0 = 1)
	if exp == 0 {

		if _, ok := new(big.Rat).SetString(baseStr); !ok {

			return "", big.NewFloat(0),
				errGenericCfg.GetErrorNativeToFunction(
					"'ok' == false",
					"",
					"",
					"",
					fmt.Sprintf("'baseStr' = %v", baseStr),
					fmt.Sprintf("new(big.Rat).SetString(baseStr) returned 'ok' == false"))
		}

		ratOne := big.NewRat(1, 1)

		resStr := ratOne.FloatString(int(reqFracDigits))

		resFloat, ok := new(big.Float).SetPrec(calcBits(1, reqFracDigits, 1)).SetString(resStr)

		if !ok {

			return "", big.NewFloat(0), errGenericCfg.GetErrorFromFunctionCall(
				"",
				"new(big.Float).SetPrec(calcBits(1, reqFracDigits, 1)).SetString(resStr)",
				fmt.Sprintf("baseStr: %s", baseStr),
				"",
				fmt.Errorf("new(big.Rat).SetString() returned 'ok' == false"))

		}

		return resStr, resFloat, nil
	}

	// Fast path: exp == 1 (exact rounding with no multiplication loop)
	if exp == 1 {
		baseRat, ok := new(big.Rat).SetString(baseStr)
		if !ok {
			return "", nil, fmt.Errorf("invalid base number string: %q", baseStr)
		}
		resStr := baseRat.FloatString(int(reqFracDigits))
		finalPrec := calcBits(uint(baseRat.Num().BitLen()), reqFracDigits, 1)
		resFloat, _ := new(big.Float).SetPrec(finalPrec).SetString(resStr)
		return resStr, resFloat, nil
	}

	// TIER 1: Exact Rational Math for moderate exponents & fractional precision
	if exp <= 250 && reqFracDigits <= 1000 {
		baseRat, ok := new(big.Rat).SetString(baseStr)
		if !ok {
			return "", nil, fmt.Errorf("invalid base number string: %q", baseStr)
		}

		resRat := powRat(baseRat, exp)
		resStr := resRat.FloatString(int(reqFracDigits))

		finalPrec := calcBits(uint(resRat.Num().BitLen()), reqFracDigits, 1)
		resFloat, _ := new(big.Float).SetPrec(finalPrec).SetString(resStr)
		return resStr, resFloat, nil
	}

	// TIER 2: Guard-buffered big.Float for large exponents (>250) or >1,000 fractional digits
	coarse, ok := new(big.Float).SetPrec(64).SetString(baseStr)
	if !ok {
		return "", nil, fmt.Errorf("invalid base number string: %q", baseStr)
	}

	// Estimate integer bits needed using absolute magnitude with overflow protection
	var mant big.Float
	binExp := coarse.MantExp(&mant)
	var estIntBits uint64
	if binExp > 0 {
		hi, lo := bits.Mul64(uint64(binExp), uint64(exp))
		if hi != 0 {
			return "", nil, fmt.Errorf("result magnitude overflows precision limits")
		}
		estIntBits = lo
	}

	workPrec, err := calcBitsSafe(estIntBits, reqFracDigits, exp)
	if err != nil {
		return "", nil, err
	}

	baseFloat, _ := new(big.Float).SetPrec(workPrec).SetString(baseStr)
	resFloat := powFloat(baseFloat, exp, workPrec)
	resStr := resFloat.Text('f', int(reqFracDigits))

	// Re-parse the formatted string so the returned big.Float matches resStr
	finalFloat, _ := new(big.Float).SetPrec(workPrec).SetString(resStr)
	return resStr, finalFloat, nil
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
func validateBaseFormat(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("base number string cannot be empty")
	}
	if strings.TrimSpace(s) != s {
		return fmt.Errorf("base number string cannot have leading or trailing whitespace: %q", s)
	}
	if strings.ContainsAny(s, "/_") {
		return fmt.Errorf("invalid characters in base number string: %q", s)
	}

	// Split at 'e' or 'E' if scientific notation is present
	mantissa := s
	if idx := strings.IndexAny(s, "eE"); idx != -1 {
		mantissa = s[:idx]
		exponent := s[idx+1:]

		if len(mantissa) == 0 || len(exponent) == 0 {
			return fmt.Errorf("malformed scientific notation: %q", s)
		}
		// In the exponent, '-' is only allowed at index 0 of the exponent part
		if expMinus := strings.LastIndexByte(exponent, '-'); expMinus > 0 {
			return fmt.Errorf("exponent minus sign must immediately follow 'e'/'E': %q", s)
		}

		// Parse the exponent as an int64 to catch both malformed digits and huge magnitudes
		sciExp, err := strconv.ParseInt(exponent, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid or out-of-range scientific notation exponent in %q: %w", s, err)
		}
		if sciExp > maxSciExp || sciExp < -maxSciExp {
			return fmt.Errorf("scientific notation exponent %d exceeds allowed range [-%d, %d]", sciExp, maxSciExp, maxSciExp)
		}
	}

	// In the mantissa (base), '-' is only allowed at index 0
	if idx := strings.LastIndexByte(mantissa, '-'); idx > 0 {
		return fmt.Errorf("minus sign must be in the first string position: %q", s)
	}

	return nil
}

func calcBits(intBits uint, reqFracDigits uint, exp uint) uint {
	fracBits := uint(math.Ceil(float64(reqFracDigits) * 3.322))
	guardBits := uint(166 + bits.Len(exp)*2)
	return intBits + fracBits + guardBits
}

func calcBitsSafe(intBits uint64, reqFracDigits uint, exp uint) (uint, error) {
	fracBits := uint64(math.Ceil(float64(reqFracDigits) * 3.322))
	guardBits := uint64(166 + bits.Len(exp)*2)
	total := intBits + fracBits + guardBits
	if total > uint64(big.MaxPrec) {
		return 0, fmt.Errorf("required precision (%d bits) exceeds big.MaxPrec", total)
	}
	return uint(total), nil
}

func powRat(base *big.Rat, exp uint) *big.Rat {
	result := big.NewRat(1, 1)
	b := new(big.Rat).Set(base)
	for exp > 0 {
		if exp&1 == 1 {
			result.Mul(result, b)
		}
		exp >>= 1
		if exp > 0 {
			b.Mul(b, b)
		}
	}
	return result
}

func powFloat(base *big.Float, exp uint, prec uint) *big.Float {
	result := new(big.Float).SetPrec(prec).SetFloat64(1)
	b := new(big.Float).SetPrec(prec).Set(base)
	for exp > 0 {
		if exp&1 == 1 {
			result.Mul(result, b)
		}
		exp >>= 1
		if exp > 0 {
			b.Mul(b, b)
		}
	}
	return result
}
