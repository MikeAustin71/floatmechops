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

	copySource, err := new(Decimal).NewFromFloatNumStr("6.9857", 8, big.ToNearestAway)

	if err != nil {
		t.Errorf("NewFromFloatNumStr() failed: %s", err.Error())

		return
	}

	receiver := (*Decimal)(nil)

	err = receiver.CopyIn(copySource)

	if err == nil {
		t.Errorf("\n" +
			"Expected CopyIn to return an error for a nil pointer receiver\n" +
			"  Actual CopyIn returned a nil error")

		return
	}

	return
}

// TestDecimal_CopyIn_NilCopySource
//
// Verifies that CopyIn returns an error when the incoming
// copySource *Decimal input parameter is a nil pointer.
func TestDecimal_CopyIn_NilCopySource(t *testing.T) {

	receiver, err := new(Decimal).NewFromFloatNumStr("0", 8, big.ToNearestAway)

	if err != nil {
		t.Errorf("NewFromFloatNumStr() failed: %s", err.Error())

		return
	}

	copySource := (*Decimal)(nil)

	err = receiver.CopyIn(copySource)

	if err == nil {
		t.Errorf("\n" +
			"Expected CopyIn to return an error for a nil copySource\n" +
			"  Actual CopyIn returned a nil error")

		return
	}

	return
}

