package main

import (
	"fmt"
	"math/big"
)

func main() {
	base, _ := new(big.Float).SetPrec(300).SetString("6.9857")
	res, _ := new(big.Float).SetPrec(300).SetString("6.9857")
	exponent := 9
	exponent--
	for i := 0; i < exponent; i++ {
		res.Mul(res, base)
	}
	fmt.Println("6.9857^9 = ", res.Text('f', 8))
}

/*
Original Version
func main() {
	base, _ := new(big.Float).SetPrec(300).SetString("6.9857")
	res := new(big.Float).SetPrec(300)
	for i := 0; i < 9; i++ {
		res.Mul(res, base)
	}
	fmt.Println("6.9857^9 = ", res.Text('f', 12))
}


*/
