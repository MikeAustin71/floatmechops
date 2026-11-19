package floatmech

import (
	"fmt"
	"math/big"

	ePref "github.com/MikeAustin71/errpref/v2"
)

// RoundingType
//
// Represents the IEEE 754-2008 rounding mode used when converting a numeric
// value to a lower precision or to an integer representation.
//
// RoundingType defines six rounding modes derived from IEEE 754-2008 sections
// 4.3.1 and 4.3.2, plus one Go-specific extension (AwayFromZero). The modes
// govern how inexact results are rounded when they cannot be represented
// exactly.
//
// Comparative Summary
// -------------------
//
// The table below illustrates how different exact numbers are rounded to the
// nearest whole unit under each mode:
//
// | Input Value | ToNearestEven (§4.3.1) | ToNearestAway (§4.3.1) | ToZero (§4.3.2) | AwayFromZero (None) | ToNegativeInf (§4.3.2) | ToPositiveInf (§4.3.2) |
// | :---------: | :----------------------: | :----------------------: | :---------------: | :-------------------: | :----------------------: | :----------------------: |
// | **+2.5**    |      +2 *(even)*       |           +3           |       +2          |         +3          |           +2           |           +3           |
// | **+1.5**    |      +2 *(even)*       |           +2           |       +1          |         +2          |           +1           |           +2           |
// | **+1.2**    |           +1           |           +1           |       +1          |         +2          |           +1           |           +2           |
// | **-1.2**    |           -1           |           -1           |       -1          |         -2          |           -2           |           -1           |
// | **-1.5**    |      -2 *(even)*       |           -2           |       -1          |         -2          |           -2           |           -1           |
// | **-2.5**    |      -2 *(even)*       |           -3           |       -2          |         -3          |           -3           |           -2           |
//
// Contracts/Assumptions
// ---------------------
//
// - The zero value of RoundingType is an uninitialized / invalid state.
// - All rounding modes are deterministic and stateless.
// - ToNearestEven is the default IEEE 754 binary rounding mode.
type RoundingType int

