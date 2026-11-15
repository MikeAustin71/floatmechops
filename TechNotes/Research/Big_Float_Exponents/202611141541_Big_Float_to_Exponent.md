



# AI Studio 2026-11-14 15:41 Big Float to Exponent



# Gemini 4 Hydrogen Question No 1 2026-11-14 15:42

New Session Title: 'go' big.Float Raised to an Exponent

Assume that I am writing a 'go' program.

In terms of architecture and design, what are my options for creating one or more custom types with methods which will calculate the result of a big.Float value raised to an exponent of type uint.



# Gemini 4 Hydrogen Response No 1

In Go, you cannot add methods directly to `big.Float` because it lives in the `math/big` package, and `big.Float` does not ship with a built-in `Pow` or `Exp` method. 

When designing custom types to provide `base ^ exp` (typically implemented via binary exponentiation / exponentiation by squaring), your architectural options fall into four main patterns, along with a few critical cross-cutting design choices specific to `math/big`.

---

### 1. Type-Design Patterns

#### Option A: Struct Embedding (`*big.Float` or `big.Float`)
You embed `big.Float` inside your custom struct so that all existing `big.Float` methods are promoted to your type automatically.

```go
type Float struct {
    *big.Float
}

func NewFloat(x float64) *Float {
    return &Float{Float: big.NewFloat(x)}
}

func (z *Float) Pow(x *Float, exp uint) *Float {
    // Binary exponentiation updating z.Float
    return z
}
```
* **Pros:** Zero boilerplate to keep standard methods (`Add`, `Mul`, `Text`, `Cmp`, etc.). Drop-in compatibility wherever you need `big.Float` behavior.
* **Cons:** **Method-chaining breaks.** Calling `myFloat.Add(a, b)` returns a `*big.Float`, *not* your `*Float` type, so you cannot chain `.Add(...).Pow(...)` unless you manually shadow/override every `big.Float` method. Also, if you embed `*big.Float` (a pointer), a zero-value `Float{}` will panic on method calls if not initialized via a constructor. *(Embedding `big.Float` by value avoids the nil panic, but introduces value-copying risks—see Section 2).*

#### Option B: Encapsulated Wrapper (Composition)
You make `big.Float` an unexported field inside a custom struct and explicitly define the API you want to expose.

```go
type Decimal struct {
    val  *big.Float
    prec uint
    mode big.RoundingMode
}

func (d *Decimal) Pow(exp uint) *Decimal { ... }
func (d *Decimal) Add(other *Decimal) *Decimal { ... }
func (d *Decimal) Raw() *big.Float { ... } // optional escape hatch
```
* **Pros:** 
  * Clean, consistent API (you can make it fluent/chainable, or even immutable).
  * Protects invariants—for example, guaranteeing every instance uses the exact same precision and rounding mode.
  * Hides `math/big` quirks (like accidental value copying).
* **Cons:** Highest boilerplate if callers also need standard arithmetic (`Add`, `Sub`, `Mul`, `Div`), because you must write forwarding methods for each one you expose.

#### Option C: Defined Type (`type MyFloat big.Float`)
You declare a new named type with `big.Float` as its underlying type.

```go
type MyFloat big.Float

func (z *MyFloat) Pow(x *MyFloat, exp uint) *MyFloat {
    // Convert between *MyFloat and *big.Float with zero allocation:
    bfZ := (*big.Float)(z)
    bfX := (*big.Float)(x)
    ...
    return z
}
```
* **Pros:** Zero memory overhead and zero-cost pointer conversions (`(*big.Float)(m)`).
* **Cons:** **Strips all existing methods.** `MyFloat` will *not* inherit `Add`, `Mul`, `SetString`, etc. Callers must constantly cast back and forth or you must re-wrap every method. In practice, this is usually less ergonomic than Option A or Option B.

#### Option D: A "Context / Calculator" Type (Separating Data from Behavior)
Instead of wrapping the number itself, you create a type that holds **configuration** (precision, guard bits, rounding mode) and operates on standard `*big.Float` values.

```go
type MathContext struct {
    Prec      uint
    GuardBits uint
    Mode      big.RoundingMode
}

func (mc *MathContext) Pow(z, base *big.Float, exp uint) *big.Float {
    // Perform calculation using mc.Prec + mc.GuardBits,
    // then round result into z at mc.Prec.
    return z
}
```
* **Pros:** 
  * Very idiomatic in Go when you only need a few specialized operations (`Pow`, `Sqrt`, `Log`, etc.).
  * Zero friction with external code or other packages that already pass around `*big.Float`.
  * Natural place to store workspace/scratch `big.Float` buffers (or a `sync.Pool`) to eliminate heap allocations in tight loops, as well as precision policies.
* **Cons:** You lose the object-oriented `num.Pow(exp)` ergonomics, and callers manage raw `*big.Float` instances.

---

### 2. Cross-Cutting Architectural Decisions

Regardless of which type pattern you choose, raising a `big.Float` to a `uint` involves three important Go and numerical design decisions:

#### 1. Mutability Style: `math/big` Convention vs. Immutable Values
* **Destination-Receiver (`z.Pow(x, exp)`):** Follows the `math/big` standard library pattern. The receiver `z` is mutated to hold the result and returned.
  * *Why choose it:* Allows callers to reuse memory buffers and drastically reduces garbage collection pressure in heavy computation.
* **Immutable / Functional (`res := x.Pow(exp)`):** `x` is never mutated; a newly allocated value is returned.
  * *Why choose it:* Safer against subtle aliasing bugs (e.g., if `z` and `x` point to the same underlying `big.Float` during intermediate squaring steps) and much easier to reason about in business/financial logic.

#### 2. Precision & Guard Bits (Why a second type or context often helps)
Unlike `big.Int.Exp`, floating-point exponentiation accumulates rounding error across $\approx \log_2(\text{exp})$ multiplications. 
* If a user wants a result accurate to `P` bits, computing intermediate squarings at precision `P` can cause drift in the least significant bits—especially for large `uint` exponents.
* Designing your custom type (Option B) or calculator context (Option D) to temporarily expand working precision (e.g., `workingPrec = targetPrec + guardBits`) during the squaring loop and round back to `targetPrec` at the end produces much cleaner, more predictable results.

#### 3. Never Copy Initialized `big.Float` Values
Internally, `big.Float` contains a slice (`[]big.Word`) for its mantissa. Copying a `big.Float` by value copies the slice header, meaning the copy and the original share the same backing array—mutating one can silently corrupt the other.
* **Design rule:** Always use **pointer receivers** (`*MyFloat`), return pointers from constructors, and avoid passing or returning your custom struct by value if it contains an initialized `big.Float`.

---

### Summary: Which Should You Pick?

| Scenario                                                     | Recommended Architecture                                     |
| :----------------------------------------------------------- | :----------------------------------------------------------- |
| You are building a rich domain type (e.g., `Money`, `Decimal`, `Scientific`) where `Pow` is one of several controlled operations. | **Option B (Encapsulated Wrapper)** with immutable methods.  |
| You mostly use `math/big.Float` everywhere, just need a `.Pow()` helper method alongside standard `.Add()`/`.Mul()` calls, and don't care about chaining. | **Option A (Struct Embedding)**.                             |
| You want maximum interoperability with existing `*big.Float` code, high performance (scratch buffer reuse), and controlled precision/guard bits. | **Option D (Context / Calculator Struct)** following the `z.Pow(x, exp)` destination pattern. |
| You want a multi-type separation of concerns (e.g., a `Float` value type + a `PowConfig` or `Engine` type). | Combine **Option A/B** for data with **Option D** for execution policy (e.g., `engine.Pow(base, exp)`). |



# Gemini 4 Hydrogen Question No 2 2026-11-14 16:08

I envisioned using Option B. Please outline the options for the internal calculation in func (d *Decimal) Pow(exp uint) *Decimal { ... }. 



# Gemini 4 Hydrogen Response No 2

With **Option B (Encapsulated Wrapper)** and the signature `func (d *Decimal) Pow(exp uint) *Decimal`, the idiomatic design treats `d` as **immutable**—allocating and returning a new `*Decimal` without mutating the receiver.

Before looking at the calculation strategies, every implementation should start with the same **edge-case fast paths** and end with a **clean rounding step**.

---

### 1. Pre-Loop Fast Paths (Common to All Options)

Handling special values up front avoids unnecessary work, prevents `0 * Inf` (`ErrNaN` panics), and protects against overflow on edge inputs:

* **`exp == 0`:** Return `1` (standard IEEE 754 convention, even for `0^0` or `Inf^0`).
* **`exp == 1`:** Return a copy of `d` (rounded to `d.prec` / `d.mode`).
* **`d` is `±0`:** Return `+0` (or `-0` if `d.val.Signbit() && exp%2 == 1`).
* **`d` is `±Inf`:** Return `+Inf` (or `-Inf` if `d.val.Signbit() && exp%2 == 1`).
* **`|d| == 1`:** Return `1` (or `-1` if negative and `exp` is odd). This avoids looping up to 64 times if a caller passes a huge exponent like `math.MaxUint`.

---

### 2. Options for the Internal Calculation

#### Option 1: Direct Binary Exponentiation (at Target Precision)
You run standard exponentiation by squaring directly at `d.prec`. Because `exp` is a `uint`, the loop executes **at most 64 iterations** even for `math.MaxUint64`.

```go
func (d *Decimal) Pow(exp uint) *Decimal {
    // ... fast paths for exp == 0, 1, ±0, ±Inf, ±1 ...

    res := new(big.Float).SetPrec(d.prec).SetMode(d.mode).SetInt64(1)
    base := new(big.Float).SetPrec(d.prec).SetMode(d.mode).Set(d.val)

    for exp > 0 {
        if exp&1 == 1 {
            res.Mul(res, base)
        }
        if exp > 1 {
            base.Mul(base, base)
            // Optional early exit if base overflows to Inf or underflows to 0
            if base.IsInf() || base.Sign() == 0 {
                return d.handleLimit(base, origExp)
            }
        }
        exp >>= 1
    }

    return &Decimal{val: res, prec: d.prec, mode: d.mode}
}
```
*Note: `big.Float.Mul` explicitly supports receiver aliasing (`res.Mul(res, base)` and `base.Mul(base, base)`), so you only allocate two `*big.Float` scratch buffers for the entire call.*

* **Pros:** Simplest implementation, lowest memory footprint, fastest execution.
* **Cons:** Rounding error accumulates across up to $\approx \log_2(\text{exp})$ multiplications. The least significant bits can drift by several ULPs (units in the last place) when `exp` is large.

---

#### Option 2: Binary Exponentiation with Guard Bits (Recommended Default)
You run the exact same binary exponentiation loop, but initialize `res` and `base` with **extended working precision**, then round down **once** at the very end.

```go
// Add guard bits proportional to the worst-case number of multiplications + safety margin
workPrec := d.prec + uint(bits.Len(exp)) + 32

res := new(big.Float).SetPrec(workPrec).SetMode(d.mode).SetInt64(1)
base := new(big.Float).SetPrec(workPrec).SetMode(d.mode).Set(d.val)

// ... same squaring loop ...

// Final single rounding step back to target precision:
// Use Set() (rounds to final's prec/mode), NOT Copy() (which would overwrite prec with workPrec).
final := new(big.Float).SetPrec(d.prec).SetMode(d.mode).Set(res)
return &Decimal{val: final, prec: d.prec, mode: d.mode}
```

