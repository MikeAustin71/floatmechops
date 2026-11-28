package floatmech

// SciNotationNum
//
// # Example Scientific Notation
//
//	 2.652e+8
//
//	 significand 		= '2.652'
//	 significand integer digits = '2'
//		mantissa				= significand factional digits = '.652'
//	 exponent    		= '8'  (10^8)
//		mantissaLength	= length of fractional digits displayed in scientific notation.
type SciNotationNum struct {
	significand NumStrFloatDto // The significand consists of the leading integer and
	//	fractional digits of the scientific notation.
	exponent         NumStrIntDto // The exponent portion of the scientific notation string
	exponentChar     rune         // 	defaults to 'e'. May be customized to 'E'
	decimalSeparator rune         // The decimal separator used to separate integer and
	// 	fractional digits in the significand. The default is
	// 	the standard USA decimal separator, the decimal point ('.').
	mantissaLength uint // The length of the fractional digits in
	// 	the significand which will be displayed
	// 	when SciNotationNum.GetSciNotationStr()
	// 	is called.
	exponentUsesLeadingPlus bool // If true, positive exponent values are
	// 	prefixed with a leading plus (+) sign.
	//  '2.652e+8'
}
