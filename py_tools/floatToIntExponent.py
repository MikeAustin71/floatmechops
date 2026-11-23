from mpmath import mp

# ============================================================
# Configuration
# ============================================================

# Base value represented as a string
base_str = "10.81233"

# Positive integer exponent represented as a string
exponent_str = "8"

# Number of digits after the decimal point desired in output
fractional_digits = 100

# Additional precision used during calculations
guard_digits = 100

# ============================================================
# Setup Precision
# ============================================================

# Calculate with extra precision to reduce rounding error
mp.dps = fractional_digits + guard_digits

# Convert strings to mpmath values
base = mp.mpf(base_str)
exponent = int(exponent_str)

# ============================================================
# Calculate
# ============================================================

result = base ** exponent

# ============================================================
# Display Result
# ============================================================

print(f"Base     : {base_str}")
print(f"Exponent : {exponent}")
print()
print(f"Result ({fractional_digits} fractional digits):")
print(f"{result:.{fractional_digits}f}")