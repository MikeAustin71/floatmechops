package floatmech

import (
	"math/big"
	"sync"

	ePref "github.com/MikeAustin71/errpref/v2"
)

const defaultMaxExactBits uint = 8192

const (
	defaultFractionalDigits uint = 6
	defaultSafetyBits       uint = 8       // Extra binary cushion for decimal rounding
	maxWorkingPrec          uint = 1 << 20 // Safety cap (~1M bits) to prevent runaway allocations
)

type powStrategy uint8

const (
	strategyHeuristic powStrategy = iota // Option 2: Fixed/Heuristic Guard Bits
	strategyZiv                          // Option 3: Ziv's Adaptive Strategy
	strategyHybrid                       // Option 4: Hybrid big.Int + Option 2
)

// Decimal is an immutable, encapsulated wrapper around *big.Float (Option B).
type Decimal struct {
	mu                     sync.Mutex
	val                    *big.Float
	targetFractionalDigits uint
	prec                   uint
	mode                   big.RoundingMode
}

// CopyIn
//
// Receives an incoming 'Decimal' object (copySource)
// and copies those structure member variable values into
// the structure member variables of the host, or receiving
// 'Decimal' object.
//
// Input Parameters
// ----------------
//
//	copySource *Decimal
//	  Member variable values from this Decimal object
//	  will be copied into corresponding member variables
//	  in the host, or receiving Decimal object.
//
// Return Parameters
// -----------------
//
//	error
func (d *Decimal) CopyIn(copySource *Decimal) error {

	// NewErrorGenericConfig()
	//
	// func (egc *ErrorGenericFactory) NewErrorGenericConfig(
	//    sourceFile string,
	//    parentType string,
	//    thisFunction string) *ErrorGenericFactory
	//
	// Add identification information for this function
	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"decimal_0010.go",
		"Decimal",
		"CopyIn()")

	if d == nil {

		// GetErrorInputParameterNilPointer()
		//
		// func (egc *ErrorGenericFactory) GetErrorInputParameterNilPointer(
		//    parameter string,
		//    path string,
		//    fileName string,
		//    url string,
		//    errorContext string,
		//    errorMessage string) error
		//
		// Formats error message and error context for when nil pointers
		// are passed as Function input parameters.
		return errGenericCfg.GetErrorInputParameterNilPointer(
			"d == nil",
			"",
			"",
			"",
			"Invalid Pointer Receiver",
			"Error: Pointer Receiver 'd' is a nil pointer")
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	if copySource == nil {

		// GetErrorInputParameterNilPointer()
		//
		// func (egc *ErrorGenericFactory) GetErrorInputParameterNilPointer(
		//    parameter string,
		//    path string,
		//    fileName string,
		//    url string,
		//    errorContext string,
		//    errorMessage string) error
		//
		// Formats error message and error context for when nil pointers
		// are passed as Function input parameters.
		return errGenericCfg.GetErrorInputParameterNilPointer(
			"copySource == nil",
			"",
			"",
			"",
			"copySource is a nil pointer",
			"Error: Input parameter 'copySource' is a nil pointer")
	}

	if copySource.val == nil {

		// GetErrorInputParameterNilPointer()
		//
		// func (egc *ErrorGenericFactory) GetErrorInputParameterNilPointer(
		//    parameter string,
		//    path string,
		//    fileName string,
		//    url string,
		//    errorContext string,
		//    errorMessage string) error
		//
		// Formats error message and error context for when nil pointers
		// are passed as Function input parameters.
		return errGenericCfg.GetErrorInputParameterNilPointer(
			"copySource.val == nil",
			"",
			"",
			"",
			"copySource member variable 'val' is a nil pointer",
			"Error: Input parameter 'copySource' member variable 'val' is a nil pointer")
	}

	d.val = new(big.Float).SetMode(copySource.val.Mode()).Set(copySource.val)

	d.targetFractionalDigits = copySource.targetFractionalDigits

	d.prec = copySource.prec

	d.mode = copySource.mode

	return nil
}

