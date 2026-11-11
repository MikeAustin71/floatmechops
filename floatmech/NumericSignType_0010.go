package floatmech

type NumericSignType int

const (
	NegativeValue NumericSignType = iota - 1
	ZeroValue
	PositiveValue
)

var NumericSignTypeNames = [...]string{
	"Negative",
	"Zero",
	"Positive",
}

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

func (nst *NumericSignType) IsNegative() bool {

	if *nst == NegativeValue {
		return true
	}

	return false
}

func (nst *NumericSignType) IsPositive() bool {

	if *nst == PositiveValue {
		return true
	}

	return false
}

func (nst *NumericSignType) IsZero() bool {

	if *nst == ZeroValue {
		return true
	}

	return false
}

func (nst *NumericSignType) IsNegativeOrZero() bool {

	if *nst == NegativeValue || *nst == ZeroValue {
		return true
	}
	return false
}

func (nst *NumericSignType) IsPositiveOrZero() bool {
	if *nst == PositiveValue || *nst == ZeroValue {
		return true
	}
	return false
}

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

func (nst *NumericSignType) SetPositive() {
	*nst = PositiveValue
}

func (nst *NumericSignType) SetZero() {
	*nst = ZeroValue
}

func (nst *NumericSignType) SetNegative() {
	*nst = NegativeValue
}
