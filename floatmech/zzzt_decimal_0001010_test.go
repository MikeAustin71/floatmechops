package floatmech

import (
	"math/big"
	"testing"
)

// TestDecimal_CopyIn_NilReceiver
//
// Verifies that CopyIn returns an error when the host, or
// receiving Decimal object (pointer receiver 'd') is a nil pointer.
func TestDecimal_CopyIn_NilReceiver(t *testing.T) {

	copySource := new(Decimal).NewFromFloat64(6.9857, 8, big.ToNearestAway)

	receiver := (*Decimal)(nil)

	err := receiver.CopyIn(copySource)

	if err == nil {
		t.Errorf("\n" +
			"Expected CopyIn to return an error for a nil pointer receiver\n" +
			"  Actual CopyIn returned a nil error")
	}

	return
}

// TestDecimal_CopyIn_NilCopySource
//
// Verifies that CopyIn returns an error when the incoming
// copySource *Decimal input parameter is a nil pointer.
func TestDecimal_CopyIn_NilCopySource(t *testing.T) {

	receiver := new(Decimal).NewFromFloat64(0, 8, big.ToNearestAway)

	copySource := (*Decimal)(nil)

	err := receiver.CopyIn(copySource)

	if err == nil {
		t.Errorf("\n" +
			"Expected CopyIn to return an error for a nil copySource\n" +
			"  Actual CopyIn returned a nil error")
	}

	return
}

// TestDecimal_CopyIn_ValidCopy
//
// Verifies that CopyIn copies the member variable values of the
// incoming copySource Decimal object into the member variables of
// the host, or receiving Decimal object.
func TestDecimal_CopyIn_ValidCopy(t *testing.T) {

	copySource := new(Decimal).NewFromFloat64(6.9857, 8, big.ToNearestAway)

	receiver := new(Decimal).NewFromFloat64(1.0, 2, big.ToNearestEven)

	err := receiver.CopyIn(copySource)

	if err != nil {
		t.Errorf("\n"+
			"Expected CopyIn to succeed for a valid copySource\n"+
			"  Actual CopyIn returned: %s", err.Error())
		return
	}

	if receiver.val != copySource.val {
		t.Errorf("\n"+
			"Expected receiver.val to equal copySource.val\n"+
			"  Actual receiver.val: %s\n"+
			"  Actual copySource.val: %s",
			receiver.val.String(), copySource.val.String())
		return
	}

	if receiver.targetFractionalDigits != copySource.targetFractionalDigits {
		t.Errorf("\n"+
			"Expected receiver.targetFractionalDigits = %d\n"+
			"  Actual receiver.targetFractionalDigits = %d",
			copySource.targetFractionalDigits, receiver.targetFractionalDigits)
		return
	}

	if receiver.prec != copySource.prec {
		t.Errorf("\n"+
			"Expected receiver.prec = %d\n"+
			"  Actual receiver.prec = %d",
			copySource.prec, receiver.prec)
		return
	}

	if receiver.mode != copySource.mode {
		t.Errorf("\n"+
			"Expected receiver.mode = %s\n"+
			"  Actual receiver.mode = %s",
			copySource.mode, receiver.mode)
		return
	}

	return
}

// TestDecimal_CopyOut_NilReceiver
//
// Verifies that CopyOut returns an error when the host, or
// source Decimal object (pointer receiver 'd') is a nil pointer.
func TestDecimal_CopyOut_NilReceiver(t *testing.T) {

	copyDestination := new(Decimal).NewFromFloat64(1.0, 2, big.ToNearestEven)

	receiver := (*Decimal)(nil)

	err := receiver.CopyOut(copyDestination)

	if err == nil {
		t.Errorf("\n" +
			"Expected CopyOut to return an error for a nil pointer receiver\n" +
			"  Actual CopyOut returned a nil error")
	}

	return
}

// TestDecimal_CopyOut_NilCopyDestination
//
// Verifies that CopyOut returns an error when the outgoing
// copyDestination *Decimal input parameter is a nil pointer.
func TestDecimal_CopyOut_NilCopyDestination(t *testing.T) {

	receiver := new(Decimal).NewFromFloat64(6.9857, 8, big.ToNearestAway)

	copyDestination := (*Decimal)(nil)

	err := receiver.CopyOut(copyDestination)

	if err == nil {
		t.Errorf("\n" +
			"Expected CopyOut to return an error for a nil copyDestination\n" +
			"  Actual CopyOut returned a nil error")
	}

	return
}

// TestDecimal_CopyOut_ValidCopy
//
// Verifies that CopyOut copies the member variable values of the
// host (source) Decimal object into the member variables of the
// outgoing copyDestination Decimal object.
func TestDecimal_CopyOut_ValidCopy(t *testing.T) {

	receiver := new(Decimal).NewFromFloat64(6.9857, 8, big.ToNearestAway)

	copyDestination := new(Decimal).NewFromFloat64(1.0, 2, big.ToNearestEven)

	err := receiver.CopyOut(copyDestination)

	if err != nil {
		t.Errorf("\n"+
			"Expected CopyOut to succeed for a valid copyDestination\n"+
			"  Actual CopyOut returned: %s", err.Error())
	}

	if copyDestination.val != receiver.val {
		t.Errorf("\n"+
			"Expected copyDestination.val to equal receiver.val\n"+
			"  Actual copyDestination.val: %s\n"+
			"  Actual receiver.val: %s",
			copyDestination.val.String(), receiver.val.String())
	}

	if copyDestination.targetFractionalDigits != receiver.targetFractionalDigits {
		t.Errorf("\n"+
			"Expected copyDestination.targetFractionalDigits = %d\n"+
			"  Actual copyDestination.targetFractionalDigits = %d",
			receiver.targetFractionalDigits, copyDestination.targetFractionalDigits)
	}

	if copyDestination.prec != receiver.prec {
		t.Errorf("\n"+
			"Expected copyDestination.prec = %d\n"+
			"  Actual copyDestination.prec = %d",
			receiver.prec, copyDestination.prec)
	}

	if copyDestination.mode != receiver.mode {
		t.Errorf("\n"+
			"Expected copyDestination.mode = %s\n"+
			"  Actual copyDestination.mode = %s",
			receiver.mode, copyDestination.mode)
	}

	return
}
