package floatmech

import (
	"math"
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

	d.val = copySource.val
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
			"copyDestination is a nil pointer",
			"Error: Input parameter 'copyDestination' is a nil pointer")
	}

	copyDestination.val = d.val
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

	clonedDecimal := new(Decimal)

	clonedDecimal.val = d.val
	clonedDecimal.targetFractionalDigits = d.targetFractionalDigits
	clonedDecimal.prec = d.prec
	clonedDecimal.mode = d.mode

	return clonedDecimal, nil
}

// NewFromFloat64 constructs an immutable Decimal sized to preserve at least
// targetFractionalDigits decimal places for v's magnitude.
func (d *Decimal) NewFromFloat64(v float64, targetFractionalDigits uint, mode big.RoundingMode) *Decimal {

	if d.val == nil {
		return &Decimal{}
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	conv := defaultConverter()

	var exp2 int
	if v != 0 && !math.IsInf(v, 0) && !math.IsNaN(v) {
		_, exp2 = math.Frexp(v)
	}

	prec := conv.toPrec(targetFractionalDigits, exp2)

	// Maximum fractional precision for type: float64
	if prec < 53 {
		prec = 53
	}

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
		return new(Decimal).NewFromFloat64(0, defaultFractionalDigits, big.ToNearestEven).powWithStrategy(exp, strategy)
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