// CopyOut
//
// Receives an outgoing 'Decimal' object (copyDestination)
// and copies the host Decimal structure member variable values into
// the structure member variables of this copyDestination 'Decimal'
// object.
//
// Input Parameters
// ----------------
//
//	copyDestination *Decimal
//	  Member variable values from the host (source) Decimal
//	  object will be copied into corresponding member variables
//	  of the copyDestination Decimal object.
//
// Return Parameters
// -----------------
//
//	error
func (d *Decimal) CopyOut(copyDestination *Decimal) error {

	// NewErrorGenericConfig()
	//
	// func (egc *ErrorGenericFactory) NewErrorGenericConfig(
	//    sourceFile string,
	//    parentType string,
	//    thisFunction string) *ErrorGenericFactory
	//
	// Add identification information for this function
	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"decimal_0010.go",
		"Decimal",
		"CopyOut()")

	if d == nil {

		// GetErrorInputParameterNilPointer()
		//
		// func (egc *ErrorGenericFactory) GetErrorInputParameterNilPointer(
		//    parameter string,
		//    path string,
		//    fileName string,
		//    url string,
		//    errorContext string,
		//    errorMessage string) error
		//
		// Formats error message and error context for when nil pointers
		// are passed as Function input parameters.
		return errGenericCfg.GetErrorInputParameterNilPointer(
			"d == nil",
			"",
			"",
			"",
			"Invalid Pointer Receiver",
			"Error: Pointer Receiver 'd' is a nil pointer")
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	if copyDestination == nil {

		// GetErrorInputParameterNilPointer()
		//
		// func (egc *ErrorGenericFactory) GetErrorInputParameterNilPointer(
		//    parameter string,
		//    path string,
		//    fileName string,
		//    url string,
		//    errorContext string,
		//    errorMessage string) error
		//
		// Formats error message and error context for when nil pointers
		// are passed as Function input parameters.
		return errGenericCfg.GetErrorInputParameterNilPointer(
			"copyDestination == nil",
			"",
			"",
			"",
			"		copyDestination is a nil pointer",
			"Error: Input parameter 'copyDestination' is a nil pointer")
	}

	if d.val == nil {

		// GetErrorNativeToFunction()
		//
		// func (egc *ErrorGenericFactory) GetErrorNativeToFunction(
		//    parameter string,
		//    path string,
		//    fileName string,
		//    url string,
		//    errorContext string,
		//    errorMessage string) error
		//
		// Formats error message and error context for errors
		// originating in the current function's internal logic.
		return errGenericCfg.GetErrorNativeToFunction(
			"d.val == nil",
			"",
			"",
			"",
			"Host member variable 'val' is a nil pointer",
			"This instance of type Decimal is invalid. Decimal member variable 'val' is a nil pointer.")
	}

	copyDestination.val = new(big.Float).SetMode(d.val.Mode()).Set(d.val)

	copyDestination.targetFractionalDigits = d.targetFractionalDigits

	copyDestination.prec = d.prec

	copyDestination.mode = d.mode

	return nil
}

