# Decimal.PowExact — Stage 2 Plan (Ziv-Engine Precision Defect)

Date created: 2026-11-21
Status: **PLANNED** (work begins next session)
Predecessor: `0060_NewFromFloatNumStr_Stage1_Report.md`

---

## 1. Scope

Stage 2 addresses the **single remaining failing test** exposed by the Stage 1 migration:

- `TestDecimal_PowExact_00020` — 10.81233^8 (tfd=14, `big.ToNearestAway`)
  - expected (exact): `186790307.75412594147937`
  - engine returned:  `186790307.75412594147769`
  - magnitude of error: ~1.68e-4, ≈ 100 ulp in the last digit

The base is now the exact decimal literal (confirmed by Stage 1), so the error is
**entirely in the Ziv engine path** (`powZiv` / `zivBounder`).



## 1-a User Addendum - 2026-11-22 15:51-Hours

I increased the value of **`const floatNumStrPadding uint`** to **`100`** from its prior value of **`30`**. I then ran the test that previously failed, **`TestDecimal_PowExact_00020`**. The test still failed with the same results, even after increasing **`floatNumStrPadding`** to a value of **`100`**.  My results showed no impact resulting form the higher **`floatNumStrPadding`** value:

```
=== RUN   TestDecimal_PowExact_00020
    zzzt_decimal_0001030_test.go:104: 
        Expected decimal.PowExact = 186790307.75412594147937
        Instead decimal.PowExact  = 186790307.75412594147769
--- FAIL: TestDecimal_PowExact_00020 (0.00s)

```

I subsequently reset the value of **`floatNumStrPadding`** to its original value of ***30***.



---

## 2. Ground Truth (verified, immutable reference)

| case | tfd | exact value (big.Int) |
|------|-----|-----------------------|
| 6.9857^9 | 8 | 39617710.92824073 |
| 7.8715^6 | 10 | 237872.9172319539 |
| 10.81233^8 | 14 | 186790307.75412594147937 |
| 15.81233^4 | 14 | 62514.89090130012128 |
| 1.0691372836118^24 | 14 | 4.97511810598396 |

Source: `AppTest/main.go` (pure big.Int arithmetic). Cross-checked by user in Python
and by `AppTest/Test001/verify_pow.go` (user-fixed).

---

## 3. Investigation Plan

### 3.1 Isolate the divergent multiplication

The Ziv engine computes `base^8` via 8 multiplications (or a divide-and-conquer tree).
For each intermediate product, compare the engine's partial result against the exact
partial product computed by `exactPow`-style integer arithmetic.

Goal: identify **which multiplication** first diverges beyond the working precision.

### 3.2 Check working precision at each step

For the failing case, `effectiveTfd = 14 + 30 = 44` decimal digits.
The derived working precision is:

```
fracBits = ceil(44 * log2(10)) ≈ ceil(44 * 3.32193) ≈ 147
intBits  = floor(log2(10.81233)) = 3
safetyBits = 8
total    ≈ 158 bits
```

Verify that 158 bits is sufficient to hold 10.81233^8 (≈ 1.87e8) to 44 decimal
digits. It should be — the integer part is only 9 digits — but confirm empirically.

### 3.3 Check the Ziv bounder tolerance

`zivBounder` (constructed by `newZivBounder` at `decimal_0004000.go:22`), with
`propagateMulErr` (`decimal_0004000.go:33`) and `tryRound`
(`decimal_0004000.go:46`), computes an error bound and decides whether the current
result is "stable enough" to round to `tfd`. If the bound is too loose, the engine may
accept a result that is off by 1–2 ulp in the last digit.

Investigate:
- The bound formula and whether it accounts for accumulated rounding over 8 multiplications.
- Whether `tryRound` is correctly comparing the bound against the rounding threshold
  (half-ulp at `tfd`).
- Whether the 30-char padding is inflating `effectiveTfd` in a way that interacts
  poorly with the bound.

