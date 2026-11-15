package floatmech

import (
	"math"
	"math/big"
)

// -----------------------------------------------------------------------------
// Private Exact Dyadic / big.Int Exponentiation Kernel
// -----------------------------------------------------------------------------

// dyadicPowKernel decomposes a finite big.Float into m * 2^scaleExp2, computes
// m^exp exactly via big.Int.Exp, rescales the binary exponent, and rounds once
// to (targetPrec, targetMode).
type dyadicPowKernel struct {
	targetPrec   uint
	targetMode   big.RoundingMode
	maxExactBits uint
}

// canUseExact reports whether x.MinPrec() * exp fits within maxExactBits
// without risking uint multiplication overflow.
func (k dyadicPowKernel) canUseExact(x *big.Float, exp uint) bool {
	if exp <= 1 || k.maxExactBits == 0 {
		return false
	}
	minPrec := x.MinPrec()
	if minPrec == 0 || minPrec > k.maxExactBits {
		return false
	}
	// Overflow-safe check for (minPrec * exp <= k.maxExactBits):
	return exp <= k.maxExactBits/minPrec
}

// powExact executes exact big.Int exponentiation with pre- and post-exponent
// overflow/underflow guards. Returns (nil, false) if it cannot safely proceed.
func (k dyadicPowKernel) powExact(x *big.Float, exp uint) (*big.Float, bool) {
	if uint64(exp) > math.MaxInt64 {
		return nil, false
	}

	// Safety check 1 (Pre-Exp): detect guaranteed overflow (> big.MaxExp) or
	// underflow (< big.MinExp) before allocating big.Int.Exp.
	if earlyOut, triggered := k.checkExponentBounds(x, exp); triggered {
		return earlyOut, true
	}

	// Decompose x into exact integer mantissa m and binary scale: x == m * 2^scaleExp2.
	m, scaleExp2, ok := k.decompose(x)
	if !ok {
		return nil, false
	}

	// Compute m^exp exactly using big.Int.Exp.
	expInt := new(big.Int).SetUint64(uint64(exp))
	mPow := new(big.Int).Exp(m, expInt, nil)

	// Safety check 2: compute total binary scaling exponent in int64 so
	// negative exponents cannot wrap on 32-bit platforms.
	totalScale64, ok := mulInt64Safe(scaleExp2, int64(exp))
	if !ok {
		// Overflowed int64 during scaling; defer safely to Option 2.
		return nil, false
	}

	// Load mPow losslessly into a temporary big.Float (BitLen() >= 1 because m != 0).
	bitLen := uint(mPow.BitLen())
	f := new(big.Float).SetPrec(bitLen).SetMode(big.ToNearestEven).SetInt(mPow)

	mantOnly := new(big.Float)
	powExp2 := f.MantExp(mantOnly) // mPow == mantOnly * 2^powExp2, 0.5 <= |mantOnly| < 1

	finalExp64 := int64(powExp2) + totalScale64
	neg := x.Signbit() && exp&1 == 1

	// Safety check 3 (Post-Exp): trigger explicit Inf or 0 if the exact binary
	// exponent is too large (> big.MaxExp) or too small (< big.MinExp).
	if finalExp64 > int64(big.MaxExp) {
		return k.signedInf(neg), true
	}
	if finalExp64 < int64(big.MinExp) {
		return k.signedZero(neg), true
	}

	// Because finalExp64 is verified inside [big.MinExp, big.MaxExp], casting
	// to int is safe on both 32-bit and 64-bit targets.
	scaled := new(big.Float).SetMantExp(mantOnly, int(finalExp64))

	// Single exact-to-target rounding pass.
	return new(big.Float).SetPrec(k.targetPrec).SetMode(k.targetMode).Set(scaled), true
}

// checkExponentBounds checks whether |x|^exp is guaranteed to exceed big.MaxExp
// or fall below big.MinExp using the MantExp interval 2^(e-1) <= |x| < 2^e.
func (k dyadicPowKernel) checkExponentBounds(x *big.Float, exp uint) (*big.Float, bool) {
	e := int64(x.MantExp(nil))
	neg := x.Signbit() && exp&1 == 1

	// If e > 1, |x| >= 2^(e-1), so |x^exp| >= 2^((e-1)*exp).
	// Since every finite big.Float is strictly < 2^MaxExp, >= 2^MaxExp overflows.
	if e > 1 {
		minResExp, ok := mulInt64Safe(e-1, int64(exp))
		if !ok || minResExp >= int64(big.MaxExp) {
			return k.signedInf(neg), true
		}
	}

	// If e < 0, |x| < 2^e, so |x^exp| < 2^(e*exp).
	// Since the smallest nonzero big.Float is 0.5 * 2^MinExp == 2^(MinExp-1),
	// anything < 2^(MinExp-1) underflows to zero.
	if e < 0 {
		maxResExp, ok := mulInt64Safe(e, int64(exp))
		if !ok || maxResExp <= int64(big.MinExp)-1 {
			return k.signedZero(neg), true
		}
	}

	return nil, false
}

// decompose extracts an exact integer mantissa m and binary scale factor
// scaleExp2 from a finite nonzero x such that x == m * 2^scaleExp2.
func (k dyadicPowKernel) decompose(x *big.Float) (*big.Int, int64, bool) {
	mant := new(big.Float)
	e := x.MantExp(mant) // x == mant * 2^e, with 0.5 <= |mant| < 1.0

	p := mant.MinPrec()
	if p == 0 || p > uint(big.MaxExp) {
		return nil, 0, false
	}

	// mant's current exponent is 0; shifting to p turns all p significant bits
	// into an exact integer without risking overflow.
	mant.SetMantExp(mant, int(p))

	m, acc := mant.Int(nil)
	if acc != big.Exact {
		return nil, 0, false
	}

	scaleExp2 := int64(e) - int64(p)
	return m, scaleExp2, true
}

func (k dyadicPowKernel) signedInf(neg bool) *big.Float {
	return new(big.Float).SetPrec(k.targetPrec).SetMode(k.targetMode).SetInf(neg)
}

func (k dyadicPowKernel) signedZero(neg bool) *big.Float {
	z := new(big.Float).SetPrec(k.targetPrec).SetMode(k.targetMode)
	if neg {
		z.Neg(z)
	}
	return z
}

func mulInt64Safe(a, b int64) (int64, bool) {
	if a == 0 || b == 0 {
		return 0, true
	}
	if (a == math.MinInt64 && b == -1) || (b == math.MinInt64 && a == -1) {
		return 0, false
	}
	res := a * b
	if res/a != b {
		return 0, false
	}
	return res, true
}
