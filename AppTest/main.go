package main

import (
	"fmt"
	"math/big"
)

// exactPow returns the exact decimal expansion of base^exp, scaled to tfd
// fractional digits, rounded half-away-from-zero (big.ToNearestAway) using
// pure integer arithmetic on base written as an exact integer / 10^fracDigits.
func exactPow(baseStr string, exp uint, tfd uint) string {
	var numStr, fracStr string
	dot := false
	for i := 0; i < len(baseStr); i++ {
		if baseStr[i] == '.' {
			dot = true
			continue
		}
		if dot {
			fracStr += string(baseStr[i])
		} else {
			numStr += string(baseStr[i])
		}
	}
	if !dot {
		fracStr = ""
	}

	// Combine integer and fractional parts into one exact integer numerator
	// over denominator 10^fracDigits. e.g. "6.9857" -> 69857 / 10^4.
	wholeStr := numStr
	if fracStr != "" {
		wholeStr = numStr + fracStr
	} else if numStr == "" {
		wholeStr = "0"
	}

	num, ok := new(big.Int).SetString(wholeStr, 10)
	if !ok {
		panic("parse whole part: " + wholeStr)
	}

	// value = (num / 10^fracDigits)^exp = num^exp / 10^(fracDigits*exp)
	fracDigits := uint(len(fracStr))
	numExp := new(big.Int).Exp(num, new(big.Int).SetUint64(uint64(exp)), nil)
	// denominator = 10^(fracDigits*exp)
	denomPow10 := fracDigits * exp
	denom := new(big.Int).Exp(big.NewInt(10), new(big.Int).SetUint64(uint64(denomPow10)), nil)

	// scale to tfd fractional digits: S = numExp * 10^tfd / 10^(fracDigits*exp)
	scale := new(big.Int).Exp(big.NewInt(10), new(big.Int).SetUint64(uint64(tfd)), nil)
	numScaled := new(big.Int).Mul(numExp, scale)

	q := new(big.Int)
	r := new(big.Int)
	q.QuoRem(numScaled, denom, r)

	// round half-away-from-zero (positive values)
	twiceR := new(big.Int).Mul(r, big.NewInt(2))
	if twiceR.Cmp(denom) >= 0 {
		q.Add(q, big.NewInt(1))
	}

	qs := q.String()
	var intPart, fracPart string
	if len(qs) <= int(tfd) {
		intPart = "0"
		fracPart = qs
		for len(fracPart) < int(tfd) {
			fracPart = "0" + fracPart
		}
	} else {
		intPart = qs[:len(qs)-int(tfd)]
		fracPart = qs[len(qs)-int(tfd):]
	}
	return intPart + "." + fracPart
}

func main() {
	cases := []struct {
		base string
		exp  uint
		tfd  uint
	}{
		{"6.9857", 9, 8},
		{"7.8715", 6, 10},
		{"10.81233", 8, 14},
		{"15.81233", 4, 14},
		{"1.0691372836118", 24, 14},
	}

	for _, c := range cases {
		val := exactPow(c.base, c.exp, c.tfd)
		fmt.Printf("%s^%d  (%d dp) = %s\n", c.base, c.exp, c.tfd, val)
	}
}
