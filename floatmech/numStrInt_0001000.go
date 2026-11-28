package floatmech

type NumStrIntDto struct {
	NumType          NumericValueType
	NumIntegerDigits uint
	IntegerRunes     []rune
	NumericSign      NumericSignType
	ValueIsZero      bool
	ValueIsOne       bool
}
