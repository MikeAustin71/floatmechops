# Go `math/big`: RoundingMode Specification Guide

In Go's `math/big` package, `big.Float` calculations operate with arbitrary precision. When a result cannot be represented exactly within the requested bit precision (significand bits), it must be rounded. The `RoundingMode` type determines how rounding decisions and tie-breaking are executed.

Most modes directly implement the **IEEE 754-2008** (and identical IEEE 754-2019) floating-point specification for *Rounding-direction attributes* (Section 4.3).

The following definition of **`type RoundingMode`** was taken from the **Go** **`math/big`** package:

```go
// RoundingMode determines how a [Float] value is rounded to the
// desired precision. Rounding may change the [Float] value; the
// rounding error is described by the [Float]'s [Accuracy].
type RoundingMode byte

// These constants define supported rounding modes.
const (
	ToNearestEven RoundingMode = iota // == IEEE 754-2008 roundTiesToEven
	ToNearestAway                     // == IEEE 754-2008 roundTiesToAway
	ToZero                            // == IEEE 754-2008 roundTowardZero
	AwayFromZero                      // no IEEE 754-2008 equivalent
	ToNegativeInf                     // == IEEE 754-2008 roundTowardNegative
	ToPositiveInf                     // == IEEE 754-2008 roundTowardPositive
)
```





---

## 1. `ToNearestEven`

- **Go Constant:** `ToNearestEven RoundingMode = 0`
- **IEEE 754-2008 Reference:** **Section 4.3.1** (`roundTiesToEven`)
- **Also Known As:** Bankers' Rounding, Round-to-Nearest-Even (default IEEE binary mode).

### Description
Rounds to the nearest representable value. If the unrounded value falls halfway between two equally close representable values (a tie), it rounds to the value whose least significant bit (or digit) is **even**. This minimizes statistical bias across repeated operations.

### Examples (Rounding to Nearest Integer)
- **Positive:**
  - `+1.4` $\to$ `+1` *(closer to 1)*
  - `+1.5` $\to$ `+2` *(tie: 2 is even)*
  - `+2.5` $\to$ `+2` *(tie: 2 is even)*
  - `+3.5` $\to$ `+4` *(tie: 4 is even)*
- **Negative:**
  - `-1.4` $\to$ `-1` *(closer to -1)*
  - `-1.5` $\to$ `-2` *(tie: -2 is even)*
  - `-2.5` $\to$ `-2` *(tie: -2 is even)*
  - `-3.5` $\to$ `-4` *(tie: -4 is even)*

---

## 2. `ToNearestAway`

- **Go Constant:** `ToNearestAway RoundingMode = 1`
- **IEEE 754-2008 Reference:** **Section 4.3.1** (`roundTiesToAway`)
- **Also Known As:** Commercial Rounding, Round Half Up (in absolute value), Round Half Away From Zero.

### Description
Rounds to the nearest representable value. If the unrounded value is exactly halfway between two representable values, it rounds to the one with the **larger magnitude** (away from zero).

### Examples (Rounding to Nearest Integer)
- **Positive:**
  - `+1.4` $\to$ `+1` *(closer to 1)*
  - `+1.5` $\to$ `+2` *(tie: rounds away from 0)*
  - `+2.5` $\to$ `+3` *(tie: rounds away from 0)*
- **Negative:**
  - `-1.4` $\to$ `-1` *(closer to -1)*
  - `-1.5` $\to$ `-2` *(tie: rounds away from 0 / larger magnitude)*
  - `-2.5` $\to$ `-3` *(tie: rounds away from 0 / larger magnitude)*

---

## 3. `ToZero`

- **Go Constant:** `ToZero RoundingMode = 2`
- **IEEE 754-2008 Reference:** **Section 4.3.2** (`roundTowardZero`)
- **Also Known As:** Truncation, Round Toward Zero.

### Description
Directed rounding that discards excess precision toward zero. It chooses the representable value closest to, but **not greater in magnitude** than, the exact value.