* **Pros:** Virtually eliminates intermediate rounding drift with minimal extra CPU/memory cost (adding 32–96 bits is cheap in `math/big`). Cleanly leverages `d.mode` at the final step.
* **Cons:** Does not *mathematically guarantee* 100% exact rounding in rare pathological cases where the true result falls extraordinarily close to a rounding midpoint (the Table-Maker's Dilemma).

---

#### Option 3: Adaptive Guard Bits (Ziv's Strategy for Guaranteed Exact Rounding)
If your domain requires **strictly correctly rounded** results (identical to computing the power with infinite precision and rounding once):

1. Compute the result at `workPrec = d.prec + guardBits` (Option 2).
2. Track an upper bound on accumulated ULP error (roughly $\le 2 \times \text{numMuls}$ ULPs at `workPrec`).
3. Inspect the trailing `guardBits` of the result: if adding or subtracting the ULP error bound could cross a rounding boundary for `d.prec`, the result is ambiguous—double `guardBits` and rerun the loop.
4. Once unambiguous, round to `d.prec`.

* **Pros:** Mathematically guaranteed correct rounding to the last bit for any finite input and `exp`.
* **Cons:** Significantly more complex to write and verify; latency becomes non-deterministic on inputs that land near rounding ties.

---

#### Option 4: Exact Decomposition via `big.Int` (`MantExp` + `big.Int.Exp`)
Every finite non-zero `big.Float` is an exact integer mantissa scaled by a power of 2: $x = m \times 2^k$. You can extract $m$ into a `*big.Int`, raise $m^{\text{exp}}$ *exactly* using `big.Int.Exp`, and scale the binary exponent back:

1. Extract fractional mantissa & exponent: `e := d.val.MantExp(mant)` ($x = \text{mant} \times 2^e$, where $0.5 \le |\text{mant}| < 1$).
2. Shift to an exact integer: `p := mant.MinPrec()`, then `mant.SetMantExp(mant, int(p))`. Now `m, _ := mant.Int(nil)` is an exact integer such that $x = m \times 2^{e - p}$.
3. Compute $m^{\text{exp}}$ exactly with `new(big.Int).Exp(m, expBigInt, nil)`.
4. Convert back with `SetInt` and shift the exponent by `(e - int(p)) * int(exp)` (checking for `int64` overflow against `big.MaxExp` / `big.MinExp` so overflows become `Inf`/`0`).
5. Round once to `d.prec` using `d.mode`.

* **Pros:** **100% exact** intermediate computation with zero rounding drift. Often faster than `big.Float` loops for very small exponents (e.g., `exp <= 20`).
* **Cons (Major Hazard):** **Exponential memory explosion.** If `d` has 128 bits of precision and `exp` is `1_000_000`, $m^{\text{exp}}$ requires 128 million bits (~16 MB) just for one integer, and `exp = 1 << 30` will instantly OOM your process.

*(Note: Because `d.val` is already binary, you should avoid `Exp(n * Ln(x))`—transcendental approaches are slower, break on negative bases, and add unnecessary approximation error for integer `uint` exponents.)*

---

### Summary: Choosing the Right Strategy

| Strategy                                   | Complexity | Accuracy                       | Best Use Case                                                |
| :----------------------------------------- | :--------- | :----------------------------- | :----------------------------------------------------------- |
| **Option 1 (Direct Squaring)**             | Low        | Good (minor ULP drift)         | High-throughput math where the last few bits of precision aren't critical. |
| **Option 2 (Fixed/Heuristic Guard Bits)**  | Low        | Excellent                      | **Best general-purpose choice** for scientific, financial, or general `Decimal` libraries. |
| **Option 3 (Ziv's Adaptive Strategy)**     | High       | Guaranteed Exact               | Standards-compliant math libraries where IEEE-style exact rounding is a hard requirement. |
| **Option 4 (Hybrid `big.Int` + Option 2)** | Medium     | Guaranteed Exact (small `exp`) | Use Option 4 when `d.val.MinPrec() * exp` is small (e.g., `< 8192` bits), and fall back to Option 2 for larger exponents. |



# Gemini 4 Hydrogen Question No 3

'mantissa' as defined in Go's big.Float includes both the inter and fractional digits of the floating point number. Is that correct?

# Gemini 4 Hydrogen Response No 3

**Yes, conceptually that is correct**—with two important details about how `big.Float` normalizes and stores those digits.

### 1. Both integer and fractional parts live in the mantissa
In `math/big`, a nonzero finite `Float` is defined as:

$$\text{sign} \times \text{mantissa} \times 2^{\text{exponent}} \quad \text{where } 0.5 \le \text{mantissa} < 1.0$$

Because `big.Float` shifts the binary point all the way to the left of the most significant `1` bit (`0.5 <= mantissa < 1.0`), **all significant bits from both the integer part and the fractional part of your number are combined together in the mantissa** (up to `prec` bits).

**Concrete example (`6.25`):**
* In binary, `6.25` is `110.01₂` (integer part `110`, fractional part `.01`).
* `big.Float` normalizes this to:
  $$\text{mantissa} = 0.11001_2 \quad (= 0.78125_{10}), \qquad \text{exponent} = 3$$
* Check: $0.78125 \times 2^3 = 6.25$.
* Both the `110` (from the integer side) and the `01` (from the fractional side) are stored together in the mantissa.

### 2. Three nuances to keep in mind

1. **The normalized mantissa itself has `0` before the binary point (`< 1.0`):**
   Unlike IEEE-754 `float64` (which normalizes to `1.xxxx × 2^exp` and hides the leading `1` bit), `big.Float` normalizes to `0.1xxx × 2^exp`. There is **no hidden bit**—every significant bit, including the leading `1`, is explicitly stored in the mantissa's `prec` bits.
2. **Sign behavior differs slightly between the model and `MantExp`:**
   * In the type's mathematical definition, `sign` is separate and `0.5 <= mantissa < 1.0`.
   * When you call `x.MantExp(mant)`, however, `mant` inherits `x`'s sign, so `0.5 <= |mant| < 1.0`.
3. **Binary digits, not decimal digits:**
   The mantissa holds **base-2 bits**, not base-10 digits. So a number with exact decimal digits (like `6.1`) has an infinite repeating binary fractional part that gets rounded to fit into `prec` mantissa bits.

## Sources

1. [go.dev](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQF7PU8YhqebGAtnXfAwsQWAHt_-Dgq1ETT153-0i52mD9q8O8ArMtFa4OoX5qzEQkgrZA_ntNFNmSfdhbcw-tfP_dT7ntf0mpvUBqt4M84ux6Q%3D)
2. [go.dev](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQF7PU8YhqebGAtnXfAwsQWAHt_-Dgq1ETT153-0i52mD9q8O8ArMtFa4OoX5qzEQkgrZA_ntNFNmSfdhbcw-tfP_dT7ntf0mpvUBqt4M84ux6Q%3D)
3. [go.dev](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQF7PU8YhqebGAtnXfAwsQWAHt_-Dgq1ETT153-0i52mD9q8O8ArMtFa4OoX5qzEQkgrZA_ntNFNmSfdhbcw-tfP_dT7ntf0mpvUBqt4M84ux6Q%3D)
4. [go.dev](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQGt666ELjgh_aL47fl1OL9nfhL_OAsDmpMO3dQ04NAYB99WVSyL_mpod4ORIGmEG4T88a5q0fxsXtSrzSiMcka5lXDbx3cizD6b1YOXiLi6mBe470MZn6jryfPFYepxWmb1)



# Gemini 4 Hydrogen Question No 4 - Low Level Option 2

Going with **Option B (Encapsulated Wrapper)** with immutable methods., I would like to implement the following three strategies.

| **Option 2 (Fixed/Heuristic Guard Bits)** | Low    | Excellent                    | **Best general-purpose choice** for scientific, financial, or general Decimal libraries. |
| ----------------------------------------- | ------ | ---------------------------- | ------------------------------------------------------------ |
| **Option 3 (Ziv's Adaptive Strategy)**    | High   | Guaranteed Exact             | Standards-compliant math libraries where IEEE-style exact rounding is a hard requirement. |
| **Option 4 (Hybrid big.Int + Option 2)**  | Medium | Guaranteed Exact (small exp) | Use Option 4 when d.val.MinPrec() * exp is small (e.g., < 8192 bits)[[4](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQHswDv929zTuv2wMZtAonvBT2gKBFfLyrTU1LjI5tiSUOuEb3Al-Sxwrz98V2kzMTQha86RHBIDgklusz9m0EWi6kHLbmsdvzhyXaGisCWoOGNP1baIGPbU9gcnRhyhdfknzoJDpSgNZhvoWqY5trWEOFX3)], and fall back to Option 2 for larger exponents. |



The internal calculations of each operation should be implemented in a separate private method on a private custom type called by a public method on the  ***type Decimal struct***. (I wish to reuse these methods with other future components)

 Starting with **Option 2 (Fixed/Heuristic Guard Bits)**, please provide example code, with code comments, demonstrating this calculation.



# Gemini 4 Hydrogen Response No 4

To support your future components and the upcoming strategies (Options 3 and 4), the cleanest structure is to separate the **public encapsulated data type** (`Decimal`) from a **private computation engine** (`floatEngine`). 

Because `floatEngine` operates purely on `*big.Float` values using a configured precision and rounding mode, any other type you build later (e.g., `Vector`, `Matrix`, `Currency`) can reuse its methods directly without wrapping values in `Decimal`.

Below is a complete, commented implementation of **Option 2 (Fixed/Heuristic Guard Bits)**.

```go
package decimal

import (
	"math/big"
	"math/bits"
)

const defaultPrec uint = 128

// Decimal is an immutable, encapsulated wrapper around *big.Float (Option B).
type Decimal struct {
	val  *big.Float
	prec uint
	mode big.RoundingMode
}

// NewFromFloat64 constructs an immutable Decimal with explicit precision and rounding mode.
func NewFromFloat64(v float64, prec uint, mode big.RoundingMode) *Decimal {
	if prec == 0 {
		prec = defaultPrec
	}
	val := new(big.Float).SetPrec(prec).SetMode(mode).SetFloat64(v)
	return &Decimal{
		val:  val,
		prec: prec,
		mode: mode,
	}
}

// String formats the underlying value for display.
func (d *Decimal) String() string {
	if d == nil || d.val == nil {
		return "0"
	}
	return d.val.Text('g', -1)
}

// engine creates the private calculation engine matching this Decimal's settings.
func (d *Decimal) engine() floatEngine {
	if d == nil || d.prec == 0 {
		return floatEngine{prec: defaultPrec, mode: big.ToNearestEven}
	}
	return floatEngine{prec: d.prec, mode: d.mode}
}

// Pow returns a new Decimal equal to d^exp without mutating d.
func (d *Decimal) Pow(exp uint) *Decimal {
	eng := d.engine()

	var base *big.Float
	if d != nil && d.val != nil {
		base = d.val
	} else {
		base = new(big.Float).SetPrec(eng.prec).SetMode(eng.mode)
	}

	result := eng.powHeuristic(base, exp)

	return &Decimal{
		val:  result,
		prec: eng.prec,
		mode: eng.mode,
	}
}

// -----------------------------------------------------------------------------
// Private Reusable Calculation Engine
// -----------------------------------------------------------------------------

// floatEngine encapsulates precision/rounding policy and performs internal
// *big.Float calculations. Future strategies (Options 3 and 4) and other
// operations (Sqrt, Log, etc.) live here as private methods.
type floatEngine struct {
	prec uint
	mode big.RoundingMode
}

// powHeuristic computes x^exp using binary exponentiation (exponentiation by
// squaring) with heuristic guard bits, then rounds once to e.prec using e.mode.
func (e floatEngine) powHeuristic(x *big.Float, exp uint) *big.Float {
	// 1. Fast paths & IEEE 754 edge cases.
	switch exp {
	case 0:
		// x^0 == 1 for all x (including 0^0 and Inf^0).
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(1)
	case 1:
		// Use Set (not Copy) so the output is rounded to e.prec / e.mode.
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).Set(x)
	}

	// Preserve signed zero: (-0)^odd == -0, otherwise +0.
	if x.Sign() == 0 {
		z := new(big.Float).SetPrec(e.prec).SetMode(e.mode) // +0
		if x.Signbit() && exp&1 == 1 {
			z.Neg(z) // -0
		}
		return z
	}

	// Preserve signed infinity: (-Inf)^odd == -Inf, otherwise +Inf.
	if x.IsInf() {
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInf(x.Signbit() && exp&1 == 1)
	}

	// Handle |x| == 1 immediately so huge exponents (e.g., math.MaxUint)
	// do not trigger unnecessary loop iterations.
	if x.Cmp(big.NewFloat(1)) == 0 {
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(1)
	}
	if x.Cmp(big.NewFloat(-1)) == 0 {
		if exp&1 == 1 {
			return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(-1)
		}
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(1)
	}

	// 2. Calculate working precision with heuristic guard bits.
	//
	// Why bits.Len(exp) + 32?
	// The condition number of f(x) = x^exp is |exp|, meaning an intermediate
	// rounding error of 1 ULP can be amplified by a factor of up to exp
	// (roughly log2(exp) == bits.Len(exp) bits). Adding bits.Len(exp) bits
	// compensates for that amplification, and +32 provides margin for the
	// <= 64 loop multiplications and rounding boundaries.
	guardBits := uint(bits.Len(exp)) + 32
	workPrec := e.prec + guardBits

	// Guard against uint overflow if caller supplied a massive precision.
	if workPrec < e.prec {
		workPrec = e.prec
	}

	// 3. Binary exponentiation at extended working precision.
	// Keep intermediate rounding neutral (ToNearestEven) so directed modes on
	// e.mode (e.g., ToZero, ToPositiveInf) are applied once to the final value.
	res := new(big.Float).SetPrec(workPrec).SetMode(big.ToNearestEven).SetInt64(1)
	base := new(big.Float).SetPrec(workPrec).SetMode(big.ToNearestEven).Set(x)

	origExp := exp
	for exp > 0 {
		if exp&1 == 1 {
			res.Mul(res, base) // Aliasing receiver and operand is safe in math/big.
		}
		if exp > 1 {
			base.Mul(base, base)

			// Early exit on overflow or underflow: if base hits Inf or 0 during
			// squaring, x^origExp is guaranteed to hit Inf or 0 as well.
			if base.IsInf() {
				return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInf(x.Signbit() && origExp&1 == 1)
			}
			if base.Sign() == 0 {
				z := new(big.Float).SetPrec(e.prec).SetMode(e.mode) // +0
				if x.Signbit() && origExp&1 == 1 {
					z.Neg(z) // -0
				}
				return z
			}
		}
		exp >>= 1
	}

	// 4. Single final rounding pass down to target precision and rounding mode.
	// Note: Set() rounds res into final's e.prec/e.mode; Copy() would overwrite them.
	final := new(big.Float).SetPrec(e.prec).SetMode(e.mode).Set(res)
	return final
}
```

### Design Highlights to Note Before Adding Options 3 & 4
1. **Intermediate Rounding Mode (`big.ToNearestEven`):** Notice that `res` and `base` use `big.ToNearestEven` internally even if `d.mode` is a directed mode like `big.ToZero` (truncate) or `big.ToPositiveInf` (ceiling). Applying a directed rounding mode inside the loop causes systematic one-directional drift; keeping the loop unbiased and applying `e.mode` only in the final `Set(res)` gives much cleaner results.
2. **Extensibility for Options 3 & 4:** Adding the next two strategies requires zero changes to `Decimal`'s public API or external call sites—you will simply add `powZiv` (Option 3) and `powHybrid` (Option 4) as private methods on `floatEngine` (or add a strategy/policy field if `Decimal` lets callers choose). 

Whenever you are ready, we can implement **Option 3 (Ziv's Adaptive Strategy)** or **Option 4 (Hybrid `big.Int` + Option 2)** on `floatEngine`.



# Gemini 4 Hydrogen Question No 5 - Revised Code

I like the direction this is taking. At this point we need to make one modification. The public methods (public facing methods) need to incorporate the concept of 'targetFractionalDigits'. 'targetFractionalDigits' are defined as the number of fractional digits to the right of the decimal point. Further, we need to differentiate between 'targetFractionalDigits' and 'prec' (or precision as defined by big.Float). I think this means that we need to include a low level private custom type and methods which convert 'targetFractionalDigits' to 'prec' and 'prec' to 'targetFractionalDigits'.

Could you restate you last response with these modifications?

------

# Gemini 4 Hydrogen Response No 5

To cleanly separate **decimal fractional digits** (`targetFractionalDigits`) from **`big.Float` binary mantissa precision** (`prec`), the architecture now has three distinct layers:

1. **Public Layer (`Decimal`):** Speaks exclusively in terms of `targetFractionalDigits` (digits to the right of the decimal point), while also exposing `Prec()` and `ActualFractionalDigits()` so the distinction is transparent.
2. **Low-Level Private Conversion Layer (`precConverter`):** A dedicated private type responsible for converting `targetFractionalDigits` $\leftrightarrow$ `prec` using the number's binary exponent (`exp2` from `MantExp` / `Frexp`).
3. **Private Computation Layer (`floatEngine`):** Operates purely in `*big.Float` binary precision (`prec`), reusable across future types and strategies (Options 3 and 4).

### Why `exp2` is Required for Conversion
In `big.Float`, `prec` covers **both** the integer bits and the fractional bits of a number. For example, if a result grows from `12.34` to `286,138.1721`, the integer part expands from 4 bits to 19 bits. Keeping the same `prec` would starve the fractional digits. By passing the value's binary exponent `exp2` into `precConverter`, `Decimal.Pow` dynamically recalculates the `prec` needed to preserve `targetFractionalDigits` in the output.

---

### Updated Example Code (Option 2 + `precConverter`)

```go
package decimal

import (
	"math"
	"math/big"
	"math/bits"
)

const (
	defaultFractionalDigits uint = 6
	defaultSafetyBits       uint = 8       // Extra binary cushion to protect decimal rounding
	maxWorkingPrec          uint = 1 << 20 // Safety cap (~1M bits) to avoid runaway allocations
)

// Decimal is an immutable, encapsulated wrapper around *big.Float (Option B).
// Public methods operate in terms of targetFractionalDigits (base-10 digits to
// the right of the decimal point), while internally maintaining the required
// big.Float binary mantissa precision (prec).
type Decimal struct {
	val                    *big.Float
	targetFractionalDigits uint
	prec                   uint
	mode                   big.RoundingMode
}

// NewFromFloat64 constructs an immutable Decimal sized to preserve at least
// targetFractionalDigits decimal places for v's magnitude.
func NewFromFloat64(v float64, targetFractionalDigits uint, mode big.RoundingMode) *Decimal {
	conv := defaultConverter()

	var exp2 int
	if v != 0 && !math.IsInf(v, 0) && !math.IsNaN(v) {
		// Frexp uses the same normalization as big.Float.MantExp:
		// v == frac * 2^exp2 with 0.5 <= |frac| < 1.0
		_, exp2 = math.Frexp(v)
	}

	prec := conv.toPrec(targetFractionalDigits, exp2)
	val := new(big.Float).SetPrec(prec).SetMode(mode).SetFloat64(v)

	return &Decimal{
		val:                    val,
		targetFractionalDigits: targetFractionalDigits,
		prec:                   prec,
		mode:                   mode,
	}
}

// TargetFractionalDigits returns the requested base-10 digits after the decimal point.
func (d *Decimal) TargetFractionalDigits() uint {
	if d == nil {
		return 0
	}
	return d.targetFractionalDigits
}

// ActualFractionalDigits derives how many base-10 fractional digits are guaranteed
// by the underlying big.Float's binary precision at its current magnitude.
func (d *Decimal) ActualFractionalDigits() uint {
	if d == nil || d.val == nil || d.val.Sign() == 0 || d.val.IsInf() {
		return d.targetFractionalDigits
	}
	exp2 := d.val.MantExp(nil)
	return defaultConverter().toFractionalDigits(d.prec, exp2)
}

// Prec returns the underlying big.Float total mantissa precision in bits
// (covering both integer and fractional parts).
func (d *Decimal) Prec() uint {
	if d == nil {
		return 0
	}
	return d.prec
}

// WithFractionalDigits returns a new Decimal rescaled to targetFractionalDigits.
func (d *Decimal) WithFractionalDigits(targetFractionalDigits uint) *Decimal {
	if d == nil || d.val == nil {
		return NewFromFloat64(0, targetFractionalDigits, big.ToNearestEven)
	}
	var exp2 int
	if d.val.Sign() != 0 && !d.val.IsInf() {
		exp2 = d.val.MantExp(nil)
	}
	newPrec := defaultConverter().toPrec(targetFractionalDigits, exp2)
	val := new(big.Float).SetPrec(newPrec).SetMode(d.mode).Set(d.val)

	return &Decimal{
		val:                    val,
		targetFractionalDigits: targetFractionalDigits,
		prec:                   newPrec,
		mode:                   d.mode,
	}
}

// Pow returns a new Decimal equal to d^exp without mutating d, dynamically
// sizing the underlying big.Float precision so the result still preserves
// d.targetFractionalDigits decimal places after the decimal point.
func (d *Decimal) Pow(exp uint) *Decimal {
	if d == nil || d.val == nil {
		return NewFromFloat64(0, defaultFractionalDigits, big.ToNearestEven).Pow(exp)
	}

	conv := defaultConverter()
	eng := floatEngine{prec: d.prec, mode: d.mode}

	// Estimate the binary exponent of d^exp so we allocate enough mantissa bits
	// for BOTH the result's integer part and d.targetFractionalDigits.
	resultExp2 := eng.estimatePowExp2(d.val, exp)
	resultPrec := conv.toPrec(d.targetFractionalDigits, resultExp2)

	resultEng := floatEngine{prec: resultPrec, mode: d.mode}
	resultVal := resultEng.powHeuristic(d.val, exp)

	return &Decimal{
		val:                    resultVal,
		targetFractionalDigits: d.targetFractionalDigits,
		prec:                   resultPrec,
		mode:                   d.mode,
	}
}

// String formats the Decimal in fixed-point notation to targetFractionalDigits.
func (d *Decimal) String() string {
	if d == nil || d.val == nil {
		return "0"
	}
	if d.val.IsInf() {
		return d.val.String()
	}
	return d.val.Text('f', int(d.targetFractionalDigits))
}

// -----------------------------------------------------------------------------
// Low-Level Private Precision Converter
// -----------------------------------------------------------------------------

const (
	log2Of10 = math.Ln10 / math.Ln2 // ~3.321928 bits needed per decimal digit
	log10Of2 = math.Ln2 / math.Ln10 // ~0.301030 decimal digits per bit
)

// precConverter translates between base-10 targetFractionalDigits and
// big.Float binary mantissa precision (prec).
type precConverter struct {
	safetyBits uint
	maxPrec    uint
}

func defaultConverter() precConverter {
	return precConverter{
		safetyBits: defaultSafetyBits,
		maxPrec:    maxWorkingPrec,
	}
}

// toPrec converts targetFractionalDigits and a binary exponent exp2 (where
// 0.5*2^exp2 <= |x| < 2^exp2) into the total big.Float mantissa bits required.
func (c precConverter) toPrec(targetFractionalDigits uint, exp2 int) uint {
	// Bits needed to represent targetFractionalDigits to the right of the decimal point.
	fracBits := uint(math.Ceil(float64(targetFractionalDigits) * log2Of10))

	// If |x| >= 1 (exp2 > 0), the integer part consumes exp2 bits of the mantissa.
	// If |x| < 1 (exp2 <= 0), there are 0 integer bits.
	var intBits uint
	if exp2 > 0 {
		intBits = uint(exp2)
	}

	// Combine integer bits + fractional bits + safety margin, guarding overflow.
	total := intBits + fracBits
	if total < intBits { // uint overflow check
		return c.maxPrec
	}
	if total+c.safetyBits < total {
		return c.maxPrec
	}
	total += c.safetyBits

	// big.Float.SetPrec(0) maps finite values to 0, so always ensure prec >= 1.
	if total == 0 {
		total = 1
	}
	if c.maxPrec > 0 && total > c.maxPrec {
		return c.maxPrec
	}
	return total
}

// toFractionalDigits converts a big.Float binary precision (prec) and binary
// exponent (exp2) back into the number of guaranteed base-10 fractional digits.
func (c precConverter) toFractionalDigits(prec uint, exp2 int) uint {
	if prec <= c.safetyBits {
		return 0
	}
	usablePrec := prec - c.safetyBits

	var fracBits uint
	if exp2 > 0 {
		// Integer bits consume the top exp2 bits of the mantissa.
		intBits := uint(exp2)
		if usablePrec <= intBits {
			return 0
		}
		fracBits = usablePrec - intBits
	} else {
		// When |x| < 1 (exp2 <= 0), all usablePrec bits are fractional, plus
		// there are (-exp2) leading zero bits right of the binary point.
		fracBits = usablePrec + uint(-exp2)
	}

	return uint(math.Floor(float64(fracBits) * log10Of2))
}

// -----------------------------------------------------------------------------
// Private Reusable Calculation Engine
// -----------------------------------------------------------------------------

// floatEngine encapsulates binary precision/rounding policy and performs
// internal *big.Float calculations. Future strategies (Options 3 and 4) and
// other operations will live here as private methods.
type floatEngine struct {
	prec uint
	mode big.RoundingMode
}

// estimatePowExp2 estimates the binary exponent of x^exp in O(1) time
// using log2(|x^exp|) = exp * (log2(mantissa) + exp2).
func (e floatEngine) estimatePowExp2(x *big.Float, exp uint) int {
	if exp == 0 || x.Sign() == 0 || x.IsInf() {
		return 1 // 1, 0, or Inf requires no dynamic integer bit growth
	}

	mant := new(big.Float)
	exp2 := x.MantExp(mant)
	mant.Abs(mant)

	// 0.5 <= m64 <= 1.0, so log2(m64) is always finite and well-behaved.
	m64, _ := mant.Float64()
	log2X := math.Log2(m64) + float64(exp2)
	est := float64(exp) * log2X

	if est <= 0 {
		return 0 // Result magnitude < 1 (0 integer bits)
	}
	if est >= float64(big.MaxExp) {
		return big.MaxExp // Will overflow to Inf during powHeuristic
	}
	return int(math.Ceil(est)) + 1
}

// powHeuristic computes x^exp using binary exponentiation (Option 2:
// Fixed/Heuristic Guard Bits), then rounds once to e.prec using e.mode.
func (e floatEngine) powHeuristic(x *big.Float, exp uint) *big.Float {
	// 1. Fast paths & IEEE 754 edge cases.
	switch exp {
	case 0:
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(1)
	case 1:
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).Set(x)
	}

	if x.Sign() == 0 {
		z := new(big.Float).SetPrec(e.prec).SetMode(e.mode) // +0
		if x.Signbit() && exp&1 == 1 {
			z.Neg(z) // -0
		}
		return z
	}

	if x.IsInf() {
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInf(x.Signbit() && exp&1 == 1)
	}

	if x.Cmp(big.NewFloat(1)) == 0 {
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(1)
	}
	if x.Cmp(big.NewFloat(-1)) == 0 {
		if exp&1 == 1 {
			return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(-1)
		}
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(1)
	}

	// 2. Calculate working precision with heuristic guard bits.
	// bits.Len(exp) compensates for condition-number error amplification (~exp),
	// and +32 provides margin for loop multiplications and final rounding.
	guardBits := uint(bits.Len(exp)) + 32
	workPrec := e.prec + guardBits
	if workPrec < e.prec {
		workPrec = e.prec
	}

	// 3. Binary exponentiation at extended working precision with neutral rounding.
	res := new(big.Float).SetPrec(workPrec).SetMode(big.ToNearestEven).SetInt64(1)
	base := new(big.Float).SetPrec(workPrec).SetMode(big.ToNearestEven).Set(x)

	origExp := exp
	for exp > 0 {
		if exp&1 == 1 {
			res.Mul(res, base)
		}
		if exp > 1 {
			base.Mul(base, base)

			if base.IsInf() {
				return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInf(x.Signbit() && origExp&1 == 1)
			}
			if base.Sign() == 0 {
				z := new(big.Float).SetPrec(e.prec).SetMode(e.mode) // +0
				if x.Signbit() && origExp&1 == 1 {
					z.Neg(z) // -0
				}
				return z
			}
		}
		exp >>= 1
	}

	// 4. Single final rounding pass down to e.prec and e.mode.
	return new(big.Float).SetPrec(e.prec).SetMode(e.mode).Set(res)
}
```

### How the Pieces Work Together
* **`NewFromFloat64(12.34, 4, big.ToNearestEven)`:** Uses `precConverter.toPrec(4, 4)` $\rightarrow$ `4 (int) + 14 (frac) + 8 (safety) = 26` bits.
* **`d.Pow(5)`:** `estimatePowExp2` predicts the result `286138.1721` needs ~20 integer bits, so `precConverter.toPrec(4, 20)` upgrades `resultPrec` to `42` bits before `floatEngine.powHeuristic` runs with guard bits.
* **`d.ActualFractionalDigits()`:** Calls `precConverter.toFractionalDigits` to verify how many decimal fractional digits the current `prec` guarantees given the number's magnitude.

Whenever you're ready, we can add **Option 3 (Ziv's Adaptive Strategy)** or **Option 4 (Hybrid `big.Int` + Option 2)** onto `floatEngine`.

# Gemini 4 Hydrogen Question No 6  **Option 3 (Ziv's Adaptive Strategy)**

Yes, please add Option 3 (Ziv's Adaptive Strategy). Following the same private custom type and method pattern, please configure internal calculations in separate private methods on custom types.



# Gemini 4 Hydrogen Response No 6

Here is the updated implementation with **Option 3 (Ziv's Adaptive Strategy)** integrated alongside Option 2.

Following your architecture, each responsibility lives on its own private custom type with dedicated private methods:
* **`precConverter`:** Converts between `targetFractionalDigits` and `big.Float` `prec`.
* **`floatEngine`:** Shared `*big.Float` calculation kernel (`powSpecialCase`, `powHeuristic`, `powZiv`, `estimatePowExp2`).
* **`zivBounder`:** A new private custom type dedicated to tracking worst-case relative error in the squaring loop and verifying rounding unambiguity via an interval squeeze check.

### How `zivBounder` Guarantees Exact Rounding
Instead of relying on tricky bit-masking around rounding boundaries, `zivBounder.tryRound` uses a **monotonic interval squeeze check**:
1. While squaring in `powZiv`, `zivBounder.propagateMulErr` tracks a strict upper bound on relative rounding error (staying `0` while `Acc() == big.Exact`, which immediately handles exact midpoints like $1.5^2 = 2.25$ without looping).
2. From `res` and the error bound $\Delta$, it constructs a guaranteed bracket $[lo, hi]$ containing the true mathematical value $x^{\text{exp}}$ ($lo$ rounded down with `ToNegativeInf`, $hi$ rounded up with `ToPositiveInf`).
3. It rounds both $lo$ and $hi$ to `e.prec` using `e.mode`. Because all Go rounding modes are monotonic, if $\text{round}(lo) == \text{round}(hi)$, the result is **provably identical** to infinite-precision rounding. Otherwise, `powZiv` doubles the guard bits and retries.

---

### Complete Example Code (Option 2 + Option 3)

```go
package decimal

import (
	"math"
	"math/big"
	"math/bits"
)

const (
	defaultFractionalDigits uint = 6
	defaultSafetyBits       uint = 8       // Extra binary cushion for decimal rounding
	maxWorkingPrec          uint = 1 << 20 // Safety cap (~1M bits) to prevent runaway allocations
)

// Decimal is an immutable, encapsulated wrapper around *big.Float (Option B).
type Decimal struct {
	val                    *big.Float
	targetFractionalDigits uint
	prec                   uint
	mode                   big.RoundingMode
}

// NewFromFloat64 constructs an immutable Decimal sized to preserve at least
// targetFractionalDigits decimal places for v's magnitude.
func NewFromFloat64(v float64, targetFractionalDigits uint, mode big.RoundingMode) *Decimal {
	conv := defaultConverter()

	var exp2 int
	if v != 0 && !math.IsInf(v, 0) && !math.IsNaN(v) {
		_, exp2 = math.Frexp(v)
	}

	prec := conv.toPrec(targetFractionalDigits, exp2)
	val := new(big.Float).SetPrec(prec).SetMode(mode).SetFloat64(v)

	return &Decimal{
		val:                    val,
		targetFractionalDigits: targetFractionalDigits,
		prec:                   prec,
		mode:                   mode,
	}
}

// TargetFractionalDigits returns the requested base-10 digits after the decimal point.
func (d *Decimal) TargetFractionalDigits() uint {
	if d == nil {
		return 0
	}
	return d.targetFractionalDigits
}

// ActualFractionalDigits derives how many base-10 fractional digits are guaranteed
// by the underlying big.Float's binary precision at its current magnitude.
func (d *Decimal) ActualFractionalDigits() uint {
	if d == nil || d.val == nil || d.val.Sign() == 0 || d.val.IsInf() {
		return d.targetFractionalDigits
	}
	return defaultConverter().toFractionalDigits(d.prec, d.val.MantExp(nil))
}

// Prec returns the underlying big.Float total mantissa precision in bits.
func (d *Decimal) Prec() uint {
	if d == nil {
		return 0
	}
	return d.prec
}

// Pow returns d^exp using Option 2 (Fixed/Heuristic Guard Bits).
func (d *Decimal) Pow(exp uint) *Decimal {
	return d.powWithStrategy(exp, false)
}

// PowExact returns d^exp using Option 3 (Ziv's Adaptive Strategy) for
// guaranteed correctly rounded output at the computed target precision.
func (d *Decimal) PowExact(exp uint) *Decimal {
	return d.powWithStrategy(exp, true)
}

func (d *Decimal) powWithStrategy(exp uint, exact bool) *Decimal {
	if d == nil || d.val == nil {
		return NewFromFloat64(0, defaultFractionalDigits, big.ToNearestEven).powWithStrategy(exp, exact)
	}

	conv := defaultConverter()
	baseEng := floatEngine{prec: d.prec, mode: d.mode, maxPrec: maxWorkingPrec}

	resultExp2 := baseEng.estimatePowExp2(d.val, exp)
	resultPrec := conv.toPrec(d.targetFractionalDigits, resultExp2)
	resultEng := floatEngine{prec: resultPrec, mode: d.mode, maxPrec: maxWorkingPrec}

	var resultVal *big.Float
	if exact {
		resultVal = resultEng.powZiv(d.val, exp)
	} else {
		resultVal = resultEng.powHeuristic(d.val, exp)
	}

	return &Decimal{
		val:                    resultVal,
		targetFractionalDigits: d.targetFractionalDigits,
		prec:                   resultPrec,
		mode:                   d.mode,
	}
}

// String formats the Decimal in fixed-point notation to targetFractionalDigits.
func (d *Decimal) String() string {
	if d == nil || d.val == nil {
		return "0"
	}
	if d.val.IsInf() {
		return d.val.String()
	}
	return d.val.Text('f', int(d.targetFractionalDigits))
}

// -----------------------------------------------------------------------------
// Low-Level Private Precision Converter
// -----------------------------------------------------------------------------

const (
	log2Of10 = math.Ln10 / math.Ln2
	log10Of2 = math.Ln2 / math.Ln10
)

type precConverter struct {
	safetyBits uint
	maxPrec    uint
}

func defaultConverter() precConverter {
	return precConverter{
		safetyBits: defaultSafetyBits,
		maxPrec:    maxWorkingPrec,
	}
}

func (c precConverter) toPrec(targetFractionalDigits uint, exp2 int) uint {
	fracBits := uint(math.Ceil(float64(targetFractionalDigits) * log2Of10))

	var intBits uint
	if exp2 > 0 {
		intBits = uint(exp2)
	}

	total := intBits + fracBits
	if total < intBits || total+c.safetyBits < total {
		return c.maxPrec
	}
	total += c.safetyBits

	if total == 0 {
		total = 1
	}
	if c.maxPrec > 0 && total > c.maxPrec {
		return c.maxPrec
	}
	return total
}

func (c precConverter) toFractionalDigits(prec uint, exp2 int) uint {
	if prec <= c.safetyBits {
		return 0
	}
	usablePrec := prec - c.safetyBits

	var fracBits uint
	if exp2 > 0 {
		intBits := uint(exp2)
		if usablePrec <= intBits {
			return 0
		}
		fracBits = usablePrec - intBits
	} else {
		fracBits = usablePrec + uint(-exp2)
	}

	return uint(math.Floor(float64(fracBits) * log10Of2))
}

// -----------------------------------------------------------------------------
// Private Reusable Calculation Engine
// -----------------------------------------------------------------------------

type floatEngine struct {
	prec    uint
	mode    big.RoundingMode
	maxPrec uint
}

func (e floatEngine) estimatePowExp2(x *big.Float, exp uint) int {
	if exp == 0 || x.Sign() == 0 || x.IsInf() {
		return 1
	}

	mant := new(big.Float)
	exp2 := x.MantExp(mant)
	mant.Abs(mant)

	m64, _ := mant.Float64()
	est := float64(exp) * (math.Log2(m64) + float64(exp2))

	if est <= 0 {
		return 0
	}
	if est >= float64(big.MaxExp) {
		return big.MaxExp
	}
	return int(math.Ceil(est)) + 1
}

// powSpecialCase handles exp in {0, 1}, ±0, ±Inf, and ±1 in one shared place.
func (e floatEngine) powSpecialCase(x *big.Float, exp uint) (*big.Float, bool) {
	switch exp {
	case 0:
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(1), true
	case 1:
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).Set(x), true
	}

	if x.Sign() == 0 {
		z := new(big.Float).SetPrec(e.prec).SetMode(e.mode)
		if x.Signbit() && exp&1 == 1 {
			z.Neg(z)
		}
		return z, true
	}

	if x.IsInf() {
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInf(x.Signbit() && exp&1 == 1), true
	}

	if x.Cmp(big.NewFloat(1)) == 0 {
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(1), true
	}
	if x.Cmp(big.NewFloat(-1)) == 0 {
		if exp&1 == 1 {
			return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(-1), true
		}
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(1), true
	}

	return nil, false
}

// powHeuristic computes x^exp using Option 2 (Fixed/Heuristic Guard Bits).
func (e floatEngine) powHeuristic(x *big.Float, exp uint) *big.Float {
	if out, ok := e.powSpecialCase(x, exp); ok {
		return out
	}

	guardBits := uint(bits.Len(exp)) + 32
	workPrec := e.prec + guardBits
	if workPrec < e.prec || (e.maxPrec > 0 && workPrec > e.maxPrec) {
		workPrec = e.maxPrec
	}

	res := new(big.Float).SetPrec(workPrec).SetMode(big.ToNearestEven).SetInt64(1)
	base := new(big.Float).SetPrec(workPrec).SetMode(big.ToNearestEven).Set(x)

	origExp := exp
	for exp > 0 {
		if exp&1 == 1 {
			res.Mul(res, base)
		}
		if exp > 1 {
			base.Mul(base, base)
			if base.IsInf() {
				return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInf(x.Signbit() && origExp&1 == 1)
			}
			if base.Sign() == 0 {
				z := new(big.Float).SetPrec(e.prec).SetMode(e.mode)
				if x.Signbit() && origExp&1 == 1 {
					z.Neg(z)
				}
				return z
			}
		}
		exp >>= 1
	}

	return new(big.Float).SetPrec(e.prec).SetMode(e.mode).Set(res)
}

// powZiv computes x^exp using Option 3 (Ziv's Adaptive Strategy), iteratively
// increasing working precision until the final rounding to e.prec is provably
// unambiguous.
func (e floatEngine) powZiv(x *big.Float, exp uint) *big.Float {
	if out, ok := e.powSpecialCase(x, exp); ok {
		return out
	}

	// Ensure basePrec covers both target precision and x's exact mantissa bits
	// so the initial load of base = x introduces zero quantization error.
	basePrec := e.prec
	if minP := x.MinPrec(); minP > basePrec {
		basePrec = minP
	}

	// Start with enough guard bits so almost all inputs succeed on Pass 1.
	guardBits := uint(bits.Len(exp)) + 32

	for {
		workPrec := basePrec + guardBits
		capped := false
		if workPrec < basePrec || (e.maxPrec > 0 && workPrec >= e.maxPrec) {
			workPrec = e.maxPrec
			capped = true
		}

		bounder := newZivBounder(e.prec, e.mode, workPrec)
		res, errUnits, earlyOut := e.runZivPass(x, exp, bounder)
		if earlyOut != nil {
			return earlyOut
		}

		// Verify whether rounding res to (e.prec, e.mode) is unambiguous.
		if rounded, ok := bounder.tryRound(res, errUnits); ok {
			return rounded
		}

		// Safeguard: stop if we hit the maximum precision cap.
		if capped {
			return new(big.Float).SetPrec(e.prec).SetMode(e.mode).Set(res)
		}

		// Double guard bits for the next Ziv refinement pass.
		guardBits *= 2
	}
}

// runZivPass executes one binary-exponentiation pass at b.workPrec while
// accumulating a deterministic upper bound on relative rounding error.
func (e floatEngine) runZivPass(x *big.Float, exp uint, b zivBounder) (*big.Float, float64, *big.Float) {
	res := new(big.Float).SetPrec(b.workPrec).SetMode(big.ToNearestEven).SetInt64(1)
	base := new(big.Float).SetPrec(b.workPrec).SetMode(big.ToNearestEven).Set(x)

	resErr := 0.0
	baseErr := 0.0
	if base.Acc() != big.Exact {
		// Only happens if x.MinPrec() exceeded maxWorkingPrec.
		baseErr = 1.0
	}

	origExp := exp
	for exp > 0 {
		if exp&1 == 1 {
			res.Mul(res, base)
			resErr = b.propagateMulErr(resErr, baseErr, res.Acc())
		}
		if exp > 1 {
			base.Mul(base, base)
			baseErr = b.propagateMulErr(baseErr, baseErr, base.Acc())

			if base.IsInf() {
				inf := new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInf(x.Signbit() && origExp&1 == 1)
				return nil, 0, inf
			}
			if base.Sign() == 0 {
				z := new(big.Float).SetPrec(e.prec).SetMode(e.mode)
				if x.Signbit() && origExp&1 == 1 {
					z.Neg(z)
				}
				return nil, 0, z
			}
		}
		exp >>= 1
	}

	return res, resErr, nil
}

// -----------------------------------------------------------------------------
// Private Ziv Error-Bounding & Interval Verification Type
// -----------------------------------------------------------------------------

// zivBounder tracks relative rounding error in units of u = 2^-workPrec and
// verifies whether rounding a working-precision result to (targetPrec, targetMode)
// is provably unambiguous.
type zivBounder struct {
	targetPrec   uint
	targetMode   big.RoundingMode
	workPrec     uint
	unitRoundoff float64 // 2^-workPrec (upper bound on relative error for ToNearestEven)
}

func newZivBounder(targetPrec uint, targetMode big.RoundingMode, workPrec uint) zivBounder {
	return zivBounder{
		targetPrec:   targetPrec,
		targetMode:   targetMode,
		workPrec:     workPrec,
		unitRoundoff: math.Ldexp(1.0, -int(workPrec)),
	}
}

// propagateMulErr computes the upper bound on relative error (in units of u)
// after multiplying two values with error bounds errA and errB.
func (z zivBounder) propagateMulErr(errA, errB float64, acc big.Accuracy) float64 {
	// (1 + errA*u)*(1 + errB*u) - 1 = (errA + errB + errA*errB*u)*u
	prop := errA + errB + errA*errB*z.unitRoundoff
	if acc == big.Exact {
		return prop
	}
	// Multiply by (1 + u) for the new rounding step.
	return prop + 1.0 + prop*z.unitRoundoff
}

// tryRound checks if the true mathematical value in [res - delta, res + delta]
// rounds identically at both endpoints under (targetPrec, targetMode).
func (z zivBounder) tryRound(res *big.Float, errUnits float64) (*big.Float, bool) {
	// Fast path: if every operation was Exact (e.g., exact powers or exact ties
	// like 1.5^2 = 2.25), res has zero error and rounds immediately.
	if errUnits == 0 {
		return new(big.Float).SetPrec(z.targetPrec).SetMode(z.targetMode).Set(res), true
	}

	if res.IsInf() {
		return new(big.Float).SetPrec(z.targetPrec).SetMode(z.targetMode).Set(res), true
	}

	// If accumulated error exceeds 0.25 ULP of targetPrec, do not bother with
	// the interval check; signal powZiv to increase guardBits immediately.
	var guardBits int
	if z.workPrec > z.targetPrec {
		guardBits = int(z.workPrec - z.targetPrec)
	}
	if errUnits >= math.Ldexp(0.25, guardBits) {
		return nil, false
	}

	// Apply a 1% conservative multiplier to errUnits to cover both float64
	// rounding and converting relative error |Y_true - res| / |Y_true| to |res|.
	safeErrUnits := errUnits * 1.01

	// Build eps >= safeErrUnits * 2^-workPrec rounded upward (ToPositiveInf).
	eps := new(big.Float).SetPrec(64).SetMode(big.ToPositiveInf).SetFloat64(safeErrUnits)
	eps.SetMantExp(eps, eps.MantExp(nil)-int(z.workPrec))

	// delta >= |res| * eps, rounded upward.
	delta := new(big.Float).SetPrec(64).SetMode(big.ToPositiveInf).Abs(res)
	delta.Mul(delta, eps)

	// Construct [lo, hi] strictly enclosing the true value Y_true.
	lo := new(big.Float).SetPrec(z.workPrec).SetMode(big.ToNegativeInf).Sub(res, delta)
	hi := new(big.Float).SetPrec(z.workPrec).SetMode(big.ToPositiveInf).Add(res, delta)

	// Round both endpoints to target precision and target rounding mode.
	rLo := new(big.Float).SetPrec(z.targetPrec).SetMode(z.targetMode).Set(lo)
	rHi := new(big.Float).SetPrec(z.targetPrec).SetMode(z.targetMode).Set(hi)

	// By monotonicity of rounding, rLo == rHi guarantees round(Y_true) == rLo.
	if rLo.Cmp(rHi) == 0 {
		return rLo, true
	}
	return nil, false
}
```

### Summary of the Custom Type Separation
* `Decimal` $\rightarrow$ Public immutable API (`targetFractionalDigits`-based).
* `precConverter` $\rightarrow$ Private `targetFractionalDigits` $\leftrightarrow$ `prec` converter.
* `floatEngine` $\rightarrow$ Private `*big.Float` operation runner (`powHeuristic`, `powZiv`).
* `zivBounder` $\rightarrow$ Private Ziv error propagator and monotonic interval verifier.

Whenever you're ready, we can add **Option 4 (Hybrid `big.Int` + Option 2)** onto `floatEngine` (with its own small private helper type for exact `MantExp` $\leftrightarrow$ `big.Int` decomposition).



# Gemini 4 Hydrogen Question No 7

Before we continue, could you provide no more than four of the best sources with narrative text explaining  **Ziv's Adaptive Strategy**?



# Gemini 4 Hydrogen Response No 7

Here are four of the best sources that provide clear narrative explanations of **Ziv's Adaptive Strategy** (also known as Ziv's "onion-peeling" strategy), moving from the original paper to modern arbitrary-precision implementations:

1. **Abraham Ziv (1991) — [Fast Evaluation of Elementary Mathematical Functions with Correctly Rounded Last Bit](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEqinpR_7_m-5KxalL4nwq7W14f6_exlqbq7E5c8S6Tl-YuQQ77IYXo3vAjWxLpqhubntn9xYA4l8FPjWQSJN5Kw6qQw6nj8Qhjvo4LVXNohN6ZN3uABt4RJ96kFx6JieHBYUsjonuy0rADWn1qihU_lmD0OpOh3FivJeJPq3UvxF-6LI4PHtvAH1oXvPtjdGBcWJVw2H2NlInLpkZ8RKJU8fWG7cgO02AAtvBFYq8izQ==)**  
   *ACM Transactions on Mathematical Software (TOMS), Vol. 17, No. 3, pp. 410–423 (DOI: `10.1145/114697.116813`)*.  
   The foundational paper that introduced the strategy. Its narrative explains the motivation and general methodology: computing a fast initial approximation, running a rounding test to check whether the error interval straddles a rounding boundary, and peeling into higher-precision fallback stages only on the rare inputs that fail the test.

2. **Laurent Fousse, Guillaume Hanrot, Vincent Lefèvre, Patrick Pélissier, and Paul Zimmermann (2007) — [MPFR: A Multiple-Precision Binary Floating-Point Library with Correct Rounding](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGSeAOaOjMmHEuB57sm48W_t4IUltsCkFol-T3fwSdzu0g0SRpFJ_sSSQRVfv8TH49-iho0muGRr8ekoZenABRcfFxjuhbYc6T4V4uJMJKRYlljR3_iqO6sxyYIqQjcyp9qHQakQfo-XNlcahUGRyL__YzEPOpPBtotnzqE4WZh-bDIXHjMqA6ZPJHVw7WYCOSF_BbwDwwZleqr0Ap4cycSNPoC1h-TM55NbmQyT-ixyMLEUW2hVBc=)**  
   *ACM Transactions on Mathematical Software (TOMS), Vol. 33, No. 2, Article 13 (DOI: `10.1145/1236463.1236468`)*.  
   Directly relevant to your `big.Float` design because it explains how Ziv's strategy works in an **arbitrary-precision** binary floating-point library rather than fixed 64-bit hardware floats. It walks through concrete examples (such as computing $\cos(1)$ to 42 bits starting at 63 working bits) and explains why average-case overhead stays small while rare worst cases require precision escalation.

3. **Jean-Michel Muller et al. (2018) — [Handbook of Floating-Point Arithmetic](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEvTzmt3QH1KI4YHEVjTJyf2STpr9gHz1Ivc7MwKJH-tq_pLNt6LcL1Y2EaNoj65CzY315kDXWNmsOQ9aG7WGHUNqEl-M8MuDWqiLAmfCCQZBdcORWqerXKM8reJHtRl8glaCY2n9-cyJfXjnq3fy49gU-V56IAcTJStkTxJfKXDJWamK7xXdtsKvipgco=) (2nd ed., Chapter 11)**  
   *Birkhäuser (DOI: `10.1007/978-3-319-76526-6`)*.  
   The standard textbook narrative on the subject. Section 11.6 (*"The Table Maker's Dilemma and Ziv's onion peeling strategy"*) gives a structured prose explanation of why fixed extra precision cannot guarantee correct rounding, how rounding tests are constructed, and how the accuracy/performance trade-off works in practice.

4. **The FLINT Documentation — [Correctly rounded elementary functions](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEL3ErF5lbL5al7hcVZLBS2_V1n5KQebleSZhVtLEsRYHNjqnLS9OY3tBbAwShTAQfJm19OE6aPbb0_jJ0kb10PKUmqCkMyp9Ic7fh30HVByD21arN8-gTUAE7KBy1a2OM7Szq-8BqozxyzPtx0hvy5GJrihLtIaKpbLvT2_S0=)**  
   *FLINT Floating-Point Arithmetic Algorithms Guide*.  
   A modern, freely accessible engineering narrative that closely mirrors the interval-enclosure approach. It describes how an implementation certifies an output enclosure, retries with exponentially increasing guard bits under Ziv's strategy, reasons about termination, and separates exact or tiny-argument edge cases so Ziv's loop does not spin indefinitely on exact ties or midpoints.

## Sources

1. [rutgers.edu](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQH8Xda5wL2kbKqIEZ3Pk9md6cVbuOEzTXXU1y2RB_yZRAZaUUnJApLNmuRJJaFD_ZeqxB7t9_avkP5cuBvP9x-YF1GI8K2w3FkpRAfYtDqUsIGBmQlpvD0eTFgsI5QCNscNHat_eZy8kIu9MiCIn9OVTbslqp42JNx6C_B5zDkC)
2. [flintlib.org](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQHOwhizgQA92wLDkPFqGJapxEQ1topug3vzREbIWydfA5Pseb-bO5hw-Gmi26XzyXAmgVzX-MarGzfuT5gKBiYi6U6ztJtJOXQ6QLYtNIBLqKJC0aeZA5aazFVkt5sreKXSoPDHUQ%3D%3D)
3. [rutgers.edu](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQG02JJ4K6MXwwbiHa50pkVUTdHhGkPp_iRb4tXFK0BTh6SczKM6FPJX5MnnAhyYoElVdjRW30wevR3ujsZ9U2jJpsb8yw9m27K3GJTyumJ-Jv1RzsWHBO3aWx3VyQoBp5P97B_Dt5N7t0o9lihR7cY_35LuKjy-9NUmXZXTdvaDLDJjJ-M2ONiuu6svQZripZ6-wPOqbvhWTto_8DJQoUC7xNpyUrK5nbLf2YLhXw%3D%3D)
4. [nih.gov](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQEWgtgILIpu_5weVYJzvxrOc4BA0y6Eh71tOmZF9WwfYKgBovkOz-z9ZR7ZAdacHYYKQtCLRWHwrcbXHtKahSKMekg_j0Q2OAbS8plDoMQr_bSsjKdW6NaY-XmPdGg7qhkiCI-hlvx9)
5. [wikipedia.org](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQFZpiZomwK152bfyOHLoRLXz8_7QpybhO1BxXwnGbOZdP7ZLKJKvUdLWyY2Qt9Bf7L3UBEoSq4CUj394mYZiMvYB0yXuK4YHHRAWyjZD75B42tvdF-nCmIKzxddMbA26KhAuAyJaxU5OVRJ)
6. [typeset.io](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQEekmRwBC2j4EpauQZqC4qOOPPK13cFqH5NdbwvLK75QpLdPx9pL55Xc4k0RNcCCKV-qVs50FvCuL7sr2up_CdRGtiBlUA0g2fFT95YiQ5irhcsZJDkRMItavZQjWckFndS3r6m2pZP)
7. [rutgers.edu](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQH8Xda5wL2kbKqIEZ3Pk9md6cVbuOEzTXXU1y2RB_yZRAZaUUnJApLNmuRJJaFD_ZeqxB7t9_avkP5cuBvP9x-YF1GI8K2w3FkpRAfYtDqUsIGBmQlpvD0eTFgsI5QCNscNHat_eZy8kIu9MiCIn9OVTbslqp42JNx6C_B5zDkC)
8. [researchgate.net](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQGSeAOaOjMmHEuB57sm48W_t4IUltsCkFol-T3fwSdzu0g0SRpFJ_sSSQRVfv8TH49-iho0muGRr8ekoZenABRcfFxjuhbYc6T4V4uJMJKRYlljR3_iqO6sxyYIqQjcyp9qHQakQfo-XNlcahUGRyL__YzEPOpPBtotnzqE4WZh-bDIXHjMqA6ZPJHVw7WYCOSF_BbwDwwZleqr0Ap4cycSNPoC1h-TM55NbmQyT-ixyMLEUW2hVBc%3D)
9. [researchgate.net](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQHHQ62ffUALQwqKJkPy3R6hBGKQwmBD8TQZzOeO60Mk9OdnSKG0geT1pmifyNoyiUoBK1m-IZkR43YAMfREHseBVAlKR56T6yoaX1p7rEvOm12kKioF2CIPDar3c5F_IFnpgVVZ7iJe6pP1Ss_p7CJYdv_aJSyJ1nflOPZyXMpeOcZW51o_Z3JbproaEyUwC0Vc5mWMAMl1MO7lYBu7JpVvMVsjrRKkllfEJviTM-q4q1X-X2A8ifc%3D)
10. [researchgate.net](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQGSeAOaOjMmHEuB57sm48W_t4IUltsCkFol-T3fwSdzu0g0SRpFJ_sSSQRVfv8TH49-iho0muGRr8ekoZenABRcfFxjuhbYc6T4V4uJMJKRYlljR3_iqO6sxyYIqQjcyp9qHQakQfo-XNlcahUGRyL__YzEPOpPBtotnzqE4WZh-bDIXHjMqA6ZPJHVw7WYCOSF_BbwDwwZleqr0Ap4cycSNPoC1h-TM55NbmQyT-ixyMLEUW2hVBc%3D)
11. [researchgate.net](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQGSeAOaOjMmHEuB57sm48W_t4IUltsCkFol-T3fwSdzu0g0SRpFJ_sSSQRVfv8TH49-iho0muGRr8ekoZenABRcfFxjuhbYc6T4V4uJMJKRYlljR3_iqO6sxyYIqQjcyp9qHQakQfo-XNlcahUGRyL__YzEPOpPBtotnzqE4WZh-bDIXHjMqA6ZPJHVw7WYCOSF_BbwDwwZleqr0Ap4cycSNPoC1h-TM55NbmQyT-ixyMLEUW2hVBc%3D)
12. [arxiv.org](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQH9lt2cIRYtrklCNdglFr3u7GDNDuJEDj5DKwAUurJ1wDFJJMTJOLiblrFM337Tdmb8danUfDmc4lNlz4q6ZzGaZw_ROm17BmtCQ6L7GPVCXONAsH423SWh)
13. [flintlib.org](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQEL3ErF5lbL5al7hcVZLBS2_V1n5KQebleSZhVtLEsRYHNjqnLS9OY3tBbAwShTAQfJm19OE6aPbb0_jJ0kb10PKUmqCkMyp9Ic7fh30HVByD21arN8-gTUAE7KBy1a2OM7Szq-8BqozxyzPtx0hvy5GJrihLtIaKpbLvT2_S0%3D)
14. [dandelon.com](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQEXyYuB8tNDoIg3Up7cB5gWUU1QaX1BbnPPw8KagZ16WuUZZH73PtCbeNf9gXSbvZXsLkKveqYWRljQMOE4Su_jpRk1_sq9EhzKjO8eFr_QSpkGvm_K9OJO4tnrHt8cdJm_4adDatu19uIWnAAgGUrmP-OAGI6wbL6Te-yT1GWl57xOlMH78Nu3y8rzQeRm2PAHT66mDNJCh5_sU1pYYFw%3D)
15. [flintlib.org](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQHQhVOBR5-s6M3aml0RzxIepjmA09bAI9CQ_oaISwf4N4bM5S4FiUTQzQ2wglj1NJFbU2CM3FCdPDYH6qAe7CZhcAkv1cKt5719SFZaMIdKBC-l7xqqz2gfO8FatxO26rDw-WSBczIHs4SKJAfml4k6EYUH)
16. [flintlib.org](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQEL3ErF5lbL5al7hcVZLBS2_V1n5KQebleSZhVtLEsRYHNjqnLS9OY3tBbAwShTAQfJm19OE6aPbb0_jJ0kb10PKUmqCkMyp9Ic7fh30HVByD21arN8-gTUAE7KBy1a2OM7Szq-8BqozxyzPtx0hvy5GJrihLtIaKpbLvT2_S0%3D)
17. [flintlib.org](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQHkdNYmWJ7_w3HtT0EhP4HoJScUrX2EdB8YEM_yuHRORDPRnS3_QyVbmI0WqDcvjkuYKZ9SO26vKJQxq072sfF0BWZKcqlTyYVMqlPiVzKzoYkB2AawJiLuEHd0o-ivCPFlwcOukoNgDlvpJQ%3D%3D)
18. [flintlib.org](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQGuWKhNjgpwSJZMygZEJLIlEPg7Eg_deiNav2duupUdRVT1b2u1rC7W3VzaOxzLhGNC7YGBnaO2zpYBgKm3JkxoDsitmahhJ41bTbOd823C_PXh3E2IC3A%3D)

https://www.google.com/search?q=%22Handbook%20of%20Floating-Point%20Arithmetic%22%20%22Ziv%27s%20strategy%22%20OR%20%22onion-peeling%22&client=app-vertex-grounding&safesearch=active

https://www.google.com/search?q=%22Fast%20evaluation%20of%20elementary%20mathematical%20functions%20with%20correctly%20rounded%20last%20bit%22%20Ziv%201991&client=app-vertex-grounding&safesearch=active

https://www.google.com/search?q=%22Ziv%27s%20strategy%22%20%22Table-Maker%27s%20Dilemma%22%20Muller%20Lefevre&safesearch=active

https://www.google.com/search?q=%22Ziv%27s%20strategy%22%20OR%20%22Ziv%27s%20onion-peeling%22%20floating%20point%20correct%20rounding%20MPFR&safesearch=active

https://www.google.com/search?q=%22Ziv%27s%20strategy%22%20site:flintlib.org%20OR%20site:mpfr.org&safesearch=active

https://www.google.com/search?q=%22MPFR:%20A%20multiple-precision%20binary%20floating-point%20library%20with%20correct%20rounding%22%20%22Ziv%27s%20strategy%22&safesearch=active



# Gemini 4 Hydrogen Question No 8 - Do Option 4

Please add **Option 4 (Hybrid big.Int + Option 2)** onto floatEngine. Remember to include safety mechanisms which trigger if the exponent is too large or too small: 

------

**Option 4 (Hybrid big.Int + Option 2)**  Medium  Guaranteed Exact (small exp)  Use Option 4 when d.val.MinPrec() * exp is small (e.g., < 8192 bits)[[4](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQHswDv929zTuv2wMZtAonvBT2gKBFfLyrTU1LjI5tiSUOuEb3Al-Sxwrz98V2kzMTQha86RHBIDgklusz9m0EWi6kHLbmsdvzhyXaGisCWoOGNP1baIGPbU9gcnRhyhdfknzoJDpSgNZhvoWqY5trWEOFX3)], and fall back to Option 2 for larger exponents.

------



# Gemini 4 Hydrogen Response No 8

To add **Option 4 (Hybrid `big.Int` + Option 2)** into your architecture, we follow the same pattern:

* **`Decimal` (Public):** Gains `PowHybrid(exp uint) *Decimal`, wired through a clean internal `powStrategy` enum alongside `Pow` (Option 2) and `PowExact` (Option 3).
* **`floatEngine` (Private):** Gains `powHybrid(x *big.Float, exp uint) *big.Float`, which shares `powSpecialCase` and falls back to `powHeuristic` (Option 2).
* **`dyadicPowKernel` (New Private Custom Type):** Encapsulates the `big.Float` $\leftrightarrow$ `big.Int` dyadic decomposition ($x = m \times 2^k$), threshold gating, and exponent overflow/underflow safety checks.

### Safety Mechanisms Included for Too-Large / Too-Small Exponents
1. **Power exponent `exp` too small (`0` or `1`):** Short-circuited in `powSpecialCase` before any decomposition runs.
2. **Power exponent `exp` too large for `big.Int` (`x.MinPrec() * exp > maxExactBits`):** Checked via division (`exp > maxExactBits/minPrec`) so the multiplication itself can never overflow `uint`, triggering an immediate fallback to **Option 2 (`powHeuristic`)**.
3. **Scaled binary exponent too large (`> big.MaxExp`) or too small (`< big.MinExp`):** Checked both *before* `big.Int.Exp` (using `x.MantExp` bounds) and *after* `big.Int.Exp` (using `int64` arithmetic so `int32` cannot wrap on 32-bit systems), returning properly signed `±Inf` or `±0`.

---

### Option 4 Implementation (Drop-in Addition to `package decimal`)

Update `powWithStrategy` on `Decimal` and add `powHybrid` + `dyadicPowKernel` below `floatEngine`:

```go
const defaultMaxExactBits uint = 8192

type powStrategy uint8

const (
	strategyHeuristic powStrategy = iota // Option 2: Fixed/Heuristic Guard Bits
	strategyZiv                          // Option 3: Ziv's Adaptive Strategy
	strategyHybrid                       // Option 4: Hybrid big.Int + Option 2
)

// Pow returns d^exp using Option 2 (Fixed/Heuristic Guard Bits).
func (d *Decimal) Pow(exp uint) *Decimal {
	return d.powWithStrategy(exp, strategyHeuristic)
}

// PowExact returns d^exp using Option 3 (Ziv's Adaptive Strategy).
func (d *Decimal) PowExact(exp uint) *Decimal {
	return d.powWithStrategy(exp, strategyZiv)
}

// PowHybrid returns d^exp using Option 4 (Exact big.Int for small
// MinPrec()*exp, falling back to Option 2 for larger exponents).
func (d *Decimal) PowHybrid(exp uint) *Decimal {
	return d.powWithStrategy(exp, strategyHybrid)
}

func (d *Decimal) powWithStrategy(exp uint, strategy powStrategy) *Decimal {
	if d == nil || d.val == nil {
		return NewFromFloat64(0, defaultFractionalDigits, big.ToNearestEven).powWithStrategy(exp, strategy)
	}

	conv := defaultConverter()
	baseEng := floatEngine{prec: d.prec, mode: d.mode, maxPrec: maxWorkingPrec}

	resultExp2 := baseEng.estimatePowExp2(d.val, exp)
	resultPrec := conv.toPrec(d.targetFractionalDigits, resultExp2)
	resultEng := floatEngine{prec: resultPrec, mode: d.mode, maxPrec: maxWorkingPrec}

	var resultVal *big.Float
	switch strategy {
	case strategyZiv:
		resultVal = resultEng.powZiv(d.val, exp)
	case strategyHybrid:
		resultVal = resultEng.powHybrid(d.val, exp)
	default:
		resultVal = resultEng.powHeuristic(d.val, exp)
	}

	return &Decimal{
		val:                    resultVal,
		targetFractionalDigits: d.targetFractionalDigits,
		prec:                   resultPrec,
		mode:                   d.mode,
	}
}

// -----------------------------------------------------------------------------
// Addition to floatEngine
// -----------------------------------------------------------------------------

// powHybrid computes x^exp using exact big.Int exponentiation when
// x.MinPrec() * exp <= defaultMaxExactBits, safely falling back to Option 2
// (powHeuristic) when the exponent or bit-growth is too large.
func (e floatEngine) powHybrid(x *big.Float, exp uint) *big.Float {
	// 1. Handle exp in {0, 1}, ±0, ±Inf, and ±1 immediately.
	if out, ok := e.powSpecialCase(x, exp); ok {
		return out
	}

	// 2. Attempt exact big.Int calculation if within the bit-budget threshold.
	kernel := dyadicPowKernel{
		targetPrec:   e.prec,
		targetMode:   e.mode,
		maxExactBits: defaultMaxExactBits,
	}

	if kernel.canUseExact(x, exp) {
		if out, ok := kernel.powExact(x, exp); ok {
			return out
		}
	}

	// 3. Fallback to Option 2 (Fixed/Heuristic Guard Bits) for large exponents
	// or if exact decomposition was declined.
	return e.powHeuristic(x, exp)
}

// -----------------------------------------------------------------------------
// Private Exact Dyadic / big.Int Exponentiation Kernel
// -----------------------------------------------------------------------------

// dyadicPowKernel decomposes a finite big.Float into m * 2^scaleExp2, computes
// m^exp exactly via big.Int.Exp, rescales the binary exponent, and rounds once
// to (targetPrec, targetMode).
type dyadicPowKernel struct {
	targetPrec   uint
	targetMode   big.RoundingMode
	maxExactBits uint
}

// canUseExact reports whether x.MinPrec() * exp fits within maxExactBits
// without risking uint multiplication overflow.
func (k dyadicPowKernel) canUseExact(x *big.Float, exp uint) bool {
	if exp <= 1 || k.maxExactBits == 0 {
		return false
	}
	minPrec := x.MinPrec()
	if minPrec == 0 || minPrec > k.maxExactBits {
		return false
	}
	// Overflow-safe check for (minPrec * exp <= k.maxExactBits):
	return exp <= k.maxExactBits/minPrec
}

// powExact executes exact big.Int exponentiation with pre- and post-exponent
// overflow/underflow guards. Returns (nil, false) if it cannot safely proceed.
func (k dyadicPowKernel) powExact(x *big.Float, exp uint) (*big.Float, bool) {
	if uint64(exp) > math.MaxInt64 {
		return nil, false
	}

	// Safety check 1 (Pre-Exp): detect guaranteed overflow (> big.MaxExp) or
	// underflow (< big.MinExp) before allocating big.Int.Exp.
	if earlyOut, triggered := k.checkExponentBounds(x, exp); triggered {
		return earlyOut, true
	}

	// Decompose x into exact integer mantissa m and binary scale: x == m * 2^scaleExp2.
	m, scaleExp2, ok := k.decompose(x)
	if !ok {
		return nil, false
	}

	// Compute m^exp exactly using big.Int.Exp.
	expInt := new(big.Int).SetUint64(uint64(exp))
	mPow := new(big.Int).Exp(m, expInt, nil)

	// Safety check 2: compute total binary scaling exponent in int64 so
	// negative exponents cannot wrap on 32-bit platforms.
	totalScale64, ok := mulInt64Safe(scaleExp2, int64(exp))
	if !ok {
		// Overflowed int64 during scaling; defer safely to Option 2.
		return nil, false
	}

	// Load mPow losslessly into a temporary big.Float (BitLen() >= 1 because m != 0).
	bitLen := uint(mPow.BitLen())
	f := new(big.Float).SetPrec(bitLen).SetMode(big.ToNearestEven).SetInt(mPow)

	mantOnly := new(big.Float)
	powExp2 := f.MantExp(mantOnly) // mPow == mantOnly * 2^powExp2, 0.5 <= |mantOnly| < 1

	finalExp64 := int64(powExp2) + totalScale64
	neg := x.Signbit() && exp&1 == 1

	// Safety check 3 (Post-Exp): trigger explicit Inf or 0 if the exact binary
	// exponent is too large (> big.MaxExp) or too small (< big.MinExp).
	if finalExp64 > int64(big.MaxExp) {
		return k.signedInf(neg), true
	}
	if finalExp64 < int64(big.MinExp) {
		return k.signedZero(neg), true
	}

	// Because finalExp64 is verified inside [big.MinExp, big.MaxExp], casting
	// to int is safe on both 32-bit and 64-bit targets.
	scaled := new(big.Float).SetMantExp(mantOnly, int(finalExp64))

	// Single exact-to-target rounding pass.
	return new(big.Float).SetPrec(k.targetPrec).SetMode(k.targetMode).Set(scaled), true
}

// checkExponentBounds checks whether |x|^exp is guaranteed to exceed big.MaxExp
// or fall below big.MinExp using the MantExp interval 2^(e-1) <= |x| < 2^e.
func (k dyadicPowKernel) checkExponentBounds(x *big.Float, exp uint) (*big.Float, bool) {
	e := int64(x.MantExp(nil))
	neg := x.Signbit() && exp&1 == 1

	// If e > 1, |x| >= 2^(e-1), so |x^exp| >= 2^((e-1)*exp).
	// Since every finite big.Float is strictly < 2^MaxExp, >= 2^MaxExp overflows.
	if e > 1 {
		minResExp, ok := mulInt64Safe(e-1, int64(exp))
		if !ok || minResExp >= int64(big.MaxExp) {
			return k.signedInf(neg), true
		}
	}

	// If e < 0, |x| < 2^e, so |x^exp| < 2^(e*exp).
	// Since the smallest nonzero big.Float is 0.5 * 2^MinExp == 2^(MinExp-1),
	// anything < 2^(MinExp-1) underflows to zero.
	if e < 0 {
		maxResExp, ok := mulInt64Safe(e, int64(exp))
		if !ok || maxResExp <= int64(big.MinExp)-1 {
			return k.signedZero(neg), true
		}
	}

	return nil, false
}

// decompose extracts an exact integer mantissa m and binary scale factor
// scaleExp2 from a finite nonzero x such that x == m * 2^scaleExp2.
func (k dyadicPowKernel) decompose(x *big.Float) (*big.Int, int64, bool) {
	mant := new(big.Float)
	e := x.MantExp(mant) // x == mant * 2^e, with 0.5 <= |mant| < 1.0

	p := mant.MinPrec()
	if p == 0 || p > uint(big.MaxExp) {
		return nil, 0, false
	}

	// mant's current exponent is 0; shifting to p turns all p significant bits
	// into an exact integer without risking overflow.
	mant.SetMantExp(mant, int(p))

	m, acc := mant.Int(nil)
	if acc != big.Exact {
		return nil, 0, false
	}

	scaleExp2 := int64(e) - int64(p)
	return m, scaleExp2, true
}

func (k dyadicPowKernel) signedInf(neg bool) *big.Float {
	return new(big.Float).SetPrec(k.targetPrec).SetMode(k.targetMode).SetInf(neg)
}

func (k dyadicPowKernel) signedZero(neg bool) *big.Float {
	z := new(big.Float).SetPrec(k.targetPrec).SetMode(k.targetMode)
	if neg {
		z.Neg(z)
	}
	return z
}

func mulInt64Safe(a, b int64) (int64, bool) {
	if a == 0 || b == 0 {
		return 0, true
	}
	if (a == math.MinInt64 && b == -1) || (b == math.MinInt64 && a == -1) {
		return 0, false
	}
	res := a * b
	if res/a != b {
		return 0, false
	}
	return res, true
}
```

### How All Four Private Types Now Fit Together
* **`precConverter`:** Maps `targetFractionalDigits` $\leftrightarrow$ `prec`.
* **`floatEngine`:** Shared calculation runner (`powHeuristic`, `powZiv`, `powHybrid`, `powSpecialCase`).
* **`zivBounder`:** Error tracking & monotonic interval certification for Option 3 (`powZiv`).
* **`dyadicPowKernel`:** Exact `big.Int` decomposition, threshold gating (`<= 8192` bits), and exponent overflow/underflow protection for Option 4 (`powHybrid`).

Gemini 4 Hydrogen Question No 9

Gemini 4 Hydrogen Response No 9



Gemini 4 Hydrogen Question No 10

Gemini 4 Hydrogen Response No 10





# Summary

## The Plan

Going with **Option B (Encapsulated Wrapper)** with immutable methods., I would like to implement the following three strategies.

| **Option 2 (Fixed/Heuristic Guard Bits)** | Low    | Excellent                    | **Best general-purpose choice** for scientific, financial, or general Decimal libraries. |
| ----------------------------------------- | ------ | ---------------------------- | ------------------------------------------------------------ |
| **Option 3 (Ziv's Adaptive Strategy)**    | High   | Guaranteed Exact             | Standards-compliant math libraries where IEEE-style exact rounding is a hard requirement. |
| **Option 4 (Hybrid big.Int + Option 2)**  | Medium | Guaranteed Exact (small exp) | Use Option 4 when d.val.MinPrec() * exp is small (e.g., < 8192 bits)[[4](https://www.google.com/url?sa=E&q=https%3A%2F%2Fvertexaisearch.cloud.google.com%2Fgrounding-api-redirect%2FAUZIYQHswDv929zTuv2wMZtAonvBT2gKBFfLyrTU1LjI5tiSUOuEb3Al-Sxwrz98V2kzMTQha86RHBIDgklusz9m0EWi6kHLbmsdvzhyXaGisCWoOGNP1baIGPbU9gcnRhyhdfknzoJDpSgNZhvoWqY5trWEOFX3)], and fall back to Option 2 for larger exponents. |



## Components 

### Gemini 4 Hydrogen Response No 6 + Gemini 4 Hydrogen Response No 8 Add-In

```go

// ==================================================================================
//    Add In: Gemini 4 Hydrogen Response # 8
// ===================================================================================

const defaultMaxExactBits uint = 8192

type powStrategy uint8

const (
	strategyHeuristic powStrategy = iota // Option 2: Fixed/Heuristic Guard Bits
	strategyZiv                          // Option 3: Ziv's Adaptive Strategy
	strategyHybrid                       // Option 4: Hybrid big.Int + Option 2
)

// Pow returns d^exp using Option 2 (Fixed/Heuristic Guard Bits).
func (d *Decimal) Pow(exp uint) *Decimal {
	return d.powWithStrategy(exp, strategyHeuristic)
}

// PowExact returns d^exp using Option 3 (Ziv's Adaptive Strategy).
func (d *Decimal) PowExact(exp uint) *Decimal {
	return d.powWithStrategy(exp, strategyZiv)
}

// PowHybrid returns d^exp using Option 4 (Exact big.Int for small
// MinPrec()*exp, falling back to Option 2 for larger exponents).
func (d *Decimal) PowHybrid(exp uint) *Decimal {
	return d.powWithStrategy(exp, strategyHybrid)
}

func (d *Decimal) powWithStrategy(exp uint, strategy powStrategy) *Decimal {
	if d == nil || d.val == nil {
		return NewFromFloat64(0, defaultFractionalDigits, big.ToNearestEven).powWithStrategy(exp, strategy)
	}

	conv := defaultConverter()
	baseEng := floatEngine{prec: d.prec, mode: d.mode, maxPrec: maxWorkingPrec}

	resultExp2 := baseEng.estimatePowExp2(d.val, exp)
	resultPrec := conv.toPrec(d.targetFractionalDigits, resultExp2)
	resultEng := floatEngine{prec: resultPrec, mode: d.mode, maxPrec: maxWorkingPrec}

	var resultVal *big.Float
	switch strategy {
	case strategyZiv:
		resultVal = resultEng.powZiv(d.val, exp)
	case strategyHybrid:
		resultVal = resultEng.powHybrid(d.val, exp)
	default:
		resultVal = resultEng.powHeuristic(d.val, exp)
	}

	return &Decimal{
		val:                    resultVal,
		targetFractionalDigits: d.targetFractionalDigits,
		prec:                   resultPrec,
		mode:                   d.mode,
	}
}



// ======================================================================================
// Gemini Response # 6
// ======================================================================================
package decimal

import (
	"math"
	"math/big"
	"math/bits"
)

const (
	defaultFractionalDigits uint = 6
	defaultSafetyBits       uint = 8       // Extra binary cushion for decimal rounding
	maxWorkingPrec          uint = 1 << 20 // Safety cap (~1M bits) to prevent runaway allocations
)

// Decimal is an immutable, encapsulated wrapper around *big.Float (Option B).
type Decimal struct {
	val                    *big.Float
	targetFractionalDigits uint
	prec                   uint
	mode                   big.RoundingMode
}

// NewFromFloat64 constructs an immutable Decimal sized to preserve at least
// targetFractionalDigits decimal places for v's magnitude.
func NewFromFloat64(v float64, targetFractionalDigits uint, mode big.RoundingMode) *Decimal {
	conv := defaultConverter()

	var exp2 int
	if v != 0 && !math.IsInf(v, 0) && !math.IsNaN(v) {
		_, exp2 = math.Frexp(v)
	}

	prec := conv.toPrec(targetFractionalDigits, exp2)
	val := new(big.Float).SetPrec(prec).SetMode(mode).SetFloat64(v)

	return &Decimal{
		val:                    val,
		targetFractionalDigits: targetFractionalDigits,
		prec:                   prec,
		mode:                   mode,
	}
}

// TargetFractionalDigits returns the requested base-10 digits after the decimal point.
func (d *Decimal) TargetFractionalDigits() uint {
	if d == nil {
		return 0
	}
	return d.targetFractionalDigits
}

// ActualFractionalDigits derives how many base-10 fractional digits are guaranteed
// by the underlying big.Float's binary precision at its current magnitude.
func (d *Decimal) ActualFractionalDigits() uint {
	if d == nil || d.val == nil || d.val.Sign() == 0 || d.val.IsInf() {
		return d.targetFractionalDigits
	}
	return defaultConverter().toFractionalDigits(d.prec, d.val.MantExp(nil))
}

// Prec returns the underlying big.Float total mantissa precision in bits.
func (d *Decimal) Prec() uint {
	if d == nil {
		return 0
	}
	return d.prec
}

// Pow returns d^exp using Option 2 (Fixed/Heuristic Guard Bits).
func (d *Decimal) Pow(exp uint) *Decimal {
	return d.powWithStrategy(exp, false)
}

// PowExact returns d^exp using Option 3 (Ziv's Adaptive Strategy) for
// guaranteed correctly rounded output at the computed target precision.
func (d *Decimal) PowExact(exp uint) *Decimal {
	return d.powWithStrategy(exp, true)
}

func (d *Decimal) powWithStrategy(exp uint, exact bool) *Decimal {
	if d == nil || d.val == nil {
		return NewFromFloat64(0, defaultFractionalDigits, big.ToNearestEven).powWithStrategy(exp, exact)
	}

	conv := defaultConverter()
	baseEng := floatEngine{prec: d.prec, mode: d.mode, maxPrec: maxWorkingPrec}

	resultExp2 := baseEng.estimatePowExp2(d.val, exp)
	resultPrec := conv.toPrec(d.targetFractionalDigits, resultExp2)
	resultEng := floatEngine{prec: resultPrec, mode: d.mode, maxPrec: maxWorkingPrec}

	var resultVal *big.Float
	if exact {
		resultVal = resultEng.powZiv(d.val, exp)
	} else {
		resultVal = resultEng.powHeuristic(d.val, exp)
	}

	return &Decimal{
		val:                    resultVal,
		targetFractionalDigits: d.targetFractionalDigits,
		prec:                   resultPrec,
		mode:                   d.mode,
	}
}

// String formats the Decimal in fixed-point notation to targetFractionalDigits.
func (d *Decimal) String() string {
	if d == nil || d.val == nil {
		return "0"
	}
	if d.val.IsInf() {
		return d.val.String()
	}
	return d.val.Text('f', int(d.targetFractionalDigits))
}

// -----------------------------------------------------------------------------
// Low-Level Private Precision Converter
// -----------------------------------------------------------------------------

const (
	log2Of10 = math.Ln10 / math.Ln2
	log10Of2 = math.Ln2 / math.Ln10
)

type precConverter struct {
	safetyBits uint
	maxPrec    uint
}

func defaultConverter() precConverter {
	return precConverter{
		safetyBits: defaultSafetyBits,
		maxPrec:    maxWorkingPrec,
	}
}

func (c precConverter) toPrec(targetFractionalDigits uint, exp2 int) uint {
	fracBits := uint(math.Ceil(float64(targetFractionalDigits) * log2Of10))

	var intBits uint
	if exp2 > 0 {
		intBits = uint(exp2)
	}

	total := intBits + fracBits
	if total < intBits || total+c.safetyBits < total {
		return c.maxPrec
	}
	total += c.safetyBits

	if total == 0 {
		total = 1
	}
	if c.maxPrec > 0 && total > c.maxPrec {
		return c.maxPrec
	}
	return total
}

func (c precConverter) toFractionalDigits(prec uint, exp2 int) uint {
	if prec <= c.safetyBits {
		return 0
	}
	usablePrec := prec - c.safetyBits

	var fracBits uint
	if exp2 > 0 {
		intBits := uint(exp2)
		if usablePrec <= intBits {
			return 0
		}
		fracBits = usablePrec - intBits
	} else {
		fracBits = usablePrec + uint(-exp2)
	}

	return uint(math.Floor(float64(fracBits) * log10Of2))
}

// -----------------------------------------------------------------------------
// Private Reusable Calculation Engine
// -----------------------------------------------------------------------------

type floatEngine struct {
	prec    uint
	mode    big.RoundingMode
	maxPrec uint
}

func (e floatEngine) estimatePowExp2(x *big.Float, exp uint) int {
	if exp == 0 || x.Sign() == 0 || x.IsInf() {
		return 1
	}

	mant := new(big.Float)
	exp2 := x.MantExp(mant)
	mant.Abs(mant)

	m64, _ := mant.Float64()
	est := float64(exp) * (math.Log2(m64) + float64(exp2))

	if est <= 0 {
		return 0
	}
	if est >= float64(big.MaxExp) {
		return big.MaxExp
	}
	return int(math.Ceil(est)) + 1
}

// powSpecialCase handles exp in {0, 1}, ±0, ±Inf, and ±1 in one shared place.
func (e floatEngine) powSpecialCase(x *big.Float, exp uint) (*big.Float, bool) {
	switch exp {
	case 0:
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(1), true
	case 1:
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).Set(x), true
	}

	if x.Sign() == 0 {
		z := new(big.Float).SetPrec(e.prec).SetMode(e.mode)
		if x.Signbit() && exp&1 == 1 {
			z.Neg(z)
		}
		return z, true
	}

	if x.IsInf() {
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInf(x.Signbit() && exp&1 == 1), true
	}

	if x.Cmp(big.NewFloat(1)) == 0 {
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(1), true
	}
	if x.Cmp(big.NewFloat(-1)) == 0 {
		if exp&1 == 1 {
			return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(-1), true
		}
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(1), true
	}

	return nil, false
}

// powHeuristic computes x^exp using Option 2 (Fixed/Heuristic Guard Bits).
func (e floatEngine) powHeuristic(x *big.Float, exp uint) *big.Float {
	if out, ok := e.powSpecialCase(x, exp); ok {
		return out
	}

	guardBits := uint(bits.Len(exp)) + 32
	workPrec := e.prec + guardBits
	if workPrec < e.prec || (e.maxPrec > 0 && workPrec > e.maxPrec) {
		workPrec = e.maxPrec
	}

	res := new(big.Float).SetPrec(workPrec).SetMode(big.ToNearestEven).SetInt64(1)
	base := new(big.Float).SetPrec(workPrec).SetMode(big.ToNearestEven).Set(x)

	origExp := exp
	for exp > 0 {
		if exp&1 == 1 {
			res.Mul(res, base)
		}
		if exp > 1 {
			base.Mul(base, base)
			if base.IsInf() {
				return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInf(x.Signbit() && origExp&1 == 1)
			}
			if base.Sign() == 0 {
				z := new(big.Float).SetPrec(e.prec).SetMode(e.mode)
				if x.Signbit() && origExp&1 == 1 {
					z.Neg(z)
				}
				return z
			}
		}
		exp >>= 1
	}

	return new(big.Float).SetPrec(e.prec).SetMode(e.mode).Set(res)
}

// powZiv computes x^exp using Option 3 (Ziv's Adaptive Strategy), iteratively
// increasing working precision until the final rounding to e.prec is provably
// unambiguous.
func (e floatEngine) powZiv(x *big.Float, exp uint) *big.Float {
	if out, ok := e.powSpecialCase(x, exp); ok {
		return out
	}

	// Ensure basePrec covers both target precision and x's exact mantissa bits
	// so the initial load of base = x introduces zero quantization error.
	basePrec := e.prec
	if minP := x.MinPrec(); minP > basePrec {
		basePrec = minP
	}

	// Start with enough guard bits so almost all inputs succeed on Pass 1.
	guardBits := uint(bits.Len(exp)) + 32

	for {
		workPrec := basePrec + guardBits
		capped := false
		if workPrec < basePrec || (e.maxPrec > 0 && workPrec >= e.maxPrec) {
			workPrec = e.maxPrec
			capped = true
		}

		bounder := newZivBounder(e.prec, e.mode, workPrec)
		res, errUnits, earlyOut := e.runZivPass(x, exp, bounder)
		if earlyOut != nil {
			return earlyOut
		}

		// Verify whether rounding res to (e.prec, e.mode) is unambiguous.
		if rounded, ok := bounder.tryRound(res, errUnits); ok {
			return rounded
		}

		// Safeguard: stop if we hit the maximum precision cap.
		if capped {
			return new(big.Float).SetPrec(e.prec).SetMode(e.mode).Set(res)
		}

		// Double guard bits for the next Ziv refinement pass.
		guardBits *= 2
	}
}

// runZivPass executes one binary-exponentiation pass at b.workPrec while
// accumulating a deterministic upper bound on relative rounding error.
func (e floatEngine) runZivPass(x *big.Float, exp uint, b zivBounder) (*big.Float, float64, *big.Float) {
	res := new(big.Float).SetPrec(b.workPrec).SetMode(big.ToNearestEven).SetInt64(1)
	base := new(big.Float).SetPrec(b.workPrec).SetMode(big.ToNearestEven).Set(x)

	resErr := 0.0
	baseErr := 0.0
	if base.Acc() != big.Exact {
		// Only happens if x.MinPrec() exceeded maxWorkingPrec.
		baseErr = 1.0
	}

	origExp := exp
	for exp > 0 {
		if exp&1 == 1 {
			res.Mul(res, base)
			resErr = b.propagateMulErr(resErr, baseErr, res.Acc())
		}
		if exp > 1 {
			base.Mul(base, base)
			baseErr = b.propagateMulErr(baseErr, baseErr, base.Acc())

			if base.IsInf() {
				inf := new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInf(x.Signbit() && origExp&1 == 1)
				return nil, 0, inf
			}
			if base.Sign() == 0 {
				z := new(big.Float).SetPrec(e.prec).SetMode(e.mode)
				if x.Signbit() && origExp&1 == 1 {
					z.Neg(z)
				}
				return nil, 0, z
			}
		}
		exp >>= 1
	}

	return res, resErr, nil
}

// -----------------------------------------------------------------------------
// Private Ziv Error-Bounding & Interval Verification Type
// -----------------------------------------------------------------------------

// zivBounder tracks relative rounding error in units of u = 2^-workPrec and
// verifies whether rounding a working-precision result to (targetPrec, targetMode)
// is provably unambiguous.
type zivBounder struct {
	targetPrec   uint
	targetMode   big.RoundingMode
	workPrec     uint
	unitRoundoff float64 // 2^-workPrec (upper bound on relative error for ToNearestEven)
}

func newZivBounder(targetPrec uint, targetMode big.RoundingMode, workPrec uint) zivBounder {
	return zivBounder{
		targetPrec:   targetPrec,
		targetMode:   targetMode,
		workPrec:     workPrec,
		unitRoundoff: math.Ldexp(1.0, -int(workPrec)),
	}
}

// propagateMulErr computes the upper bound on relative error (in units of u)
// after multiplying two values with error bounds errA and errB.
func (z zivBounder) propagateMulErr(errA, errB float64, acc big.Accuracy) float64 {
	// (1 + errA*u)*(1 + errB*u) - 1 = (errA + errB + errA*errB*u)*u
	prop := errA + errB + errA*errB*z.unitRoundoff
	if acc == big.Exact {
		return prop
	}
	// Multiply by (1 + u) for the new rounding step.
	return prop + 1.0 + prop*z.unitRoundoff
}

// tryRound checks if the true mathematical value in [res - delta, res + delta]
// rounds identically at both endpoints under (targetPrec, targetMode).
func (z zivBounder) tryRound(res *big.Float, errUnits float64) (*big.Float, bool) {
	// Fast path: if every operation was Exact (e.g., exact powers or exact ties
	// like 1.5^2 = 2.25), res has zero error and rounds immediately.
	if errUnits == 0 {
		return new(big.Float).SetPrec(z.targetPrec).SetMode(z.targetMode).Set(res), true
	}

	if res.IsInf() {
		return new(big.Float).SetPrec(z.targetPrec).SetMode(z.targetMode).Set(res), true
	}

	// If accumulated error exceeds 0.25 ULP of targetPrec, do not bother with
	// the interval check; signal powZiv to increase guardBits immediately.
	var guardBits int
	if z.workPrec > z.targetPrec {
		guardBits = int(z.workPrec - z.targetPrec)
	}
	if errUnits >= math.Ldexp(0.25, guardBits) {
		return nil, false
	}

	// Apply a 1% conservative multiplier to errUnits to cover both float64
	// rounding and converting relative error |Y_true - res| / |Y_true| to |res|.
	safeErrUnits := errUnits * 1.01

	// Build eps >= safeErrUnits * 2^-workPrec rounded upward (ToPositiveInf).
	eps := new(big.Float).SetPrec(64).SetMode(big.ToPositiveInf).SetFloat64(safeErrUnits)
	eps.SetMantExp(eps, eps.MantExp(nil)-int(z.workPrec))

	// delta >= |res| * eps, rounded upward.
	delta := new(big.Float).SetPrec(64).SetMode(big.ToPositiveInf).Abs(res)
	delta.Mul(delta, eps)

	// Construct [lo, hi] strictly enclosing the true value Y_true.
	lo := new(big.Float).SetPrec(z.workPrec).SetMode(big.ToNegativeInf).Sub(res, delta)
	hi := new(big.Float).SetPrec(z.workPrec).SetMode(big.ToPositiveInf).Add(res, delta)

	// Round both endpoints to target precision and target rounding mode.
	rLo := new(big.Float).SetPrec(z.targetPrec).SetMode(z.targetMode).Set(lo)
	rHi := new(big.Float).SetPrec(z.targetPrec).SetMode(z.targetMode).Set(hi)

	// By monotonicity of rounding, rLo == rHi guarantees round(Y_true) == rLo.
	if rLo.Cmp(rHi) == 0 {
		return rLo, true
	}
	return nil, false
}

// ================================================================================
// Gemini Response # 8 Add in
// ================================================================================
// -----------------------------------------------------------------------------
// Addition to floatEngine
// -----------------------------------------------------------------------------

// powHybrid computes x^exp using exact big.Int exponentiation when
// x.MinPrec() * exp <= defaultMaxExactBits, safely falling back to Option 2
// (powHeuristic) when the exponent or bit-growth is too large.
func (e floatEngine) powHybrid(x *big.Float, exp uint) *big.Float {
	// 1. Handle exp in {0, 1}, ±0, ±Inf, and ±1 immediately.
	if out, ok := e.powSpecialCase(x, exp); ok {
		return out
	}

	// 2. Attempt exact big.Int calculation if within the bit-budget threshold.
	kernel := dyadicPowKernel{
		targetPrec:   e.prec,
		targetMode:   e.mode,
		maxExactBits: defaultMaxExactBits,
	}

	if kernel.canUseExact(x, exp) {
		if out, ok := kernel.powExact(x, exp); ok {
			return out
		}
	}

	// 3. Fallback to Option 2 (Fixed/Heuristic Guard Bits) for large exponents
	// or if exact decomposition was declined.
	return e.powHeuristic(x, exp)
}

// -----------------------------------------------------------------------------
// Private Exact Dyadic / big.Int Exponentiation Kernel
// -----------------------------------------------------------------------------

// dyadicPowKernel decomposes a finite big.Float into m * 2^scaleExp2, computes
// m^exp exactly via big.Int.Exp, rescales the binary exponent, and rounds once
// to (targetPrec, targetMode).
type dyadicPowKernel struct {
	targetPrec   uint
	targetMode   big.RoundingMode
	maxExactBits uint
}

// canUseExact reports whether x.MinPrec() * exp fits within maxExactBits
// without risking uint multiplication overflow.
func (k dyadicPowKernel) canUseExact(x *big.Float, exp uint) bool {
	if exp <= 1 || k.maxExactBits == 0 {
		return false
	}
	minPrec := x.MinPrec()
	if minPrec == 0 || minPrec > k.maxExactBits {
		return false
	}
	// Overflow-safe check for (minPrec * exp <= k.maxExactBits):
	return exp <= k.maxExactBits/minPrec
}

// powExact executes exact big.Int exponentiation with pre- and post-exponent
// overflow/underflow guards. Returns (nil, false) if it cannot safely proceed.
func (k dyadicPowKernel) powExact(x *big.Float, exp uint) (*big.Float, bool) {
	if uint64(exp) > math.MaxInt64 {
		return nil, false
	}

	// Safety check 1 (Pre-Exp): detect guaranteed overflow (> big.MaxExp) or
	// underflow (< big.MinExp) before allocating big.Int.Exp.
	if earlyOut, triggered := k.checkExponentBounds(x, exp); triggered {
		return earlyOut, true
	}

	// Decompose x into exact integer mantissa m and binary scale: x == m * 2^scaleExp2.
	m, scaleExp2, ok := k.decompose(x)
	if !ok {
		return nil, false
	}

	// Compute m^exp exactly using big.Int.Exp.
	expInt := new(big.Int).SetUint64(uint64(exp))
	mPow := new(big.Int).Exp(m, expInt, nil)

	// Safety check 2: compute total binary scaling exponent in int64 so
	// negative exponents cannot wrap on 32-bit platforms.
	totalScale64, ok := mulInt64Safe(scaleExp2, int64(exp))
	if !ok {
		// Overflowed int64 during scaling; defer safely to Option 2.
		return nil, false
	}

	// Load mPow losslessly into a temporary big.Float (BitLen() >= 1 because m != 0).
	bitLen := uint(mPow.BitLen())
	f := new(big.Float).SetPrec(bitLen).SetMode(big.ToNearestEven).SetInt(mPow)

	mantOnly := new(big.Float)
	powExp2 := f.MantExp(mantOnly) // mPow == mantOnly * 2^powExp2, 0.5 <= |mantOnly| < 1

	finalExp64 := int64(powExp2) + totalScale64
	neg := x.Signbit() && exp&1 == 1

	// Safety check 3 (Post-Exp): trigger explicit Inf or 0 if the exact binary
	// exponent is too large (> big.MaxExp) or too small (< big.MinExp).
	if finalExp64 > int64(big.MaxExp) {
		return k.signedInf(neg), true
	}
	if finalExp64 < int64(big.MinExp) {
		return k.signedZero(neg), true
	}

	// Because finalExp64 is verified inside [big.MinExp, big.MaxExp], casting
	// to int is safe on both 32-bit and 64-bit targets.
	scaled := new(big.Float).SetMantExp(mantOnly, int(finalExp64))

	// Single exact-to-target rounding pass.
	return new(big.Float).SetPrec(k.targetPrec).SetMode(k.targetMode).Set(scaled), true
}

// checkExponentBounds checks whether |x|^exp is guaranteed to exceed big.MaxExp
// or fall below big.MinExp using the MantExp interval 2^(e-1) <= |x| < 2^e.
func (k dyadicPowKernel) checkExponentBounds(x *big.Float, exp uint) (*big.Float, bool) {
	e := int64(x.MantExp(nil))
	neg := x.Signbit() && exp&1 == 1

	// If e > 1, |x| >= 2^(e-1), so |x^exp| >= 2^((e-1)*exp).
	// Since every finite big.Float is strictly < 2^MaxExp, >= 2^MaxExp overflows.
	if e > 1 {
		minResExp, ok := mulInt64Safe(e-1, int64(exp))
		if !ok || minResExp >= int64(big.MaxExp) {
			return k.signedInf(neg), true
		}
	}

	// If e < 0, |x| < 2^e, so |x^exp| < 2^(e*exp).
	// Since the smallest nonzero big.Float is 0.5 * 2^MinExp == 2^(MinExp-1),
	// anything < 2^(MinExp-1) underflows to zero.
	if e < 0 {
		maxResExp, ok := mulInt64Safe(e, int64(exp))
		if !ok || maxResExp <= int64(big.MinExp)-1 {
			return k.signedZero(neg), true
		}
	}

	return nil, false
}

// decompose extracts an exact integer mantissa m and binary scale factor
// scaleExp2 from a finite nonzero x such that x == m * 2^scaleExp2.
func (k dyadicPowKernel) decompose(x *big.Float) (*big.Int, int64, bool) {
	mant := new(big.Float)
	e := x.MantExp(mant) // x == mant * 2^e, with 0.5 <= |mant| < 1.0

	p := mant.MinPrec()
	if p == 0 || p > uint(big.MaxExp) {
		return nil, 0, false
	}

	// mant's current exponent is 0; shifting to p turns all p significant bits
	// into an exact integer without risking overflow.
	mant.SetMantExp(mant, int(p))

	m, acc := mant.Int(nil)
	if acc != big.Exact {
		return nil, 0, false
	}

	scaleExp2 := int64(e) - int64(p)
	return m, scaleExp2, true
}

func (k dyadicPowKernel) signedInf(neg bool) *big.Float {
	return new(big.Float).SetPrec(k.targetPrec).SetMode(k.targetMode).SetInf(neg)
}

func (k dyadicPowKernel) signedZero(neg bool) *big.Float {
	z := new(big.Float).SetPrec(k.targetPrec).SetMode(k.targetMode)
	if neg {
		z.Neg(z)
	}
	return z
}

func mulInt64Safe(a, b int64) (int64, bool) {
	if a == 0 || b == 0 {
		return 0, true
	}
	if (a == math.MinInt64 && b == -1) || (b == math.MinInt64 && a == -1) {
		return 0, false
	}
	res := a * b
	if res/a != b {
		return 0, false
	}
	return res, true
}

```







# Cost on Paid Plan

If I were on the paid plan for Gemini 4 Hydrogen, can you reasonably estimate how much this session (74,494 tokens) would have cost me?

------

On the **Gemini Developer API Standard Paid Tier** for **Gemini 4 Hydrogen** (`gemini-4-hydrogen`), a session totaling **74,494 billed tokens** would reasonably cost around **$0.21 to $0.29** (roughly **21¢–29¢**), depending on your input/output split and implicit cache hits.

*(Note: If by "paid plan" you meant a flat-rate consumer subscription like Google AI Plus/Pro/Ultra, you wouldn't be billed per token at all. The breakdown below assumes pay-as-you-go API / Google AI Studio Paid Tier billing.)*

---

### 1. Official Gemini 4 Hydrogen Paid-Tier Rates (Standard)
* **Input tokens:** **$1.50** per 1M tokens (flat across the 1M context window)
* **Cached input tokens (Implicit caching):** **$0.075** per 1M tokens (95% discount; active by default once the shared prompt prefix reaches 4,096 tokens)
* **Output tokens (includes visible text + thinking tokens):** **$7.50** per 1M tokens
* **Google Search Grounding:** First 5,000 requests/month are free across Gemini 3+ models (then $14 per 1,000 queries)

---

### 2. Why `74,494` Tokens Gives a Range
Because output (including internal reasoning/thinking tokens) costs **5×** as much as uncached input ($7.50 vs. $1.50), the exact cost depends on how those 74,494 tokens break down:

* **Hard mathematical bounds (uncached):**
  * If 100% were uncached input: $0.074494 \times \$1.50 = \mathbf{\$0.112}$
  * If 100% were output/thinking: $0.074494 \times \$7.50 = \mathbf{\$0.559}$

---

### 3. Reasonable Estimate for This Session
In a multi-turn coding and architecture conversation like ours (short prompts from you, long Go code generation and `high`-default thinking tokens from the model, with history re-sent on each turn), a typical cumulative split is roughly **~67% input (~50,000 tokens)** and **~33% output + thinking (~24,494 tokens)**:

#### Scenario A: Without Cache Hits (Conservative Estimate)
* **Input:** $50,000 \text{ tokens} \times \$1.50 / 1\text{M} = \$0.0750$
* **Output (incl. thinking):** $24,494 \text{ tokens} \times \$7.50 / 1\text{M} = \$0.1837$
* **Estimated Total:** **~$0.26** *(or ~$0.20–$0.34 across a 20%–50% output share)*

#### Scenario B: With Implicit Context Caching
Once our conversation crossed 4,096 tokens around Turn 3–4, the repeated prefix from earlier turns qualifies for implicit caching ($0.075/1M). Assuming ~15,000 fresh input tokens and ~35,000 cached prefix tokens:
* **Fresh Input:** $15,000 \times \$1.50 / 1\text{M} = \$0.0225$
* **Cached Input:** $35,000 \times \$0.075 / 1\text{M} = \$0.0026$
* **Output (incl. thinking):** $24,494 \times \$7.50 / 1\text{M} = \$0.1837$
* **Estimated Total:** **~$0.21**

*(One caveat: if `74,494` was read from a "current context window" meter rather than a cumulative session counter—meaning earlier turns were billed separately before reaching 74,494—your cumulative uncached input cost would be higher, though implicit caching still keeps the total roughly around **$0.25–$0.40** since output/thinking tokens dominate the bill.)*