### 3.4 Check `powZiv` iteration structure

`powZiv` (`decimal_0003000.go:116`) and its driver `runZivPass`
(`decimal_0003000.go:162`) — confirm:
- Is it doing sequential multiplication (`result *= base`, 8 times) or exponentiation
  by squaring?
- Does it re-derive precision at each step, or carry a fixed precision?
- Is the final rounding step using the correct `tfd` (14) or the inflated
  `effectiveTfd` (44)?

### 3.5 Scratch program

Create `AppTest/Test002/` (or extend `AppTest/main.go`) to:

1. Import `github.com/MikeAustin71/floatmechops/floatmech` (capital M, per go.mod).
2. Call `NewFromFloatNumStr("10.81233", 14, big.ToNearestAway)` then `PowExact(8)`.
3. Print the result and compare against `186790307.75412594147937`.
4. Add intermediate diagnostics: after each multiplication in the engine, print
   the partial product at high precision.

The user has confirmed `AppTest/Test001/verify_pow.go` may be freely modified.
A new `AppTest/Test002/` directory is preferred to avoid disturbing existing files.

---

## 4. Likely Root Causes (hypotheses, in order of probability)

1. **Insufficient working precision** at one or more intermediate steps — the engine
   rounds each product to ~158 bits, accumulating error over 8 multiplications that
   exceeds the final rounding threshold.

2. **Ziv bounder too loose** — the error bound underestimates accumulated rounding
   error, so `tryRound` accepts an unstable result.

3. **30-char padding interaction** — `effectiveTfd = 44` inflates the working
   precision but the final rounding to `tfd = 14` may lose precision if the
   intermediate results are not held at sufficient extra precision beyond the
   44-digit working target.

4. **Rounding mode inconsistency** — an intermediate step uses `ToNearestEven` while
   the final step uses `ToNearestAway`, producing a 1-ulp shift in the last digit.

---

## 5. Acceptance Criteria for Stage 2

- [ ] `TestDecimal_PowExact_00020` passes (result = `186790307.75412594147937`).
- [ ] All other tests continue to pass (14 currently passing).
- [ ] `go build ./...` and `go vet ./floatmech/` remain clean.
- [ ] Root cause documented with a specific code reference (file:line).
- [ ] Fix does not regress the 30-char padding stop-gap (flagged separately).

---

## 6. Files to Read First (next session)

| file | what to look for |
|------|-----------------|
| `floatmech/decimal_0003000.go:116,162` | `powZiv` / `runZivPass` — iteration structure, precision per step |
| `floatmech/decimal_0004000.go:22,33,46` | `newZivBounder`, `propagateMulErr`, `tryRound` — bound formula |
| `floatmech/decimal_0002000.go` | `toPrec` / `toFractionalDigits` — precision derivation |
| `floatmech/decimal_0001000.go:659` | `powWithStrategy` — how `powZiv` is invoked, what precision is passed |
| `floatmech/zzzt_decimal_0001030_test.go:61-112` | The failing test — exact expected/actual strings |

---

## 7. Scratch Program Import Path

```go
import (
    "github.com/MikeAustin71/floatmechops/floatmech"
    "math/big"
)
```

Note: capital `M` in `MikeAustin71` (per `go.mod`). The package name is `floatmech`
(lowercase). Public methods like `NewFromFloatNumStr` and `PowExact` are directly
callable on a `*floatmech.Decimal` instance.

---

## 8. Open Items Carried Forward from Stage 1

1. **30-char padding stop-gap** — flagged for replacement with a principled precision
   model. Do NOT change this in Stage 2; fix the engine first.
2. **`decimal_0002010.go:90-91`** — empty error-context strings. Cosmetic; defer.
3. **Docs cleanup** — `0050_PowExact_BaseSubstitution_RootCause_Analysis.md` contains
   the now-superseded base-substitution conclusion. Add a note at the top referencing
   the 0060 report once Stage 2 is complete.