### Examples (Rounding to Nearest Integer)
- **Positive:**
  - `+1.1` $\to$ `+1`
  - `+1.5` $\to$ `+1`
  - `+1.9` $\to$ `+1`
- **Negative:**
  - `-1.1` $\to$ `-1`
  - `-1.5` $\to$ `-1`
  - `-1.9` $\to$ `-1`

---

## 4. `AwayFromZero`

- **Go Constant:** `AwayFromZero RoundingMode = 3`
- **IEEE 754-2008 Reference:** **No IEEE 754-2008 equivalent** *(Go extension)*

### Description
Directed rounding that always moves **away from zero** whenever the value cannot be represented exactly. It selects the representable value closest to, but **not smaller in magnitude** than, the exact value (the polar opposite of truncation).

### Examples (Rounding to Nearest Integer)
- **Positive:**
  - `+1.1` $\to$ `+2`
  - `+1.5` $\to$ `+2`
  - `+1.9` $\to$ `+2`
- **Negative:**
  - `-1.1` $\to$ `-2`
  - `-1.5` $\to$ `-2`
  - `-1.9` $\to$ `-2`

---

## 5. `ToNegativeInf`

- **Go Constant:** `ToNegativeInf RoundingMode = 4`
- **IEEE 754-2008 Reference:** **Section 4.3.2** (`roundTowardNegative`)
- **Also Known As:** Floor, Round Toward $-\infty$.

### Description
Directed rounding downward. It chooses the representable value closest to, but **less than or equal to** ($\le$), the exact value.

### Examples (Rounding to Nearest Integer)
- **Positive:**
  - `+1.1` $\to$ `+1`
  - `+1.5` $\to$ `+1`
  - `+1.9` $\to$ `+1`
- **Negative:**
  - `-1.1` $\to$ `-2`
  - `-1.5` $\to$ `-2`
  - `-1.9` $\to$ `-2`

---

## 6. `ToPositiveInf`

- **Go Constant:** `ToPositiveInf RoundingMode = 5`
- **IEEE 754-2008 Reference:** **Section 4.3.2** (`roundTowardPositive`)
- **Also Known As:** Ceiling, Round Toward $+\infty$.

### Description
Directed rounding upward. It chooses the representable value closest to, but **greater than or equal to** ($\ge$), the exact value.

### Examples (Rounding to Nearest Integer)
- **Positive:**
  - `+1.1` $\to$ `+2`
  - `+1.5` $\to$ `+2`
  - `+1.9` $\to$ `+2`
- **Negative:**
  - `-1.1` $\to$ `-1`
  - `-1.5` $\to$ `-1`
  - `-1.9` $\to$ `-1`

---

## Comparative Summary

The table below illustrates how different exact numbers are rounded to the nearest whole unit under each mode:

| Input Value | `ToNearestEven` (§4.3.1) | `ToNearestAway` (§4.3.1) | `ToZero` (§4.3.2) | `AwayFromZero` (None) | `ToNegativeInf` (§4.3.2) | `ToPositiveInf` (§4.3.2) |
| :---------: | :----------------------: | :----------------------: | :---------------: | :-------------------: | :----------------------: | :----------------------: |
| **`+2.5`**  |      `+2` *(even)*       |           `+3`           |       `+2`        |         `+3`          |           `+2`           |           `+3`           |
| **`+1.5`**  |      `+2` *(even)*       |           `+2`           |       `+1`        |         `+2`          |           `+1`           |           `+2`           |
| **`+1.2`**  |           `+1`           |           `+1`           |       `+1`        |         `+2`          |           `+1`           |           `+2`           |
| **`-1.2`**  |           `-1`           |           `-1`           |       `-1`        |         `-2`          |           `-2`           |           `-1`           |
| **`-1.5`**  |      `-2` *(even)*       |           `-2`           |       `-1`        |         `-2`          |           `-2`           |           `-1`           |
| **`-2.5`**  |      `-2` *(even)*       |           `-3`           |       `-2`        |         `-3`          |           `-3`           |           `-2`           |