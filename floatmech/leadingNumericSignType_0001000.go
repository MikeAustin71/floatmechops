package floatmech

type LeadingNumericSignType int

const (
	LeadingMinusSign LeadingNumericSignType = iota + 1
	LeadingPlusSign
	ImpliedPlusSign // There is no sign present but a plus sign is implied
	LeadingTrailingParenSigns
	InvalidLeadingSign
)

// LeadingNumericSignTypeNames
//
// Maps LeadingNumericSignType values to their string representations.
var LeadingNumericSignTypeNames = [...]string{
	"LeadingMinus",
	"LeadingPlus",
	"ImpliedPlus",
	"LeadingTrailingParen",
	"InvalidLeadingSign",
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
func (lnst *LeadingNumericSignType) String() string {

	switch *lnst {

	case LeadingMinusSign:
		return "LeadingMinus"

	case LeadingPlusSign:
		return "LeadingPlus"

	case ImpliedPlusSign:
		return "ImpliedPlus"

	case LeadingTrailingParenSigns:
		return "LeadingTrailingParen"

	default:
		return "InvalidLeadingSign"
	}

}

func (lnst *LeadingNumericSignType) IsLeadingMinus() bool {

	if *lnst == LeadingMinusSign {
		return true
	}

	return false
}

func (lnst *LeadingNumericSignType) IsLeadingPlus() bool {

	if *lnst == LeadingPlusSign {
		return true
	}

	return false
}

func (lnst *LeadingNumericSignType) IsImpliedPlusSign() bool {

	if *lnst == ImpliedPlusSign {
		return true
	}

	return false
}

func (lnst *LeadingNumericSignType) IsLeadingTrailingParen() bool {

	if *lnst == LeadingTrailingParenSigns {
		return true
	}

	return false
}

func (lnst *LeadingNumericSignType) IsInvalidLeadingSign() bool {

	switch *lnst {

	case LeadingMinusSign:
		return false

	case LeadingPlusSign:
		return false

	case ImpliedPlusSign:
		return false

	case LeadingTrailingParenSigns:
		return false

	default:
		return true
	}
}

func (lnst *LeadingNumericSignType) SetLeadingMinusSign() {
	*lnst = LeadingMinusSign
}

func (lnst *LeadingNumericSignType) SetLeadingPlusSign() {
	*lnst = LeadingPlusSign
}

func (lnst *LeadingNumericSignType) SetImpliedPlusSign() {
	*lnst = ImpliedPlusSign
}

func (lnst *LeadingNumericSignType) SetLeadingTrailingParenSigns() {
	*lnst = LeadingTrailingParenSigns
}

func (lnst *LeadingNumericSignType) SetInvalidLeadingSign() {
	*lnst = InvalidLeadingSign
}
