package floatmech

import "math"

// -----------------------------------------------------------------------------
// Low-Level Private Precision Converter
// -----------------------------------------------------------------------------

const (
	log2Of10 = math.Ln10 / math.Ln2
	log10Of2 = math.Ln2 / math.Ln10
)

type precConverter struct {
	safetyBits uint
	maxPrec    uint
}

func defaultConverter() precConverter {
	return precConverter{
		safetyBits: defaultSafetyBits,
		maxPrec:    maxWorkingPrec,
	}
}

func (c precConverter) toPrec(targetFractionalDigits uint, exp2 int) uint {
	fracBits := uint(math.Ceil(float64(targetFractionalDigits) * log2Of10))

	var intBits uint
	if exp2 > 0 {
		intBits = uint(exp2)
	}

	total := intBits + fracBits
	if total < intBits || total+c.safetyBits < total {
		return c.maxPrec
	}
	total += c.safetyBits

	if total == 0 {
		total = 1
	}
	if c.maxPrec > 0 && total > c.maxPrec {
		return c.maxPrec
	}
	return total
}

func (c precConverter) toFractionalDigits(prec uint, exp2 int) uint {
	if prec <= c.safetyBits {
		return 0
	}
	usablePrec := prec - c.safetyBits

	var fracBits uint
	if exp2 > 0 {
		intBits := uint(exp2)
		if usablePrec <= intBits {
			return 0
		}
		fracBits = usablePrec - intBits
	} else {
		fracBits = usablePrec + uint(-exp2)
	}

	return uint(math.Floor(float64(fracBits) * log10Of2))
}