// Clone
//
// Returns a newly created *Decimal object which is an
// exact clone of the host Decimal object.
//
// Input Parameters
// ----------------
//
//	None
//
// Return Parameters
// -----------------
//
//	*Decimal
//	  Pointer to a newly created Decimal object containing duplicate
//	  member variable values to those in the host Decimal object.
//
//	error
func (d *Decimal) Clone() (*Decimal, error) {

	// NewErrorGenericConfig()
	//
	// func (egc *ErrorGenericFactory) NewErrorGenericConfig(
	//    sourceFile string,
	//    parentType string,
	//    thisFunction string) *ErrorGenericFactory
	//
	// Add identification information for this function
	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"decimal_0010.go",
		"Decimal",
		"Clone()")

	if d == nil {

		// GetErrorInputParameterNilPointer()
		//
		// func (egc *ErrorGenericFactory) GetErrorInputParameterNilPointer(
		//    parameter string,
		//    path string,
		//    fileName string,
		//    url string,
		//    errorContext string,
		//    errorMessage string) error
		//
		// Formats error message and error context for when nil pointers
		// are passed as Function input parameters.
		return &Decimal{}, errGenericCfg.GetErrorInputParameterNilPointer(
			"d == nil",
			"",
			"",
			"",
			"Invalid Pointer Receiver",
			"Error: Pointer Receiver 'd' is a nil pointer")
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	if d.val == nil {

		// GetErrorNativeToFunction()
		//
		// func (egc *ErrorGenericFactory) GetErrorNativeToFunction(
		//    parameter string,
		//    path string,
		//    fileName string,
		//    url string,
		//    errorContext string,
		//    errorMessage string) error
		//
		// Formats error message and error context for errors
		// originating in the current function's internal logic.
		return &Decimal{}, errGenericCfg.GetErrorNativeToFunction(
			"d.val == nil",
			"",
			"",
			"",
			"Host member variable 'val' is a nil pointer",
			"This instance of type Decimal is invalid. Decimal member variable 'val' is a nil pointer.")
	}

	clonedDecimal := new(Decimal)

	clonedDecimal.val = new(big.Float).SetMode(d.val.Mode()).Set(d.val)

	clonedDecimal.targetFractionalDigits = d.targetFractionalDigits

	clonedDecimal.prec = d.prec

	clonedDecimal.mode = d.mode

	return clonedDecimal, nil
}

// NewFromFloatNumStr
//
// Constructs an immutable Decimal sized to preserve at least
// targetFractionalDigits decimal places for the value's magnitude.
//
// The input numeric string value is padded with a fixed number
// of trailing '0' characters prior to parsing, raising the
// parsed value's magnitude so that the derived working
// precision retains the requested number of fractional digits.
//
// Input Parameters
// ----------------
//
//	vNumStr string
//	  Numeric string value used to construct the Decimal. If
//	  'vNumStr' is an empty string, an error is returned.
//
//	targetFractionalDigits uint
//	  Number of base-10 fractional digits to be preserved by the
//	  constructed Decimal.
//
//	mode big.RoundingMode
//	  Rounding mode applied when the parsed value is rounded to
//	  the working precision.
//
// Return Parameters
// -----------------
//
//	*Decimal
//	  Pointer to the newly constructed Decimal object.
//
//	error
func (d *Decimal) NewFromFloatNumStr(vNumStr string, targetFractionalDigits uint, mode big.RoundingMode) (*Decimal, error) {

	// NewErrorGenericConfig()
	//
	// func (egc *ErrorGenericFactory) NewErrorGenericConfig(
	//    sourceFile string,
	//    parentType string,
	//    thisFunction string) *ErrorGenericFactory
	//
	// Add identification information for this function
	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"decimal_0001000.go",
		"Decimal",
		"NewFromFloatNumStr()")

	if d == nil {

		// GetErrorInputParameterNilPointer()
		//
		// func (egc *ErrorGenericFactory) GetErrorInputParameterNilPointer(
		//    parameter string,
		//    path string,
		//    fileName string,
		//    url string,
		//    errorContext string,
		//    errorMessage string) error
		//
		// Formats error message and error context for when nil pointers
		// are passed as Function input parameters.
		return &Decimal{}, errGenericCfg.GetErrorInputParameterNilPointer(
			"d == nil",
			"",
			"",
			"",
			"Invalid Pointer Receiver",
			"Error: Pointer Receiver 'd' is a nil pointer")
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	constructedDecimal, err := new(numStrConstructorMechanics).constructFromNumStr(
		vNumStr,
		targetFractionalDigits,
		mode)

	if err != nil {

		// GetErrorFromFunctionCall()
		//
		// func (egc *ErrorGenericFactory) GetErrorFromFunctionCall(
		//    calledFunctionType string,
		//    calledFunction string,
		//    errorContext string,
		//    errorMessage string,
		//    returnedCalledFunctionError error) error
		//
		// Adds error context to errors returned by other
		// Functions.
		err2 := errGenericCfg.GetErrorFromFunctionCall(
			"numStrConstructorMechanics",
			"constructFromNumStr()",
			"",
			"",
			err)

		return &Decimal{}, err2
	}

	return constructedDecimal, nil
}

