# Decimal.PowExact — `NewFromFloatNumStr` Migration — Stage 1 Report

Date: 2026-11-21
Package: `floatmech`
Module: `github.com/MikeAustin71/floatmechops`
Go: 1.27.3 (windows/amd64)
Related:
- `0030_Required_Corrective_Action_Precision_Discrepancy_Error.md`
- `0040_Decimal_Fix_After_Action_Report.md` (the 6.9857^9 issue)
- `0050_PowExact_BaseSubstitution_RootCause_Analysis.md` (base-substitution finding)

---

## 1. Objective

Stage 1 of a staged plan agreed with the user:

1. **Create** `NewFromFloatNumStr`, a string-based constructor that replaces `NewFromFloat64`.
   The incoming numeric string is padded with a fixed number of trailing `0` characters
   before parsing so the derived working precision retains the requested fractional digits.
2. **Migrate** every test call site from `NewFromFloat64` to `NewFromFloatNumStr`.
3. **Verify** results against independent `big.Int` / `big.Float` ground truth.
4. **Report** the results, including the successful cases with their exact `big.Int` values.

Downstream Ziv-engine precision failures are explicitly **deferred** to Stage 2.

### Padding directive (stop-gap)

Per the user, append **exactly 30 `"0"` characters** to every `vNumStr` before parsing.
This is a deliberate stop-gap and is flagged for revisit; it is not intended as a
final design.

---

## 2. What was done

### 2.1 New private mechanics — `decimal_0002010.go`

New file added. It defines `numStrConstructorMechanics` and the private method
`constructFromNumStr(vNumStr string, targetFractionalDigits uint, mode big.RoundingMode)
(*Decimal, error)`.

Mechanics, in order:

1. Pad the input: append `floatNumStrPadding = 30` zero characters to `vNumStr`.
2. Parse the padded string with `(*big.Float).SetString` at infinite precision.
3. Decompose via `MantExp` to obtain the binary exponent.
4. Derive the working precision with
   `defaultConverter().toPrec(effectiveTfd, exp2)`, where
   `effectiveTfd = targetFractionalDigits + floatNumStrPadding`.
5. Finalize: `SetPrec(prec).SetMode(mode).Set(exact)`.

### 2.2 Public entry point — `decimal_0001000.go`

`NewFromFloat64` was **replaced** by `NewFromFloatNumStr` at `decimal_0001000.go:390`.
The public method:

- uses a pointer receiver and acquires `d.mu` (public entry points are mutex-protected);
- validates the nil receiver;
- delegates all real work to the private
  `new(numStrConstructorMechanics).constructFromNumStr(...)` in `decimal_0002010.go`;
- wraps any returned error with `ErrorGenericFactory`.

This is the intended public/private split: the public entry point owns validation and
locking; the private mechanics file owns the calculation.

### 2.3 Zero-base call fixed

The zero-base construction inside the `PowExact` path
(`decimal_0001000.go:666`) now calls
`new(Decimal).NewFromFloatNumStr("0", defaultFractionalDigits, big.ToNearestEven)`.

### 2.4 Test migration

All four test files were migrated to `NewFromFloatNumStr` with string arguments and
explicit error handling:

| Test file | Notes |
|-----------|-------|
| `zzzt_decimal_0001000_test.go` | Migrated. **Expected value corrected** (see §5). |
| `zzzt_decimal_0001010_test.go` | Migrated. Two `=` → `:=` compile fixes (receiver and copyDestination). |
| `zzzt_decimal_0001020_test.go` | Migrated. |
| `zzzt_decimal_0001030_test.go` | Migrated. |

A full-repo grep confirms **no remaining `NewFromFloat64` references in Go source**
(only stale mentions remain in docs / research notes).

### 2.5 Build and vet

- `go build ./...` — clean.
- `go vet ./floatmech/` — clean.

---

## 3. Ground truth — independent `big.Int` / `big.Float`

Computed by the scratch program `AppTest/main.go` (exact integer arithmetic), and
independently cross-checked by the user in Python. `AppTest/Test001/verify_pow.go`
was also fixed by the user (its `res` accumulator was initialized to `0` instead of
`1`) and now agrees.

| case | fractional digits | exact `big.Int` / `big.Float` value |
|------|-------------------|-------------------------------------|
| 6.9857^9 | 8 | **39617710.92824073** |
| 7.8715^6 | 10 | **237872.9172319539** |
| 10.81233^8 | 14 | **186790307.75412594147937** |
| 15.81233^4 | 14 | **62514.89090130012128** |
| 1.0691372836118^24 | 14 | **4.97511810598396** |

---

## 4. Test results

Command: `go test ./floatmech/ -run 'TestDecimal' -v`

Result: **14 PASS / 1 FAIL.**

### Passing

- `TestDecimal_Pow_0010` — 6.9857^9 → `39617710.92824073` (matches ground truth).
- `TestDecimal_NewFromFloatNumStr_0010` — 6.9857^9 → `39617710.92824073`.
- `TestDecimal_CopyIn_NilReceiver`
- `TestDecimal_CopyIn_NilCopySource`
- `TestDecimal_CopyIn_ValidCopy`
- `TestDecimal_CopyOut_NilReceiver`
- `TestDecimal_CopyOut_NilCopyDestination`
- `TestDecimal_CopyOut_ValidCopy`
- `TestDecimal_CopyIn_NilCopySourceVal`
- `TestDecimal_CopyOut_NilHostVal`
- `TestDecimal_Clone_NilHostVal`
- `TestDecimal_PowExact_00010` — 7.8715^6 → `237872.9172319539` (matches ground truth).
- `TestDecimal_PowExact_00030` — 15.81233^4 → `62514.89090130012128` (matches ground truth).
- `TestDecimal_PowExact_00040` — 1.0691372836118^24 → `4.97511810598396` (matches ground truth).

