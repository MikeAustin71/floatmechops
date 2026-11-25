# Stage 2 — Ziv Engine Precision Defect — Progress Notes

Last updated: 2026-11-25 (context checkpoint)

## Objective
- Identify and fix the root cause of `TestDecimal_PowExact_00020` failure (`10.81233^8`),
  where the engine returns `186790307.75412594147769` instead of the exact
  `186790307.75412594147937` (off by ~168 ulp at tfd=14).

## Failing Test
- `TestDecimal_PowExact_00020` at `zzzt_decimal_0001030_test.go:61`.
- Input: `vNumStr = "10.81233"`, `targetFractionalDigits = 14`, `mode = big.ToNearestAway`, `exponent = 8`.
- Expected: `186790307.75412594147937`
- Actual:   `186790307.75412594147769`
- Base construction confirmed correct (Stage 1 complete). Error is entirely in the exponentiation engine path.
- `floatNumStrPadding` tested at 30 and 100 — no impact — confirms defect is downstream of base construction.

## Current Engine File Layout (RESTRUCTURED — superseded earlier notes)
> NOTE: The engine was refactored since 0070. The Ziv path is no longer the
> whole story; an **exact dyadic kernel** now exists.

- `decimal_0001000.go`: `Decimal` struct, `powStrategy` enum, `powWithStrategy` entry point.
- `decimal_0002000.go`: `precConverter.toPrec` / `toFractionalDigits` — working-precision derivation.
- `decimal_0002010.go`: `constructFromNumStr` / `NewFromFloatNumStr` — Stage 1 base construction (confirmed correct).
- `decimal_0003000.go`: `powZiv`, `estimatePowExp2`, `powSpecialCase` — Ziv engine exponentiation; instantiates `dyadicPowKernel` (line 211).
- `decimal_0004000.go`: `zivBounder`, `propagateMulErr`, `tryRound` — bounding logic.
- `decimal_0005000.go`: **NEW** `dyadicPowKernel` — exact `big.Int.Exp` path (decompose → exact pow → single round).

## Key Kernel Facts (decimal_0005000.go)
- `canUseExact`: returns false unless `exp > 1` and `minPrec <= maxExactBits` and `exp <= maxExactBits/minPrec`.
- `powExact`: `checkExponentBounds` → `decompose` (x = m * 2^scale) → `big.Int.Exp(m, exp)` → `SetMantExp` → single rounding to `(targetPrec, targetMode)`.
- `decompose`: `mant.SetMantExp(mant, int(p))` (right-shift by p = MinPrec) then `mant.Int(nil)` must be `big.Exact`.
- `maxExactBits` is set where the kernel is constructed in `decimal_0003000.go` (line 211).

## Work State
### Completed
- Stage 1: `NewFromFloatNumStr` created; all tests migrated; base-substitution root cause eliminated.
- Root-cause analysis (0050) confirmed `PowExact` is correct relative to its input base.
- Padding experiment (30→100) confirmed defect is downstream of base construction.

### Active
- Determine **which engine path** (exact `dyadicPowKernel` vs. Ziv) is actually taken for
  `10.81233^8` at tfd=14, and pinpoint why the final rounding is ~168 ulp low.

### Blocked
- (none)

## Next Move
1. Read `decimal_0003000.go` (kernel construction + `maxExactBits` + dispatch) to determine
   whether `canUseExact` is TRUE or FALSE for this input.
2. Read `decimal_0002000.go` (`toPrec`) to confirm the `targetPrec` / working precision.
3. Add a temporary diagnostic (in `AppTest/`) to print `canUseExact` result and the
   intermediate `m`, `scaleExp2`, and the single rounding — to see if the exact path is
   being taken and where the precision is lost.
4. Implement the fix in whichever engine path is implicated; re-run `TestDecimal_PowExact_00020`.

## Relevant Files
- `D:\GoProjects\MikeAustin71\floatmechops\floatmech\zzzt_decimal_0001030_test.go`
- `D:\GoProjects\MikeAustin71\floatmechops\floatmech\decimal_0003000.go`
- `D:\GoProjects\MikeAustin71\floatmechops\floatmech\decimal_0004000.go`
- `D:\GoProjects\MikeAustin71\floatmechops\floatmech\decimal_0002000.go`
- `D:\GoProjects\MikeAustin71\floatmechops\floatmech\decimal_0005000.go`
- `D:\GoProjects\MikeAustin71\floatmechops\AppTest\Test003\bigFloatPow.go` (user verification harness; `calculateBigFloatPrecision` uses `*3.32192809488736 + 64` safety)
