package floatmech

import (
	"math"
	"math/big"
	"math/bits"
)

// -----------------------------------------------------------------------------
// Private Reusable Calculation Engine
// -----------------------------------------------------------------------------

type floatEngine struct {
	prec    uint
	mode    big.RoundingMode
	maxPrec uint
}

func (e floatEngine) estimatePowExp2(x *big.Float, exp uint) int {
	if exp == 0 || x.Sign() == 0 || x.IsInf() {
		return 1
	}

	mant := new(big.Float)
	exp2 := x.MantExp(mant)
	mant.Abs(mant)

	m64, _ := mant.Float64()
	est := float64(exp) * (math.Log2(m64) + float64(exp2))

	if est <= 0 {
		return 0
	}
	if est >= float64(big.MaxExp) {
		return big.MaxExp
	}
	return int(math.Ceil(est)) + 1
}

// powSpecialCase handles exp in {0, 1}, ±0, ±Inf, and ±1 in one shared place.
func (e floatEngine) powSpecialCase(x *big.Float, exp uint) (*big.Float, bool) {
	switch exp {
	case 0:
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(1), true
	case 1:
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).Set(x), true
	}

	if x.Sign() == 0 {
		z := new(big.Float).SetPrec(e.prec).SetMode(e.mode)
		if x.Signbit() && exp&1 == 1 {
			z.Neg(z)
		}
		return z, true
	}

	if x.IsInf() {
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInf(x.Signbit() && exp&1 == 1), true
	}

	if x.Cmp(big.NewFloat(1)) == 0 {
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(1), true
	}
	if x.Cmp(big.NewFloat(-1)) == 0 {
		if exp&1 == 1 {
			return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(-1), true
		}
		return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInt64(1), true
	}

	return nil, false
}

// powHeuristic computes x^exp using Option 2 (Fixed/Heuristic Guard Bits).
func (e floatEngine) powHeuristic(x *big.Float, exp uint) *big.Float {
	if out, ok := e.powSpecialCase(x, exp); ok {
		return out
	}

	guardBits := uint(bits.Len(exp)) + 32
	workPrec := e.prec + guardBits
	if workPrec < e.prec || (e.maxPrec > 0 && workPrec > e.maxPrec) {
		workPrec = e.maxPrec
	}

	res := new(big.Float).SetPrec(workPrec).SetMode(big.ToNearestEven).SetInt64(1)
	base := new(big.Float).SetPrec(workPrec).SetMode(big.ToNearestEven).Set(x)

	origExp := exp
	for exp > 0 {
		if exp&1 == 1 {
			res.Mul(res, base)
		}
		if exp > 1 {
			base.Mul(base, base)
			if base.IsInf() {
				return new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInf(x.Signbit() && origExp&1 == 1)
			}
			if base.Sign() == 0 {
				z := new(big.Float).SetPrec(e.prec).SetMode(e.mode)
				if x.Signbit() && origExp&1 == 1 {
					z.Neg(z)
				}
				return z
			}
		}
		exp >>= 1
	}

	return new(big.Float).SetPrec(e.prec).SetMode(e.mode).Set(res)
}

// powZiv computes x^exp using Option 3 (Ziv's Adaptive Strategy), iteratively
// increasing working precision until the final rounding to e.prec is provably
// unambiguous.
func (e floatEngine) powZiv(x *big.Float, exp uint) *big.Float {
	if out, ok := e.powSpecialCase(x, exp); ok {
		return out
	}

	// Ensure basePrec covers both target precision and x's exact mantissa bits
	// so the initial load of base = x introduces zero quantization error.
	basePrec := e.prec
	if minP := x.MinPrec(); minP > basePrec {
		basePrec = minP
	}

	// Start with enough guard bits so almost all inputs succeed on Pass 1.
	guardBits := uint(bits.Len(exp)) + 32

	for {
		workPrec := basePrec + guardBits
		capped := false
		if workPrec < basePrec || (e.maxPrec > 0 && workPrec >= e.maxPrec) {
			workPrec = e.maxPrec
			capped = true
		}

		bounder := newZivBounder(e.prec, e.mode, workPrec)
		res, errUnits, earlyOut := e.runZivPass(x, exp, bounder)
		if earlyOut != nil {
			return earlyOut
		}

		// Verify whether rounding res to (e.prec, e.mode) is unambiguous.
		if rounded, ok := bounder.tryRound(res, errUnits); ok {
			return rounded
		}

		// Safeguard: stop if we hit the maximum precision cap.
		if capped {
			return new(big.Float).SetPrec(e.prec).SetMode(e.mode).Set(res)
		}

		// Double guard bits for the next Ziv refinement pass.
		guardBits *= 2
	}
}

// runZivPass executes one binary-exponentiation pass at b.workPrec while
// accumulating a deterministic upper bound on relative rounding error.
func (e floatEngine) runZivPass(x *big.Float, exp uint, b zivBounder) (*big.Float, float64, *big.Float) {
	res := new(big.Float).SetPrec(b.workPrec).SetMode(big.ToNearestEven).SetInt64(1)
	base := new(big.Float).SetPrec(b.workPrec).SetMode(big.ToNearestEven).Set(x)

	resErr := 0.0
	baseErr := 0.0
	if base.Acc() != big.Exact {
		// Only happens if x.MinPrec() exceeded maxWorkingPrec.
		baseErr = 1.0
	}

	origExp := exp
	for exp > 0 {
		if exp&1 == 1 {
			res.Mul(res, base)
			resErr = b.propagateMulErr(resErr, baseErr, res.Acc())
		}
		if exp > 1 {
			base.Mul(base, base)
			baseErr = b.propagateMulErr(baseErr, baseErr, base.Acc())

			if base.IsInf() {
				inf := new(big.Float).SetPrec(e.prec).SetMode(e.mode).SetInf(x.Signbit() && origExp&1 == 1)
				return nil, 0, inf
			}
			if base.Sign() == 0 {
				z := new(big.Float).SetPrec(e.prec).SetMode(e.mode)
				if x.Signbit() && origExp&1 == 1 {
					z.Neg(z)
				}
				return nil, 0, z
			}
		}
		exp >>= 1
	}

	return res, resErr, nil
}

// powHybrid computes x^exp using exact big.Int exponentiation when
// x.MinPrec() * exp <= defaultMaxExactBits, safely falling back to Option 2
// (powHeuristic) when the exponent or bit-growth is too large.
func (e floatEngine) powHybrid(x *big.Float, exp uint) *big.Float {
	// 1. Handle exp in {0, 1}, ±0, ±Inf, and ±1 immediately.
	if out, ok := e.powSpecialCase(x, exp); ok {
		return out
	}

	// 2. Attempt exact big.Int calculation if within the bit-budget threshold.
	kernel := dyadicPowKernel{
		targetPrec:   e.prec,
		targetMode:   e.mode,
		maxExactBits: defaultMaxExactBits,
	}

	if kernel.canUseExact(x, exp) {
		if out, ok := kernel.powExact(x, exp); ok {
			return out
		}
	}

	// 3. Fallback to Option 2 (Fixed/Heuristic Guard Bits) for large exponents
	// or if exact decomposition was declined.
	return e.powHeuristic(x, exp)
}
