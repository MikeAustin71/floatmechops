package floatmech

import (
	"math/big"
	"testing"
)

func TestRoundingType_String(t *testing.T) {
	tests := []struct {
		name  string
		value RoundingType
		want  string
	}{
		{
			name:  "ToNearestEvenVal",
			value: ToNearestEvenVal,
			want:  "ToNearestEven",
		},
		{
			name:  "ToNearestAwayVal",
			value: ToNearestAwayVal,
			want:  "ToNearestAway",
		},
		{
			name:  "ToZeroVal",
			value: ToZeroVal,
			want:  "ToZero",
		},
		{
			name:  "AwayFromZeroVal",
			value: AwayFromZeroVal,
			want:  "AwayFromZero",
		},
		{
			name:  "ToNegativeInfVal",
			value: ToNegativeInfVal,
			want:  "ToNegativeInf",
		},
		{
			name:  "ToPositiveInfVal",
			value: ToPositiveInfVal,
			want:  "ToPositiveInf",
		},
		{
			name:  "invalid_value_above_max",
			value: 6,
			want:  "NOT SET - INVALID",
		},
		{
			name:  "invalid_negative_value",
			value: -1,
			want:  "NOT SET - INVALID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := tt.value
			got := rt.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRoundingType_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value RoundingType
		want  bool
	}{
		{
			name:  "ToNearestEvenVal",
			value: ToNearestEvenVal,
			want:  true,
		},
		{
			name:  "ToNearestAwayVal",
			value: ToNearestAwayVal,
			want:  true,
		},
		{
			name:  "ToZeroVal",
			value: ToZeroVal,
			want:  true,
		},
		{
			name:  "AwayFromZeroVal",
			value: AwayFromZeroVal,
			want:  true,
		},
		{
			name:  "ToNegativeInfVal",
			value: ToNegativeInfVal,
			want:  true,
		},
		{
			name:  "ToPositiveInfVal",
			value: ToPositiveInfVal,
			want:  true,
		},
		{
			name:  "invalid_value_above_max",
			value: 6,
			want:  false,
		},
		{
			name:  "invalid_negative_value",
			value: -1,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := tt.value
			if got := rt.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundingType_GetBigMathRoundingMode(t *testing.T) {
	tests := []struct {
		name      string
		value     RoundingType
		want      big.RoundingMode
		wantError bool
	}{
		{
			name:  "ToNearestEvenVal",
			value: ToNearestEvenVal,
			want:  big.ToNearestEven,
		},
		{
			name:  "ToNearestAwayVal",
			value: ToNearestAwayVal,
			want:  big.ToNearestAway,
		},
		{
			name:  "ToZeroVal",
			value: ToZeroVal,
			want:  big.ToZero,
		},
		{
			name:  "AwayFromZeroVal",
			value: AwayFromZeroVal,
			want:  big.AwayFromZero,
		},
		{
			name:  "ToNegativeInfVal",
			value: ToNegativeInfVal,
			want:  big.ToNegativeInf,
		},
		{
			name:  "ToPositiveInfVal",
			value: ToPositiveInfVal,
			want:  big.ToPositiveInf,
		},
		{
			name:      "invalid_value_above_max",
			value:     6,
			wantError: true,
		},
		{
			name:      "invalid_negative_value",
			value:     -1,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := tt.value
			got, err := rt.GetBigMathRoundingMode()

			if tt.wantError {
				if err == nil {
					t.Errorf("GetBigMathRoundingMode() expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("GetBigMathRoundingMode() unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("GetBigMathRoundingMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundingType_SetToNearestEven(t *testing.T) {
	rt := ToPositiveInfVal
	rt.SetToNearestEven()
	if rt != ToNearestEvenVal {
		t.Errorf("after SetToNearestEven() = %v, want %v", rt, ToNearestEvenVal)
	}
}

func TestRoundingType_SetToNearestAway(t *testing.T) {
	rt := ToZeroVal
	rt.SetToNearestAway()
	if rt != ToNearestAwayVal {
		t.Errorf("after SetToNearestAway() = %v, want %v", rt, ToNearestAwayVal)
	}
}

func TestRoundingType_SetToZero(t *testing.T) {
	rt := ToNearestEvenVal
	rt.SetToZero()
	if rt != ToZeroVal {
		t.Errorf("after SetToZero() = %v, want %v", rt, ToZeroVal)
	}
}

func TestRoundingType_SetAwayFromZero(t *testing.T) {
	rt := ToNegativeInfVal
	rt.SetAwayFromZero()
	if rt != AwayFromZeroVal {
		t.Errorf("after SetAwayFromZero() = %v, want %v", rt, AwayFromZeroVal)
	}
}

func TestRoundingType_SetToNegativeInf(t *testing.T) {
	rt := ToPositiveInfVal
	rt.SetToNegativeInf()
	if rt != ToNegativeInfVal {
		t.Errorf("after SetToNegativeInf() = %v, want %v", rt, ToNegativeInfVal)
	}
}

func TestRoundingType_SetToPositiveInf(t *testing.T) {
	rt := AwayFromZeroVal
	rt.SetToPositiveInf()
	if rt != ToPositiveInfVal {
		t.Errorf("after SetToPositiveInf() = %v, want %v", rt, ToPositiveInfVal)
	}
}

func TestRoundingType_Values(t *testing.T) {
	if ToNearestEvenVal != 0 {
		t.Errorf("ToNearestEvenVal = %d, want 0", int(ToNearestEvenVal))
	}

	if ToNearestAwayVal != 1 {
		t.Errorf("ToNearestAwayVal = %d, want 1", int(ToNearestAwayVal))
	}

	if ToZeroVal != 2 {
		t.Errorf("ToZeroVal = %d, want 2", int(ToZeroVal))
	}

	if AwayFromZeroVal != 3 {
		t.Errorf("AwayFromZeroVal = %d, want 3", int(AwayFromZeroVal))
	}

	if ToNegativeInfVal != 4 {
		t.Errorf("ToNegativeInfVal = %d, want 4", int(ToNegativeInfVal))
	}

	if ToPositiveInfVal != 5 {
		t.Errorf("ToPositiveInfVal = %d, want 5", int(ToPositiveInfVal))
	}
}

func TestRoundingType_AllValidValuesRoundTrip(t *testing.T) {
	values := []RoundingType{
		ToNearestEvenVal,
		ToNearestAwayVal,
		ToZeroVal,
		AwayFromZeroVal,
		ToNegativeInfVal,
		ToPositiveInfVal,
	}

	for _, v := range values {
		if !v.IsValid() {
			t.Errorf("IsValid() = false for valid value %v", v)
		}

		if v.String() == "NOT SET - INVALID" {
			t.Errorf("String() = %q for valid value %v", v.String(), v)
		}

		if _, err := v.GetBigMathRoundingMode(); err != nil {
			t.Errorf("GetBigMathRoundingMode() unexpected error for %v: %v", v, err)
		}
	}
}