// TargetFractionalDigits returns the requested base-10 digits after the decimal point.
func (d *Decimal) TargetFractionalDigits() uint {

	if d == nil {
		return 0
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	return d.targetFractionalDigits
}

// ActualFractionalDigits derives how many base-10 fractional digits are guaranteed
// by the underlying big.Float's binary precision at its current magnitude.
func (d *Decimal) ActualFractionalDigits() uint {

	if d.val == nil {
		return 0
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	if d.val.Sign() == 0 || d.val.IsInf() {
		return d.targetFractionalDigits
	}

	return defaultConverter().toFractionalDigits(d.prec, d.val.MantExp(nil))
}

// Prec returns the underlying big.Float total mantissa precision in bits.
func (d *Decimal) Prec() uint {

	if d == nil {
		return 0
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	return d.prec
}

// IsValid
//
// Tests the internal value stored in structure element 'val'.
// If 'val' is a nil pointer, the method returns 'false'.
// If 'val' is not nil, the method returns 'true'.
//
// Input Parameters
// ----------------
//
//	None
//
// Return Parameters
// -----------------
//
//	bool
//	  Returns 'false' if the internal value 'val' is a nil pointer.
//	  Returns 'true' if the internal value 'val' is not nil.
func (d *Decimal) IsValid() bool {

	if d == nil {
		return false
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	return d.val != nil
}

// IsError
//
// Tests the internal value stored in structure element 'val'.
// If 'val' is a nil pointer, the method returns an error formatted
// via GetErrorNativeToFunction. If 'val' is not nil, the method
// returns nil.
//
// Input Parameters
// ----------------
//
//	None
//
// Return Parameters
// -----------------
//
//	error
//	  Returns an error if the internal value 'val' is a nil pointer.
//	  Returns nil if the internal value 'val' is not nil.
func (d *Decimal) IsError() error {

	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"decimal_0001000.go",
		"Decimal",
		"IsError()")

	if d == nil {

		return errGenericCfg.GetErrorNativeToFunction(
			"d == nil",
			"",
			"",
			"",
			"Invalid Pointer Receiver",
			"This instance of type Decimal is invalid. Decimal member variable 'val' is a nil pointer.")
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	if d.val == nil {

		return errGenericCfg.GetErrorNativeToFunction(
			"val == nil",
			"",
			"",
			"",
			"Invalid member variable",
			"This instance of type Decimal is invalid. Decimal member variable 'val' is a nil pointer.")
	}

	return nil
}

// Pow returns d^exp using Option 2 (Fixed/Heuristic Guard Bits).
func (d *Decimal) Pow(exp uint) *Decimal {

	if d == nil {
		return &Decimal{}
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	return d.powWithStrategy(exp, strategyHeuristic)
}

// PowExact returns d^exp using Option 3 (Ziv's Adaptive Strategy).
func (d *Decimal) PowExact(exp uint) *Decimal {

	if d == nil {
		return &Decimal{}
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	return d.powWithStrategy(exp, strategyZiv)
}

// PowHybrid returns d^exp using Option 4 (Exact big.Int for small
// MinPrec()*exp, falling back to Option 2 for larger exponents).
func (d *Decimal) PowHybrid(exp uint) *Decimal {

	if d == nil {
		return &Decimal{}
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	return d.powWithStrategy(exp, strategyHybrid)
}

// String formats the Decimal in fixed-point notation to targetFractionalDigits.
func (d *Decimal) String() string {

	if d == nil {
		return ""
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	if d.val == nil {
		return "0"
	}

	if d.val.IsInf() {
		return d.val.String()
	}

	return d.val.Text('f', int(d.targetFractionalDigits))
}

func (d *Decimal) powWithStrategy(exp uint, strategy powStrategy) *Decimal {

	if d == nil {
		return &Decimal{}
	}

	if d.val == nil {
		zeroBase, err := new(Decimal).NewFromFloatNumStr("0", defaultFractionalDigits, big.ToNearestEven)

		if err != nil {
			return &Decimal{}
		}

		return zeroBase.powWithStrategy(exp, strategy)
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