// ToNearestEvenVal, ToNearestAwayVal, ToZeroVal, AwayFromZeroVal,
// ToNegativeInfVal, ToPositiveInfVal
//
// Enumerated constants representing the IEEE 754-2008 rounding modes supported
// by RoundingType. The values correspond to IEEE 754-2008 section 4.3.1
// (roundTiesToEven, roundTiesToAway) and section 4.3.2 (roundTowardZero,
// roundTowardNegative, roundTowardPositive), plus one Go-specific extension
// (AwayFromZero) with no IEEE 754-2008 equivalent.
//
// ToNearestEvenVal (0) - IEEE 754-2008 roundTiesToEven (Bankers' Rounding):
//
//	Rounds to the nearest representable value. If the unrounded value falls
//	halfway between two equally close representable values (a tie), it rounds
//	to the value whose least significant bit (or digit) is even. This
//	minimizes statistical bias across repeated operations.
//
//	Examples (rounding to nearest integer):
//
//	  +1.4 -> +1  (closer to 1)
//	  +1.5 -> +2  (tie: 2 is even)
//	  +2.5 -> +2  (tie: 2 is even)
//	  +3.5 -> +4  (tie: 4 is even)
//	  -1.4 -> -1  (closer to -1)
//	  -1.5 -> -2  (tie: -2 is even)
//	  -2.5 -> -2  (tie: -2 is even)
//	  -3.5 -> -4  (tie: -4 is even)
//
// ToNearestAwayVal (1) - IEEE 754-2008 roundTiesToAway (Commercial Rounding):
//
//	Rounds to the nearest representable value. If the unrounded value is
//	exactly halfway between two representable values, it rounds to the one
//	with the larger magnitude (away from zero).
//
//	Examples (rounding to nearest integer):
//
//	  +1.4 -> +1  (closer to 1)
//	  +1.5 -> +2  (tie: rounds away from 0)
//	  +2.5 -> +3  (tie: rounds away from 0)
//	  -1.4 -> -1  (closer to -1)
//	  -1.5 -> -2  (tie: rounds away from 0 / larger magnitude)
//	  -2.5 -> -3  (tie: rounds away from 0 / larger magnitude)
//
// ToZeroVal (2) - IEEE 754-2008 roundTowardZero (Truncation):
//
//	Directed rounding that discards excess precision toward zero. It chooses
//	the representable value closest to, but not greater in magnitude than, the
//	exact value.
//
//	Examples (rounding to nearest integer):
//
//	  +1.1 -> +1
//	  +1.5 -> +1
//	  +1.9 -> +1
//	  -1.1 -> -1
//	  -1.5 -> -1
//	  -1.9 -> -1
//
// AwayFromZeroVal (3) - Go extension (no IEEE 754-2008 equivalent):
//
//	Directed rounding that always moves away from zero whenever the value
//	cannot be represented exactly. It selects the representable value closest
//	to, but not smaller in magnitude than, the exact value (the polar opposite
//	of truncation).
//
//	Examples (rounding to nearest integer):
//
//	  +1.1 -> +2
//	  +1.5 -> +2
//	  +1.9 -> +2
//	  -1.1 -> -2
//	  -1.5 -> -2
//	  -1.9 -> -2
//
// ToNegativeInfVal (4) - IEEE 754-2008 roundTowardNegative (Floor):
//
//	Directed rounding downward. It chooses the representable value closest to,
//	but less than or equal to (<=), the exact value.
//
//	Examples (rounding to nearest integer):
//
//	  +1.1 -> +1
//	  +1.5 -> +1
//	  +1.9 -> +1
//	  -1.1 -> -2
//	  -1.5 -> -2
//	  -1.9 -> -2
//
// ToPositiveInfVal (5) - IEEE 754-2008 roundTowardPositive (Ceiling):
//
//	Directed rounding upward. It chooses the representable value closest to,
//	but greater than or equal to (>=), the exact value.
//
//	Examples (rounding to nearest integer):
//
//	  +1.1 -> +2
//	  +1.5 -> +2
//	  +1.9 -> +2
//	  -1.1 -> -1
//	  -1.5 -> -1
//	  -1.9 -> -1
const (
	ToNearestEvenVal RoundingType = iota // == IEEE 754-2008 roundTiesToEven
	ToNearestAwayVal                     // == IEEE 754-2008 roundTiesToAway
	ToZeroVal                            // == IEEE 754-2008 roundTowardZero
	AwayFromZeroVal                      // no IEEE 754-2008 equivalent
	ToNegativeInfVal                     // == IEEE 754-2008 roundTowardNegative
	ToPositiveInfVal                     // == IEEE 754-2008 roundTowardPositive
)

// String
//
// Returns the string representation of the RoundingType value.
//
// Input Parameters
// ----------------
//
// # None
//
// Return Parameters
// -----------------
//
// string
//
//	Returns the string representation of the RoundingType value.
func (rt *RoundingType) String() string {
	switch *rt {
	case 0:
		return "ToNearestEven"
	case 1:
		return "ToNearestAway"
	case 2:
		return "ToZero"
	case 3:
		return "AwayFromZero"
	case 4:
		return "ToNegativeInf"
	case 5:
		return "ToPositiveInf"
	default:
		return "NOT SET - INVALID"
	}
}

// IsValid
//
// Returns true if the RoundingType value is one of the six supported IEEE
// 754-2008 rounding modes (ToNearestEvenVal, ToNearestAwayVal, ToZeroVal,
// AwayFromZeroVal, ToNegativeInfVal, or ToPositiveInfVal).
//
// Input Parameters
// ----------------
//
// # None
//
// Return Parameters
// -----------------
//
// bool
//
//	Returns true if the value is a valid rounding mode, otherwise false.
func (rt *RoundingType) IsValid() bool {
	switch *rt {

	case ToNearestEvenVal:
		return true

	case ToNearestAwayVal:
		return true

	case ToZeroVal:
		return true

	case AwayFromZeroVal:
		return true

	case ToNegativeInfVal:
		return true

	case ToPositiveInfVal:
		return true

	default:

	}

	return false
}

