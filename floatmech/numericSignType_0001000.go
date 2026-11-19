package floatmech

// NumericSignType
//
// Represents the sign of a numeric value as Negative, Zero, or Positive.
type NumericSignType int

// NegativeValue, ZeroValue, PositiveValue
//
// Enumerated constants representing the possible sign states of a numeric value.
const (
	NegativeValue NumericSignType = iota - 1
	ZeroValue
	PositiveValue
)

// NumericSignTypeNames
//
// Maps NumericSignType values to their string representations.
var NumericSignTypeNames = [...]string{
	"Negative",
	"Zero",
	"Positive",
}

// String
//
// Returns the string representation of the NumericSignType value.
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
//	Returns the string representation of the NumericSignType value.
func (nst *NumericSignType) String() string {
	switch *nst {
	case -1:
		return "Negative"
	case 0:
		return "Zero"
	case 1:
		return "Positive"
	default:
		return "NOT SET - INVALID"
	}
}

// IsNegative
//
// Returns true if the NumericSignType value is NegativeValue.
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
//	Returns true if the value is NegativeValue, otherwise false.
func (nst *NumericSignType) IsNegative() bool {

	if *nst == NegativeValue {
		return true
	}

	return false
}

// IsPositive
//
// Returns true if the NumericSignType value is PositiveValue.
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
//	Returns true if the value is PositiveValue, otherwise false.
func (nst *NumericSignType) IsPositive() bool {

	if *nst == PositiveValue {
		return true
	}

	return false
}

// IsZero
//
// Returns true if the NumericSignType value is ZeroValue.
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
//	Returns true if the value is ZeroValue, otherwise false.
func (nst *NumericSignType) IsZero() bool {

	if *nst == ZeroValue {
		return true
	}

	return false
}

// IsNegativeOrZero
//
// Returns true if the NumericSignType value is NegativeValue or ZeroValue.
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
//	Returns true if the value is NegativeValue or ZeroValue, otherwise false.
func (nst *NumericSignType) IsNegativeOrZero() bool {

	if *nst == NegativeValue || *nst == ZeroValue {
		return true
	}
	return false
}

// IsPositiveOrZero
//
// Returns true if the NumericSignType value is PositiveValue or ZeroValue.
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
//	Returns true if the value is PositiveValue or ZeroValue, otherwise false.
func (nst *NumericSignType) IsPositiveOrZero() bool {
	if *nst == PositiveValue || *nst == ZeroValue {
		return true
	}
	return false
}

// IsValid
//
// Returns true if the NumericSignType value is NegativeValue, ZeroValue, or PositiveValue.
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
//	Returns true if the value is a valid sign state, otherwise false.
func (nst *NumericSignType) IsValid() bool {

	switch *nst {

	case NegativeValue:
		return true

	case ZeroValue:
		return true

	case PositiveValue:
		return true

	default:

	}

	return false
}

// SetPositive
//
// Sets the NumericSignType value to PositiveValue.
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
func (nst *NumericSignType) SetPositive() {
	*nst = PositiveValue
}

// SetZero
//
// Sets the NumericSignType value to ZeroValue.
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
func (nst *NumericSignType) SetZero() {
	*nst = ZeroValue
}

// SetNegative
//
// Sets the NumericSignType value to NegativeValue.
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
func (nst *NumericSignType) SetNegative() {
	*nst = NegativeValue
}
