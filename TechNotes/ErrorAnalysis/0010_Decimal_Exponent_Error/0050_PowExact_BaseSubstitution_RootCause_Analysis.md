# Decimal.PowExact — 10.81233^8 / 15.81233^4 / 1.0691372836118^24 Root-Cause Analysis

Date: 2026-11-21
Files: `decimal_0001000.go` (NewFromFloat64, powWithStrategy), `decimal_0003000.go` (powZiv), `decimal_0004000.go` (zivBounder)
Related: `0030_Required_Corrective_Action_Precision_Discrepancy_Error.md`, `0040_Decimal_Fix_After_Action_Report.md` (the 6.9857^9 issue)

## 1. Symptom

Three of the four new `PowExact` tests fail, each off in the last few fractional digits:

| case | expected | PowExact (code) | exact binary64 base | exact decimal-string base |
|------|----------|-----------------|---------------------|---------------------------|
| 7.8715^6, tfd=10 | 237872.9172319539 | 237872.9172319539 (PASS) | 237872.9172319539 | 237872.9172319539 |
| 10.81233^8, tfd=14 | 186790307.75412594147937 | 186790307.75412584893490 | 186790307.75412584893490 | 186790307.75412594147937 |
| 15.81233^4, tfd=14 | 62514.89090130012128 | 62514.89090130011070 | 62514.89090130011070 | 62514.89090130012128 |
| 1.0691372836118^24, tfd=14 | 4.97511810598396 | 4.97511810598397 | 4.97511810598397 | 4.97511810598396 |

## 2. Method

Diagnostic test `zzzt_diag_temp_test.go` computed each result three ways:
- [A] exact (4000-bit) exponentiation of the **stored base** `x` (the exact binary64 value `NewFromFloat64` holds);
- [B] exact exponentiation of the **decimal string** `fmt.Sprintf("%v", base)`;
- [C] the real `PowExact` code path.

All three bases store the value as `Acc = Exact` (`MinPrec = 52` or `53`). The 53-bit floor in `NewFromFloat64` (`decimal_0001000.go:378-380`) is present and working.

## 3. Finding

**`PowExact` is correct relative to the base it is given.** For all three failing cases,
`PowExact` (C) **exactly equals** [A] — exact exponentiation of the stored binary64 base.
The "expected" values equal [B] — exact exponentiation of the *decimal literal*, which is a
different number.

Concretely:
- The float64 `10.81233` is the binary64 value `10.8123299999999997316…`, **not** `10.81233`.
- Raising that to the 8th power gives `186790307.7541258489…`, which `PowExact` returns.
- Raising the decimal `10.81233` to the 8th power gives `186790307.7541259414…` (the "expected").

The discrepancy is a base-substitution error amplified by exponentiation, not a precision
failure in the Ziv machinery. `powZiv`/`zivBounder` are converging and verifying correctly;
the interval check passes because the true value (of the stored base) is being represented
with enough precision to round unambiguously to `tfd` digits.

## 4. Why the earlier "1% → 1.1%" fixes did not cure this

Those changes widened the error bound so `tryRound` rejects unstable results and forces more
guard bits. They are orthogonal to this problem: there is no instability to reject. The stored
base is simply a different quantity than the typed decimal. Widen the bound all you like —
`PowExact` will keep returning the exact power of the binary64 base.

## 5. Distinction from the 6.9857^9 case (0040 report)

- 6.9857^9: the stored base had been *truncated below the full 53-bit mantissa*, so the base
  itself was wrong. Fix = enforce the 53-bit floor. (That floor is now present.)
- These cases: the stored base is the **full, exact** binary64 value (Acc=Exact). It is a
  different number than the user's decimal literal by construction of float64. No amount of
  working precision recovers the decimal-literal value from the binary64 base.

## 6. Recommendation (decision, not a code defect)

This is a **semantic** choice, not a bug. Pick one and document it:

- **Option A (current): "pow of the exact float64 value."** `PowExact` is then correct by
  construction. The tests' expected values are wrong — they should be derived from the
  binary64 base (e.g., via `math.Pow`/big.Float at the stored base). Keep the code; fix the
  test expectations. This matches mpmath/Python `**` behavior.

- **Option B: "pow of the intended decimal literal."** Requires `NewFromFloat64` to store a
  base that carries the decimal identity, e.g. parse the shortest round-trip decimal
  `strconv.FormatFloat(v, 'f', -1, 64)` (verified to equal the intended literal for all four
  cases) via `SetString` into `big.Float` at sufficient precision, instead of `SetFloat64`.
  Trade-off: this makes `NewFromFloat64` no longer store the exact binary64 value, which may
  change every other operation's semantics and is a much larger, riskier change.

Recommendation: adopt **Option A** for `PowExact` and correct the test expectations; treat
Option B as a separate, deliberate API change only if a true decimal-literal semantic is
required across the whole `Decimal` type.

## 7. Files involved
- `decimal_0001000.go:366-390` — `NewFromFloat64` (53-bit floor; stores exact binary64 base).
- `decimal_0001000.go:588-625` — `powWithStrategy` dispatch (Ziv path for `PowExact`).
- `decimal_0003000.go:116-158` — `powZiv` (converges and verifies correctly).
- `decimal_0004000.go:33-92` — `propagateMulErr` / `tryRound` (bound is sound; not the cause).
