import json
import math
from mpmath import mp, mpf, nstr

# Test cases: (name, base_str, exp, req_frac_digits)
CASES = [
    ("zero_to_zero", "0", 0, 10),
    ("neg_to_zero", "-2.5", 0, 5),
    ("sci_to_zero", "-1.25e-10", 0, 4),
    ("exp_one_neg", "-2.56789", 1, 3),
    ("exp_one_sci", "-1.25e-10", 1, 15),
    ("tier1_pos", "1.23456789", 25, 50),
    ("tier1_neg_even", "-2.5", 4, 6),
    ("tier1_neg_odd", "-2.5", 5, 6),
    ("tier1_500_digits", "1.000123456789", 120, 500),
    ("tier1_1000_digits", "2.718281828459045", 200, 1000),
    ("tier2_exp_over_250", "1.005", 1000, 500),
    ("tier2_neg_exp_1001", "-1.0025", 1001, 500),
    ("tier2_sci_notation", "-1.25e-2", 300, 750),
    ("tier2_5500_digits", "1.4142135623730950488", 300, 5500),
]

def compute_expected(base_str: str, exp: int, req_frac: int) -> str:
    if exp == 0:
        return "1" if req_frac == 0 else "1." + ("0" * req_frac)

    # Estimate integer digits to ensure mp.dps covers integer + fractional + 100 guard digits
    rough_val = abs(float(base_str))
    if rough_val == 0.0:
        int_digits = 1
    else:
        est_log10 = math.log10(rough_val) * exp
        int_digits = max(1, int(math.ceil(est_log10)) + 5)

    mp.dps = int_digits + req_frac + 100

    base = mpf(base_str)
    res = base ** exp

    # Format fixed-point decimal string to req_frac decimal places
    # Use mp.nstr with enough total significant digits to cover integer + fractional parts
    total_sig_figs = int_digits + req_frac + 20
    raw = nstr(res, total_sig_figs, min_fixed=-mp.dps, max_fixed=mp.dps)

    # Round half-to-even cleanly via Python's Decimal on the high-precision string
    from decimal import Decimal, ROUND_HALF_EVEN
    quant = Decimal("1") if req_frac == 0 else Decimal("1e-" + str(req_frac))

    from decimal import (
        Decimal,
        ROUND_HALF_EVEN,
        localcontext
    )

    quant = (
        Decimal("1")
        if req_frac == 0
        else Decimal("1e-" + str(req_frac))
    )

    with localcontext() as ctx:

        ctx.prec = int_digits + req_frac + 200

        rounded = Decimal(raw).quantize(
            quant,
            rounding=ROUND_HALF_EVEN
        )

    return f"{rounded:.{req_frac}f}"

def main():
    vectors = []
    for name, base_str, exp, req_frac in CASES:
        expected = compute_expected(base_str, exp, req_frac)
        vectors.append({
            "name": name,
            "base": base_str,
            "exp": exp,
            "req_frac_digits": req_frac,
            "expected": expected,
        })

    import os
    output_file = os.path.join(os.path.dirname(__file__), "test_vectors.json")
    with open(output_file, "w", encoding="utf-8") as f:
        json.dump(vectors, f, indent=2)

    print(f"Generated {len(vectors)} test vectors in {output_file}")

if __name__ == "__main__":
    main()