### Failing

- `TestDecimal_PowExact_00020` — 10.81233^8:
  - expected (exact) = `186790307.75412594147937`
  - engine returned  = `186790307.75412594147769`

  The two agree through the first 12 fractional digits
  (`…75412594147`) and diverge in the final two digits.

---

## 5. The corrected 6.9857^9 expected value

The original test asserted `…92824071`. That value was **wrong**. The independent
`big.Int`/`big.Float` computation gives `39617710.92824073`, and the engine returns
the same. The expected constant in `zzzt_decimal_0001000_test.go` was corrected from
`071` to `073`. The test now passes against a verified exact value.

---

## 6. Re-characterization of the 10.81233^8 failure

This is the key finding of Stage 1.

### 6.1 The 0050 conclusion no longer holds

`0050_PowExact_BaseSubstitution_RootCause_Analysis.md` concluded that the 10.81233^8
failure was a **base-substitution** error: `NewFromFloat64` stored the binary64 base
`10.8123299999999997316…`, and `PowExact` was correctly computing the power of *that*
base, not of the decimal literal. Under that framing, the code was correct and the test
expectation was the thing to change.

### 6.2 Why the migration invalidates that framing

With `NewFromFloatNumStr`, the base is now constructed from the **decimal string**
`"10.81233"` (zero-padded, then parsed). The stored base is therefore the exact decimal
literal, not a binary64 surrogate. That removes the base-substitution cause entirely:

- All four initial-value string checks in the migrated tests pass, confirming the stored
  base is the decimal literal to the requested precision.
- The expected value `…94147937` is the correct target for the decimal-literal base.

### 6.3 What the failure now is

The engine returns `…94147769`. That is a **genuine Ziv-engine precision / rounding
discrepancy**: the result is off by roughly `1.68e-4`, about 100 ulp in the final digit
at 14 fractional digits. This is a real defect in the `powZiv` / `zivBounder` path, not
a test-expectation problem.

This is a **newly exposed** defect. Before the migration it was masked by the
base-substitution artifact; the migration resolved the base-substitution issue and
surfaced the underlying engine bug.

### 6.4 Migration impact summary

| case | before migration | after migration |
|------|------------------|-----------------|
| 7.8715^6 | FAIL (base-substitution) | **PASS** |
| 10.81233^8 | FAIL (base-substitution) | FAIL (**engine precision**) |
| 15.81233^4 | FAIL (base-substitution) | **PASS** |
| 1.0691372836118^24 | FAIL (base-substitution) | **PASS** |

3 of the 4 previously-failing cases now pass. The single remaining failure has a
different, deeper cause.

---

## 7. Known issues / open items

1. **30-character padding is a stop-gap.** It inflates the working precision by 30
   decimal digits of headroom. It is not a principled precision model and must be
   replaced with a deliberate precision derivation in a later stage.
2. **10.81233^8 is a real engine defect** in the Ziv path (`powZiv` / `zivBounder`).
   Deferred to Stage 2. The engine must be brought to within rounding distance of the
   exact decimal-literal result `186790307.75412594147937`.
3. **`decimal_0002010.go` lines 90-91** carry empty error-context string literals
   (`""`). Flagged for cleanup so error context is meaningful.
4. **Scratch programs** live under `AppTest/` (not the OS temp dir, per directive).
   `AppTest/Test001/verify_pow.go` is now correct (the `res = 0` accumulator bug is
   fixed) and agrees with `AppTest/main.go`.

---

## 8. Files involved

| file | role |
|------|------|
| `floatmech/decimal_0001000.go:390` | `NewFromFloatNumStr` public entry point |
| `floatmech/decimal_0002010.go` | `numStrConstructorMechanics.constructFromNumStr` (private mechanics, padding) |
| `floatmech/decimal_0002000.go` | `precConverter.toPrec` / `toFractionalDigits` |
| `floatmech/decimal_0003000.go` | `powZiv` (deferred Stage 2) |
| `floatmech/zzzt_decimal_0001000_test.go` | 6.9857^9 tests; expected value corrected 071→073 |
| `floatmech/zzzt_decimal_0001010_test.go` | CopyIn / CopyOut tests (migrated) |
| `floatmech/zzzt_decimal_0001020_test.go` | nil-val tests (migrated) |
| `floatmech/zzzt_decimal_0001030_test.go` | PowExact tests (migrated); 10.81233^8 still fails |
| `AppTest/main.go` | `big.Int` ground-truth scratch program |
| `AppTest/Test001/verify_pow.go` | user-fixed verification program |

---

## 9. Stage 1 verdict

- Constructor migration: **complete**. `NewFromFloatNumStr` is in place, all call sites
  migrated, build and vet clean.
- Correctness vs. exact ground truth: **4 of 5 power cases verified exact**, including
  the corrected 6.9857^9 value.
- One genuine engine defect remains (10.81233^8), now correctly attributed to the Ziv
  path and deferred to Stage 2.
- The 30-character padding stop-gap is in place and explicitly flagged for replacement.
