package floatmech

import (
	"math"
	"math/big"
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
	val                    *big.Float
	targetFractionalDigits uint
	prec                   uint
	mode                   big.RoundingMode
}

// NewFromFloat64 constructs an immutable Decimal sized to preserve at least
// targetFractionalDigits decimal places for v's magnitude.
func (d *Decimal) NewFromFloat64(v float64, targetFractionalDigits uint, mode big.RoundingMode) *Decimal {
	conv := defaultConverter()

	var exp2 int
	if v != 0 && !math.IsInf(v, 0) && !math.IsNaN(v) {
		_, exp2 = math.Frexp(v)
	}

	prec := conv.toPrec(targetFractionalDigits, exp2)
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
	return d.targetFractionalDigits
}

// ActualFractionalDigits derives how many base-10 fractional digits are guaranteed
// by the underlying big.Float's binary precision at its current magnitude.
func (d *Decimal) ActualFractionalDigits() uint {

	/*	if d == nil || d.val == nil || d.val.Sign() == 0 || d.val.IsInf() {
			return d.targetFractionalDigits
		}
	*/

	if d.val == nil || d.val.Sign() == 0 || d.val.IsInf() {
		return d.targetFractionalDigits
	}

	return defaultConverter().toFractionalDigits(d.prec, d.val.MantExp(nil))
}

// Prec returns the underlying big.Float total mantissa precision in bits.
func (d *Decimal) Prec() uint {
	if d == nil {
		return 0
	}
	return d.prec
}

// Pow returns d^exp using Option 2 (Fixed/Heuristic Guard Bits).
func (d *Decimal) Pow(exp uint) *Decimal {
	return d.powWithStrategy(exp, strategyHeuristic)
}

// PowExact returns d^exp using Option 3 (Ziv's Adaptive Strategy).
func (d *Decimal) PowExact(exp uint) *Decimal {
	return d.powWithStrategy(exp, strategyZiv)
}

// PowHybrid returns d^exp using Option 4 (Exact big.Int for small
// MinPrec()*exp, falling back to Option 2 for larger exponents).
func (d *Decimal) PowHybrid(exp uint) *Decimal {
	return d.powWithStrategy(exp, strategyHybrid)
}

// String formats the Decimal in fixed-point notation to targetFractionalDigits.
func (d *Decimal) String() string {
	if d == nil || d.val == nil {
		return "0"
	}
	if d.val.IsInf() {
		return d.val.String()
	}
	return d.val.Text('f', int(d.targetFractionalDigits))
}

func (d *Decimal) powWithStrategy(exp uint, strategy powStrategy) *Decimal {
	if d == nil || d.val == nil {
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