// GetBigMathRoundingMode
//
// Returns the equivalent math/big RoundingMode value for the current
// RoundingType instance. The RoundingType enumeration closely parallels the
// math/big RoundingMode type; this method maps the RoundingType value onto
// its math/big counterpart so it can be used directly with big.Float and the
// math/big rounding API.
//
// Input Parameters
// ----------------
//
// # None
//
// Return Parameters
// -----------------
//
// big.RoundingMode
//
//	Returns the equivalent math/big RoundingMode value for the RoundingType
//	instance.
//
// error
//
//	Returns nil on success, or a non-nil error if the RoundingType instance is
//	in an uninitialized / invalid state with no math/big RoundingMode
//	equivalent.
func (rt *RoundingType) GetBigMathRoundingMode() (big.RoundingMode, error) {
	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"RoundingType_0010.go",
		"RoundingType",
		"GetBigMathRoundingMode()")

	switch *rt {
	case ToNearestEvenVal:
		return big.ToNearestEven, nil

	case ToNearestAwayVal:
		return big.ToNearestAway, nil

	case ToZeroVal:
		return big.ToZero, nil

	case AwayFromZeroVal:
		return big.AwayFromZero, nil

	case ToNegativeInfVal:
		return big.ToNegativeInf, nil

	case ToPositiveInfVal:
		return big.ToPositiveInf, nil

	default:
		return 0, errGenericCfg.GetErrorNativeToFunction(
			"",
			"",
			"",
			"",
			fmt.Sprintf("RoundingType value (%d) has no math/big RoundingMode equivalent", int(*rt)),
			"RoundingType instance is in an uninitialized / invalid state")
	}
}

// SetToNearestEven
//
// Sets the RoundingType value to ToNearestEvenVal (IEEE 754-2008
// roundTiesToEven / Bankers' Rounding).
//
// Input Parameters
// ----------------
//
// # None
//
// Return Parameters
// -----------------
//
// None
func (rt *RoundingType) SetToNearestEven() {
	*rt = ToNearestEvenVal
}

// SetToNearestAway
//
// Sets the RoundingType value to ToNearestAwayVal (IEEE 754-2008
// roundTiesToAway / Commercial Rounding).
//
// Input Parameters
// ----------------
//
// # None
//
// Return Parameters
// -----------------
//
// None
func (rt *RoundingType) SetToNearestAway() {
	*rt = ToNearestAwayVal
}

// SetToZero
//
// Sets the RoundingType value to ToZeroVal (IEEE 754-2008
// roundTowardZero / Truncation).
//
// Input Parameters
// ----------------
//
// # None
//
// Return Parameters
// -----------------
//
// None
func (rt *RoundingType) SetToZero() {
	*rt = ToZeroVal
}

// SetAwayFromZero
//
// Sets the RoundingType value to AwayFromZeroVal (Go extension, no IEEE
// 754-2008 equivalent).
//
// Input Parameters
// ----------------
//
// # None
//
// Return Parameters
// -----------------
//
// None
func (rt *RoundingType) SetAwayFromZero() {
	*rt = AwayFromZeroVal
}

// SetToNegativeInf
//
// Sets the RoundingType value to ToNegativeInfVal (IEEE 754-2008
// roundTowardNegative / Floor).
//
// Input Parameters
// ----------------
//
// # None
//
// Return Parameters
// -----------------
//
// None
func (rt *RoundingType) SetToNegativeInf() {
	*rt = ToNegativeInfVal
}

// SetToPositiveInf
//
// Sets the RoundingType value to ToPositiveInfVal (IEEE 754-2008
// roundTowardPositive / Ceiling).
//
// Input Parameters
// ----------------
//
// # None
//
// Return Parameters
// -----------------
//
// None
func (rt *RoundingType) SetToPositiveInf() {
	*rt = ToPositiveInfVal
}
