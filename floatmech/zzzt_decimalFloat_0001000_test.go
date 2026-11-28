package floatmech

import (
	"encoding/json"
	"os"
	"testing"
)

type powDecimalVector struct {
	Name          string `json:"name"`
	Base          string `json:"base"`
	Exp           uint   `json:"exp"`
	ReqFracDigits uint   `json:"req_frac_digits"`
	Expected      string `json:"expected"`
}

func TestPowDecimal_00010(t *testing.T) {
	data, err := os.ReadFile("test_vectors.json")
	if err != nil {
		t.Fatalf("failed to read test vectors: %v", err)
	}

	var vectors []powDecimalVector
	if err := json.Unmarshal(data, &vectors); err != nil {
		t.Fatalf("failed to unmarshal test vectors: %v", err)
	}

	df := &DecimalFloat{}

	for _, v := range vectors {
		t.Run(v.Name, func(t *testing.T) {
			resStr, _, err := df.PowDecimal(v.Base, v.Exp, v.ReqFracDigits)
			if err != nil {
				t.Errorf("PowDecimal(%s, %d, %d) returned error: %v", v.Base, v.Exp, v.ReqFracDigits, err)
				return
			}

			if resStr != v.Expected {
				t.Errorf("PowDecimal(%s, %d, %d) = %s; want %s", v.Base, v.Exp, v.ReqFracDigits, resStr, v.Expected)
				return
			}
		})
	}
	return
}
