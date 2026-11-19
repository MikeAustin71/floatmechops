package floatmech

import (
	"math"
	"math/big"
)

// -----------------------------------------------------------------------------
// Private Ziv Error-Bounding & Interval Verification Type
// -----------------------------------------------------------------------------

// zivBounder tracks relative rounding error in units of u = 2^-workPrec and
// verifies whether rounding a working-precision result to (targetPrec, targetMode)
// is provably unambiguous.
type zivBounder struct {
	targetPrec   uint
	targetMode   big.RoundingMode
	workPrec     uint
	unitRoundoff float64 // 2^-workPrec (upper bound on relative error for ToNearestEven)
}

func newZivBounder(targetPrec uint, targetMode big.RoundingMode, workPrec uint) zivBounder {
	return zivBounder{
		targetPrec:   targetPrec,
		targetMode:   targetMode,
		workPrec:     workPrec,
		unitRoundoff: math.Ldexp(1.0, -int(workPrec)),
	}
}

// propagateMulErr computes the upper bound on relative error (in units of u)
// after multiplying two values with error bounds errA and errB.
func (z zivBounder) propagateMulErr(errA, errB float64, acc big.Accuracy) float64 {
	// (1 + errA*u)*(1 + errB*u) - 1 = (errA + errB + errA*errB*u)*u
	prop := errA + errB + errA*errB*z.unitRoundoff
	if acc == big.Exact {
		return prop
	}
	// Multiply by (1 + 1.1*u) for the new rounding step.
	// The 1.1 factor is a safety margin to account for accumulated drift in binary exponentiation.
	return prop + 1.1 + prop*z.unitRoundoff
}

// tryRound checks if the true mathematical value in [res - delta, res + delta]
// rounds identically at both endpoints under (targetPrec, targetMode).
func (z zivBounder) tryRound(res *big.Float, errUnits float64) (*big.Float, bool) {
	// Fast path: if every operation was Exact (e.g., exact powers or exact ties
	// like 1.5^2 = 2.25), res has zero error and rounds immediately.
	if errUnits == 0 {
		return new(big.Float).SetPrec(z.targetPrec).SetMode(z.targetMode).Set(res), true
	}

	if res.IsInf() {
		return new(big.Float).SetPrec(z.targetPrec).SetMode(z.targetMode).Set(res), true
	}

	// If accumulated error exceeds 0.25 ULP of targetPrec, do not bother with
	// the interval check; signal powZiv to increase guardBits immediately.
	var guardBits int
	if z.workPrec > z.targetPrec {
		guardBits = int(z.workPrec - z.targetPrec)
	}
	if errUnits >= math.Ldexp(0.25, guardBits) {
		return nil, false
	}

	// Apply a 10% conservative multiplier to errUnits (increased from 1% per error report)
	// to cover both float64 rounding and converting relative error |Y_true - res| / |Y_true| to |res|.
	safeErrUnits := errUnits * 1.1

	// Build eps >= safeErrUnits * 2^-workPrec rounded upward (ToPositiveInf).
	eps := new(big.Float).SetPrec(64).SetMode(big.ToPositiveInf).SetFloat64(safeErrUnits)
	eps.SetMantExp(eps, eps.MantExp(nil)-int(z.workPrec))

	// delta >= |res| * eps, rounded upward.
	delta := new(big.Float).SetPrec(64).SetMode(big.ToPositiveInf).Abs(res)
	delta.Mul(delta, eps)

	// Construct [lo, hi] strictly enclosing the true value Y_true.
	lo := new(big.Float).SetPrec(z.workPrec).SetMode(big.ToNegativeInf).Sub(res, delta)
	hi := new(big.Float).SetPrec(z.workPrec).SetMode(big.ToPositiveInf).Add(res, delta)

	// Round both endpoints to target precision and target rounding mode.
	rLo := new(big.Float).SetPrec(z.targetPrec).SetMode(z.targetMode).Set(lo)
	rHi := new(big.Float).SetPrec(z.targetPrec).SetMode(z.targetMode).Set(hi)

	// By monotonicity of rounding, rLo == rHi guarantees round(Y_true) == rLo.
	if rLo.Cmp(rHi) == 0 {
		return rLo, true
	}
	return nil, false
}