// TestDecimal_CopyIn_ValidCopy
//
// Verifies that CopyIn copies the member variable values of the
// incoming copySource Decimal object into the member variables of
// the host, or receiving Decimal object.
func TestDecimal_CopyIn_ValidCopy(t *testing.T) {

	targetFractionalDigits := uint(8)

	copySource, err := new(Decimal).NewFromFloatNumStr("6.9857", targetFractionalDigits, big.ToNearestAway)

	if err != nil {
		t.Errorf("NewFromFloatNumStr() failed: %s", err.Error())

		return
	}

	receiver, err := new(Decimal).NewFromFloatNumStr("1.0", 2, big.ToNearestEven)

	if err != nil {
		t.Errorf("NewFromFloatNumStr() failed: %s", err.Error())

		return
	}

	err = receiver.CopyIn(copySource)

	if err != nil {
		t.Errorf("\n"+
			"Expected CopyIn to succeed for a valid copySource\n"+
			"  Actual CopyIn returned: %s", err.Error())

		return
	}

	if receiver.mode != copySource.mode {
		t.Errorf("\n"+
			"Expected receiver.mode = %s\n"+
			"  Actual receiver.mode = %s",
			copySource.mode, receiver.mode)

		return
	}

	copySourceValStr := copySource.val.Text('f', int(targetFractionalDigits))

	receiverValStr := receiver.val.Text('f', int(targetFractionalDigits))

	if receiverValStr != copySourceValStr {
		t.Errorf("\n"+
			"Expected receiverValStr to equal copySourceValStr\n"+
			"  Actual receiverValStr: %s\n"+
			"  Actual copySourceValStr: %s",
			receiverValStr, copySourceValStr)

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

	return
}

// TestDecimal_CopyOut_NilReceiver
//
// Verifies that CopyOut returns an error when the host, or
// source Decimal object (pointer receiver 'd') is a nil pointer.
func TestDecimal_CopyOut_NilReceiver(t *testing.T) {

	copyDestination, err := new(Decimal).NewFromFloatNumStr("1.0", 2, big.ToNearestEven)

	if err != nil {
		t.Errorf("NewFromFloatNumStr() failed: %s", err.Error())

		return
	}

	receiver := (*Decimal)(nil)

	err = receiver.CopyOut(copyDestination)

	if err == nil {
		t.Errorf("\n" +
			"Expected CopyOut to return an error for a nil pointer receiver\n" +
			"  Actual CopyOut returned a nil error")

		return
	}

	return
}

// TestDecimal_CopyOut_NilCopyDestination
//
// Verifies that CopyOut returns an error when the outgoing
// copyDestination *Decimal input parameter is a nil pointer.
func TestDecimal_CopyOut_NilCopyDestination(t *testing.T) {

	receiver, err := new(Decimal).NewFromFloatNumStr("6.9857", 8, big.ToNearestAway)

	if err != nil {
		t.Errorf("NewFromFloatNumStr() failed: %s", err.Error())

		return
	}

	copyDestination := (*Decimal)(nil)

	err = receiver.CopyOut(copyDestination)

	if err == nil {
		t.Errorf("\n" +
			"Expected CopyOut to return an error for a nil copyDestination\n" +
			"  Actual CopyOut returned a nil error")

		return
	}

	return
}

// TestDecimal_CopyOut_ValidCopy
//
// Verifies that CopyOut copies the member variable values of the
// host (source) Decimal object into the member variables of the
// outgoing copyDestination Decimal object.
func TestDecimal_CopyOut_ValidCopy(t *testing.T) {

	targetFractionalDigits := uint(8)

	receiver, err := new(Decimal).NewFromFloatNumStr("6.9857", targetFractionalDigits, big.ToNearestAway)

	if err != nil {
		t.Errorf("NewFromFloatNumStr() failed: %s", err.Error())

		return
	}

	copyDestination, err := new(Decimal).NewFromFloatNumStr("1.0", 2, big.ToNearestEven)

	if err != nil {
		t.Errorf("NewFromFloatNumStr() failed: %s", err.Error())

		return
	}

	err = receiver.CopyOut(copyDestination)

	if err != nil {
		t.Errorf("\n"+
			"Expected CopyOut to succeed for a valid copyDestination\n"+
			"  Actual CopyOut returned: %s", err.Error())

		return
	}

	if copyDestination.mode != receiver.mode {
		t.Errorf("\n"+
			"Expected copyDestination.mode = %s\n"+
			"  Actual copyDestination.mode = %s",
			receiver.mode, copyDestination.mode)

		return
	}

	copyDestinationValStr := copyDestination.val.Text('f', int(targetFractionalDigits))

	receiverValStr := receiver.val.Text('f', int(targetFractionalDigits))

	if copyDestinationValStr != receiverValStr {
		t.Errorf("\n"+
			"Expected copyDestinationValStr to equal receiverValStr\n"+
			"  Actual copyDestinationValStr: %s\n"+
			"  Actual receiverValStr: %s",
			copyDestinationValStr, receiverValStr)

		return
	}

	if copyDestination.targetFractionalDigits != receiver.targetFractionalDigits {
		t.Errorf("\n"+
			"Expected copyDestination.targetFractionalDigits = %d\n"+
			"  Actual copyDestination.targetFractionalDigits = %d",
			receiver.targetFractionalDigits, copyDestination.targetFractionalDigits)

		return
	}

	if copyDestination.prec != receiver.prec {
		t.Errorf("\n"+
			"Expected copyDestination.prec = %d\n"+
			"  Actual copyDestination.prec = %d",
			receiver.prec, copyDestination.prec)

		return
	}

	return
}

// TestDecimal_Clone_NilReceiver
//
// Verifies that Clone returns an error when the host Decimal object
// (pointer receiver 'd') is a nil pointer.
func TestDecimal_Clone_NilReceiver(t *testing.T) {

	receiver := (*Decimal)(nil)

	cloned, err := receiver.Clone()

	if err == nil {
		t.Errorf("\n" +
			"Expected Clone to return an error for a nil pointer receiver\n" +
			"  Actual Clone returned a nil error")

		return
	}

	if cloned == nil {
		t.Errorf("\n" +
			"Expected Clone to return a non-nil Decimal object even on error\n" +
			"  Actual Clone returned a nil Decimal object")

		return
	}

	return
}

// TestDecimal_Clone_NilVal
//
// Verifies that Clone returns an error when the host Decimal object's
// internal 'val' member variable is a nil pointer.
func TestDecimal_Clone_NilVal(t *testing.T) {

	receiver := &Decimal{val: nil}

	cloned, err := receiver.Clone()

	if err == nil {
		t.Errorf("\n" +
			"Expected Clone to return an error for a nil internal 'val'\n" +
			"  Actual Clone returned a nil error")

		return
	}

	if cloned == nil {
		t.Errorf("\n" +
			"Expected Clone to return a non-nil Decimal object even on error\n" +
			"  Actual Clone returned a nil Decimal object")

		return
	}

	return
}

// TestDecimal_Clone_ValidClone
//
// Verifies that Clone returns a newly created Decimal object which is
// an exact clone of the host Decimal object.
func TestDecimal_Clone_ValidClone(t *testing.T) {

	targetFractionalDigits := uint(8)
	valStr := "6.9857"

	receiver, err := new(Decimal).NewFromFloatNumStr(valStr, targetFractionalDigits, big.ToNearestAway)

	if err != nil {
		t.Errorf("NewFromFloatNumStr() failed: %s", err.Error())

		return
	}

	cloned, err := receiver.Clone()

	if err != nil {
		t.Errorf("\n"+
			"Expected Clone to succeed for a valid Decimal object\n"+
			"  Actual Clone returned: %s", err.Error())

		return
	}

	if cloned == nil {
		t.Errorf("\n" +
			"Expected Clone to return a non-nil Decimal object\n" +
			"  Actual Clone returned a nil Decimal object")

		return
	}

	if cloned == receiver {
		t.Errorf("\n" +
			"Expected Clone to return a new object (different pointer)\n" +
			"  Actual Clone returned the same pointer as the receiver")

		return
	}

	if cloned.mode != receiver.mode {
		t.Errorf("\n"+
			"Expected cloned.mode = %s\n"+
			"  Actual cloned.mode = %s",
			receiver.mode, cloned.mode)

		return
	}

	receiverValStr := receiver.val.Text('f', int(targetFractionalDigits))
	clonedValStr := cloned.val.Text('f', int(targetFractionalDigits))

	if clonedValStr != receiverValStr {
		t.Errorf("\n"+
			"Expected clonedValStr to equal receiverValStr\n"+
			"  Actual clonedValStr: %s\n"+
			"  Actual receiverValStr: %s",
			clonedValStr, receiverValStr)

		return
	}

	if cloned.targetFractionalDigits != receiver.targetFractionalDigits {
		t.Errorf("\n"+
			"Expected cloned.targetFractionalDigits = %d\n"+
			"  Actual cloned.targetFractionalDigits = %d",
			receiver.targetFractionalDigits, cloned.targetFractionalDigits)

		return
	}

	if cloned.prec != receiver.prec {
		t.Errorf("\n"+
			"Expected cloned.prec = %d\n"+
			"  Actual cloned.prec = %d",
			receiver.prec, cloned.prec)

		return
	}

	return
}
