package floatmech

// NumericValueType
//
// Represents the fundamental kind of a numeric value as either IntegerValue or
// FloatingPointValue.
type NumericValueType int

// IntegerValue, FloatingPointValue
//
// Enumerated constants representing the possible kinds of a numeric value.
const (
	IntegerValue NumericValueType = iota + 1
	FloatingPointValue
)

// NumericValueTypeNames
//
// Maps NumericValueType values to their string representations.
// Index 0 represents the uninitialized / invalid state.
var NumericValueTypeNames = [...]string{
	"INVALID",
	"IntegerValue",
	"FloatingPointValue",
}

// String
//
// Returns the string representation of the NumericValueType value.
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
//	Returns the string representation of the NumericValueType value, or
//	"INVALID" for an out-of-range value.
func (nvt *NumericValueType) String() string {
	if int(*nvt) < 1 || int(*nvt) >= len(NumericValueTypeNames) {
		return NumericValueTypeNames[0]
	}

	return NumericValueTypeNames[*nvt]
}

// IsInteger
//
// Returns true if the NumericValueType value is IntegerValue.
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
//	Returns true if the value is IntegerValue, otherwise false.
func (nvt *NumericValueType) IsInteger() bool {
	if *nvt == IntegerValue {
		return true
	}

	return false
}

// IsFloatingPoint
//
// Returns true if the NumericValueType value is FloatingPointValue.
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
//	Returns true if the value is FloatingPointValue, otherwise false.
func (nvt *NumericValueType) IsFloatingPoint() bool {
	if *nvt == FloatingPointValue {
		return true
	}

	return false
}

// IsValid
//
// Returns true if the NumericValueType value is IntegerValue or FloatingPointValue.
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
//	Returns true if the value is a valid numeric value type, otherwise false.
func (nvt *NumericValueType) IsValid() bool {
	switch *nvt {

	case IntegerValue:
		return true

	case FloatingPointValue:
		return true

	default:

	}

	return false
}

// SetInteger
//
// Sets the NumericValueType value to IntegerValue.
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
func (nvt *NumericValueType) SetInteger() {
	*nvt = IntegerValue
}

// SetFloatingPoint
//
// Sets the NumericValueType value to FloatingPointValue.
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
func (nvt *NumericValueType) SetFloatingPoint() {
	*nvt = FloatingPointValue
}
