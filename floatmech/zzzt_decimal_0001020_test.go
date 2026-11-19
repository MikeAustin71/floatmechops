package floatmech

import (
	"math/big"
	"testing"
)

// TestDecimal_CopyIn_NilCopySourceVal
//
// Verifies that CopyIn returns an error when the incoming
// copySource *Decimal object has a nil pointer 'val' member.
func TestDecimal_CopyIn_NilCopySourceVal(t *testing.T) {

	receiver := new(Decimal).NewFromFloat64(1.0, 2, big.ToNearestEven)

	copySource := &Decimal{}

	err := receiver.CopyIn(copySource)

	if err == nil {
		t.Errorf("\n" +
			"Expected CopyIn to return an error for a nil copySource.val\n" +
			"  Actual CopyIn returned a nil error")
	}

	return
}

// TestDecimal_CopyOut_NilHostVal
//
// Verifies that CopyOut returns an error when the host Decimal
// object has a nil pointer 'val' member.
func TestDecimal_CopyOut_NilHostVal(t *testing.T) {

	host := &Decimal{}

	copyDestination := new(Decimal).NewFromFloat64(1.0, 2, big.ToNearestEven)

	err := host.CopyOut(copyDestination)

	if err == nil {
		t.Errorf("\n" +
			"Expected CopyOut to return an error for a nil host val\n" +
			"  Actual CopyOut returned a nil error")
	}

	return
}

// TestDecimal_Clone_NilHostVal
//
// Verifies that Clone returns an error when the host Decimal
// object has a nil pointer 'val' member.
func TestDecimal_Clone_NilHostVal(t *testing.T) {

	host := &Decimal{}

	_, err := host.Clone()

	if err == nil {
		t.Errorf("\n" +
			"Expected Clone to return an error for a nil host val\n" +
			"  Actual Clone returned a nil error")
	}

	return
}
