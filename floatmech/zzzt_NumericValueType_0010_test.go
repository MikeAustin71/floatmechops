package floatmech

import "testing"

func TestNumericValueType_String(t *testing.T) {
	tests := []struct {
		name  string
		value NumericValueType
		want  string
	}{
		{
			name:  "IntegerValue",
			value: IntegerValue,
			want:  "IntegerValue",
		},
		{
			name:  "FloatingPointValue",
			value: FloatingPointValue,
			want:  "FloatingPointValue",
		},
		{
			name:  "invalid_zero_value",
			value: 0,
			want:  "INVALID",
		},
		{
			name:  "invalid_value_greater_than_floating_point",
			value: 3,
			want:  "INVALID",
		},
		{
			name:  "invalid_negative_value",
			value: -1,
			want:  "INVALID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nvt := tt.value
			got := nvt.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNumericValueType_IsInteger(t *testing.T) {
	tests := []struct {
		name  string
		value NumericValueType
		want  bool
	}{
		{
			name:  "IntegerValue",
			value: IntegerValue,
			want:  true,
		},
		{
			name:  "FloatingPointValue",
			value: FloatingPointValue,
			want:  false,
		},
		{
			name:  "invalid_zero_value",
			value: 0,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nvt := tt.value
			if got := nvt.IsInteger(); got != tt.want {
				t.Errorf("IsInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumericValueType_IsFloatingPoint(t *testing.T) {
	tests := []struct {
		name  string
		value NumericValueType
		want  bool
	}{
		{
			name:  "IntegerValue",
			value: IntegerValue,
			want:  false,
		},
		{
			name:  "FloatingPointValue",
			value: FloatingPointValue,
			want:  true,
		},
		{
			name:  "invalid_zero_value",
			value: 0,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nvt := tt.value
			if got := nvt.IsFloatingPoint(); got != tt.want {
				t.Errorf("IsFloatingPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumericValueType_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value NumericValueType
		want  bool
	}{
		{
			name:  "IntegerValue",
			value: IntegerValue,
			want:  true,
		},
		{
			name:  "FloatingPointValue",
			value: FloatingPointValue,
			want:  true,
		},
		{
			name:  "invalid_zero_value",
			value: 0,
			want:  false,
		},
		{
			name:  "invalid_value_greater_than_floating_point",
			value: 3,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nvt := tt.value
			if got := nvt.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumericValueType_SetInteger(t *testing.T) {
	nvt := FloatingPointValue
	nvt.SetInteger()
	if nvt != IntegerValue {
		t.Errorf("after SetInteger() = %v, want %v", nvt, IntegerValue)
	}
}

func TestNumericValueType_SetFloatingPoint(t *testing.T) {
	nvt := IntegerValue
	nvt.SetFloatingPoint()
	if nvt != FloatingPointValue {
		t.Errorf("after SetFloatingPoint() = %v, want %v", nvt, FloatingPointValue)
	}
}

func TestNumericValueType_Values(t *testing.T) {
	if IntegerValue != 1 {
		t.Errorf("IntegerValue = %d, want 1", int(IntegerValue))
	}
	if FloatingPointValue != 2 {
		t.Errorf("FloatingPointValue = %d, want 2", int(FloatingPointValue))
	}
}

func TestNumericValueTypeNames(t *testing.T) {
	want := [...]string{"INVALID", "IntegerValue", "FloatingPointValue"}
	if len(NumericValueTypeNames) != len(want) {
		t.Fatalf("len(NumericValueTypeNames) = %d, want %d", len(NumericValueTypeNames), len(want))
	}
	for i := range want {
		if NumericValueTypeNames[i] != want[i] {
			t.Errorf("NumericValueTypeNames[%d] = %q, want %q", i, NumericValueTypeNames[i], want[i])
		}
	}
}
