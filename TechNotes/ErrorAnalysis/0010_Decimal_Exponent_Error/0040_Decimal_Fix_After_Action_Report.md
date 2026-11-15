# 2026-11-15 02:20-Hours Decimal_Fixed_Passing

## 1. Problem Summary
A precision bug was identified in the `PowExact` operation where calculating $6.9857^9$ returned `39617710.92773922` instead of the expected value (approximately `39617710.92824073`). The root cause was located in `NewFromFloat64` within `decimal_0010.go`, which calculated the precision (`prec`) based on the number of fractional digits. For inputs with a small number of fractional digits, `prec` could fall below 53 bits, causing the truncation of the `float64` mantissa and introducing relative error that was amplified during exponentiation.

## 2. Problem Solving Strategy employed
The strategy involved:
- **Root Cause Analysis**: Analyzing the precision calculation in `NewFromFloat64` and determining that it did not guarantee the preservation of the full 53-bit mantissa of a `float64`.
- **Fix Implementation**: Modifying the precision logic to ensure a minimum `prec` value of 53 bits, regardless of the number of fractional digits.
- **Verification**: Implementing a regression test case in `zzzt_decimal_0010_test.go` to specifically validate the result of $6.9857^9$ against the expected high-precision value.

## 3. Source Code Files Modified
- `floatmech/decimal_0010.go`: Updated `NewFromFloat64` to enforce minimum 53-bit precision.
- `floatmech/zzzt_decimal_0010_test.go`: Added/updated test cases to verify the fix.

## 4. Test Use to Verify Solution
The solution was verified using the Go test suite, specifically the `TestDecimal_Pow_0010` function in `floatmech/zzzt_decimal_0010_test.go`, which confirmed that the result of the exponentiation now matches the expected precision.



# Background

Prior to this fix, source code file **`floatmech/decimal_0040.go`**, method ***zivBounder.propagateMulErr()*** was modified to enhance the error bound calculation. This method computes the upper bound on relative error (in units of $u$) after multiplying two values with error bounds $errA$ and $errB$, based on the derivation:
$(1 + errA \cdot u)(1 + errB \cdot u) - 1 = (errA + errB + errA \cdot errB \cdot u) \cdot u$

Specifically, the safety margin factor was increased from **1.0** to **1.1** to account for accumulated drift in binary exponentiation.

```go
// propagateMulErr computes the upper bound on relative error (in units of u)
// after multiplying two values with error bounds errA and errB.
func (z zivBounder) propagateMulErr(errA, errB float64, acc big.Accuracy) float64 {
	// (1 + errA*u)*(1 + errB*u) - 1 = (errA + errB + errA*errB*u)*u
	prop := errA + errB + errA*errB*z.unitRoundoff
	if acc == big.Exact {
		return prop
	}
	// Multiply by (1 + 1.1*u) for the new rounding step.
	// The 1.1 factor is a safety margin to account for accumulated drift in binary exponentiation.
	return prop + 1.1 + prop*z.unitRoundoff
}
```
