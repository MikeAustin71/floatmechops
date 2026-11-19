package floatmech

import "testing"

func TestNumericSignType_String(t *testing.T) {
	tests := []struct {
		name  string
		value NumericSignType
		want  string
	}{
		{
			name:  "NegativeValue",
			value: NegativeValue,
			want:  "Negative",
		},
		{
			name:  "ZeroValue",
			value: ZeroValue,
			want:  "Zero",
		},
		{
			name:  "PositiveValue",
			value: PositiveValue,
			want:  "Positive",
		},
		{
			name:  "invalid_value_less_than_negative",
			value: -2,
			want:  "NOT SET - INVALID",
		},
		{
			name:  "invalid_value_greater_than_positive",
			value: 2,
			want:  "NOT SET - INVALID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nst := tt.value
			got := nst.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNumericSignType_IsNegative(t *testing.T) {
	tests := []struct {
		name  string
		value NumericSignType
		want  bool
	}{
		{
			name:  "NegativeValue",
			value: NegativeValue,
			want:  true,
		},
		{
			name:  "ZeroValue",
			value: ZeroValue,
			want:  false,
		},
		{
			name:  "PositiveValue",
			value: PositiveValue,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nst := tt.value
			if got := nst.IsNegative(); got != tt.want {
				t.Errorf("IsNegative() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumericSignType_IsPositive(t *testing.T) {
	tests := []struct {
		name  string
		value NumericSignType
		want  bool
	}{
		{
			name:  "NegativeValue",
			value: NegativeValue,
			want:  false,
		},
		{
			name:  "ZeroValue",
			value: ZeroValue,
			want:  false,
		},
		{
			name:  "PositiveValue",
			value: PositiveValue,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nst := tt.value
			if got := nst.IsPositive(); got != tt.want {
				t.Errorf("IsPositive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumericSignType_IsZero(t *testing.T) {
	tests := []struct {
		name  string
		value NumericSignType
		want  bool
	}{
		{
			name:  "NegativeValue",
			value: NegativeValue,
			want:  false,
		},
		{
			name:  "ZeroValue",
			value: ZeroValue,
			want:  true,
		},
		{
			name:  "PositiveValue",
			value: PositiveValue,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nst := tt.value
			if got := nst.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumericSignType_IsNegativeOrZero(t *testing.T) {
	tests := []struct {
		name  string
		value NumericSignType
		want  bool
	}{
		{
			name:  "NegativeValue",
			value: NegativeValue,
			want:  true,
		},
		{
			name:  "ZeroValue",
			value: ZeroValue,
			want:  true,
		},
		{
			name:  "PositiveValue",
			value: PositiveValue,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nst := tt.value
			if got := nst.IsNegativeOrZero(); got != tt.want {
				t.Errorf("IsNegativeOrZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumericSignType_IsPositiveOrZero(t *testing.T) {
	tests := []struct {
		name  string
		value NumericSignType
		want  bool
	}{
		{
			name:  "NegativeValue",
			value: NegativeValue,
			want:  false,
		},
		{
			name:  "ZeroValue",
			value: ZeroValue,
			want:  true,
		},
		{
			name:  "PositiveValue",
			value: PositiveValue,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nst := tt.value
			if got := nst.IsPositiveOrZero(); got != tt.want {
				t.Errorf("IsPositiveOrZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumericSignType_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value NumericSignType
		want  bool
	}{
		{
			name:  "NegativeValue",
			value: NegativeValue,
			want:  true,
		},
		{
			name:  "ZeroValue",
			value: ZeroValue,
			want:  true,
		},
		{
			name:  "PositiveValue",
			value: PositiveValue,
			want:  true,
		},
		{
			name:  "invalid_value_less_than_negative",
			value: -2,
			want:  false,
		},
		{
			name:  "invalid_value_greater_than_positive",
			value: 2,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nst := tt.value
			if got := nst.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumericSignType_SetPositive(t *testing.T) {
	nst := NegativeValue
	nst.SetPositive()
	if nst != PositiveValue {
		t.Errorf("after SetPositive() = %v, want %v", nst, PositiveValue)
	}
}

func TestNumericSignType_SetZero(t *testing.T) {
	nst := PositiveValue
	nst.SetZero()
	if nst != ZeroValue {
		t.Errorf("after SetZero() = %v, want %v", nst, ZeroValue)
	}
}

func TestNumericSignType_SetNegative(t *testing.T) {
	nst := ZeroValue
	nst.SetNegative()
	if nst != NegativeValue {
		t.Errorf("after SetNegative() = %v, want %v", nst, NegativeValue)
	}
}

func TestNumericSignType_Values(t *testing.T) {
	if NegativeValue != -1 {
		t.Errorf("NegativeValue = %d, want -1", int(NegativeValue))
	}
	if ZeroValue != 0 {
		t.Errorf("ZeroValue = %d, want 0", int(ZeroValue))
	}
	if PositiveValue != 1 {
		t.Errorf("PositiveValue = %d, want 1", int(PositiveValue))
	}
}

func TestNumericSignTypeNames(t *testing.T) {
	want := [...]string{"Negative", "Zero", "Positive"}
	if len(NumericSignTypeNames) != len(want) {
		t.Fatalf("len(NumericSignTypeNames) = %d, want %d", len(NumericSignTypeNames), len(want))
	}
	for i := range want {
		if NumericSignTypeNames[i] != want[i] {
			t.Errorf("NumericSignTypeNames[%d] = %q, want %q", i, NumericSignTypeNames[i], want[i])
		}
	}
}
