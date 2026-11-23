[TOC]

# Calculated Mantissa for ***go*** Big Float Positive Integer Exponentiation Version 0010 2026-11-22



## Background

This analysis is intended to identify the worst case number of digits required for a ***go*** Big Float Mantissa to accurately represent the total number of numerical digits resulting from a number (integer or floating point) being raised to a positive integer exponent.

## Calculations

Actual results for this analysis were computed with the Python Version 3.14 [***mpmath*** Library](https://mpmath.org/doc/current/#advanced-mathematics).



## Collected Raw Data

[See PDF File of Collected Raw Data](D:\Programming\go\big_Math_Package\Big_Float\Big_Float_Mantissa\Big_Float_Mantissa_Emperical_Data.pdf)



## Actual Results



### $999.999 ^ 5$ (Fractional digits of accuracy = 100)

#### Actual Results

$X^0$ =                                  1.000
$X^1$ =                              999.999
$X^2$  =                      999 998.000 001 000 000 000 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000

$X^3$  =              999 997 000.002 999 999 000 000 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000

$X^4$ =        999 996 000 005.999 996 000 001 000 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000

$X^5$ = 999 995 000 009 999.990 000 004 999 999 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000

#### Digits in The Result: Integer + Fractional = Total Digits

$X^0$ =   3 +   3 =   1
$X^1$ =   3 +   3 =   6
$X^2$ =   6 +   6 = 12
$X^3$ =   9 +   9 = 18
$X^4$ = 12 + 12 = 24
$X^5$ = 15 + 15 = 30



--------------------------------------------------------------------------------
### $777.777^5$  (Fractional digits of accuracy = 100)

#### Actual Results

$X^0$  =                                 1.000
$X^1$  =                             777.777
$X^2$  =                     604 937.061 729 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
$X^3$  =               470 506 133.060 396 433 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
$X^4$  =        365 948 848 653.315 956 469 441 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
$X^5$  =  284 626 597 659 030.124 674 932 412 657 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000

#### Digits in The Result: Integer + Fractional = Total Digits

$X^0$ =   3 +   3 =   1
$X^1$ =   3 +   3 =   6
$X^2$ =   6 +   6 = 12
$X^3$ =   9 +   9 = 18
$X^4$ = 12 + 12 = 24
$X^5$ = 15 + 15 = 30

--------------------------------------------------------------------------------

### $444.444^5$  (Fractional digits of accuracy = 100)

#### Actual Results

$X^0$                                 1.000
$X^1$                            444.444
$X^2$                    197 530.469 136 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
$X^3$                87 791 231.824 680 384 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
$X^4$         39 018 286 237.088 248 586 496 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
$X^5$   17 341 443 208 356.449 554 776 628 224 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000

#### Digits in The Result: Integer + Fractional = Total Digits

$X^0$ =   1 +   0 =   1
$X^1$ =   3 +   3 =   6
$X^2$ =   6 +   6 = 12
$X^3$ =   8 +   9 = 17
$X^4$ = 11 + 12 = 23
$X^5$ = 14 + 15 = 29

--------------------------------------------------------------------------------

### $333.333^5$ (Fractional digits of accuracy = 100)

#### Actual Results

$X^0$                              1.000
$X^1$                          333.333
$X^2$                  111 110.888 889 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
$X^3$             37 036 925.926 037 037 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
$X^4$       12 345 629 629.703 703 654 321 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
$X^5$      411 505 761 358.024 650 205 781 893 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000

#### Numerical Digits in The Result: Integer Digits + Fractional Digits = Total Digits

$X^0$ =   1 +   0 =   1
$X^1$ =   3 +   3 =   6
$X^2$ =   6 +   6 = 12
$X^3$ =   8 +   9 = 17
$X^4$ = 11 + 12 = 23
$X^5$ = 12 + 15 = 27

--------------------------------------------------------------------------------

### $111.111^5$ (Fractional digits of accuracy = 100)

#### Actual Results

$X^0$                               1.000
$X^1$                           111.111
$X^2$                      12 345.654 321 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
$X^3$                 1 371 737.997 260 631 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
$X^4$             152 415 180.613 625 971 041 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
$X^5$        16 935 003 133.160 595 268 336 551 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000



#### Numerical Digits in The Result: Integer Digits + Fractional Fractional = Total Digits

$X^0$ =   1 +   0 =  1
$X^1$ =   3 +   3 =  6
$X^2$ =   5 +   6 = 11
$X^3$ =   7 +   9 = 16
$X^4$ =   9 + 12 = 21
$X^5$ = 11 + 15 = 26




--------------------------------------------------------------------------------


**Prediction Using 'Worst Case Formula'**
--------------------------------------------------------------------------------

### Exponentiation Parameters

baseStr := "10.81233"

baseStr Actual Fractional Digits = "812 33" = 5

exponentValue = 8

requiredFractionalDigits = 14

baseStr Integer and Fractional = Total Digits

baseStr = 2 + 5 = 7 Total Digits

### Expected Total Digits

Total Base Digits = 7
Number of Multiplications = 8
Total Additional Digits Through Multiplication = 48
    Number of Multiplications x 6 = 8 x 6 = 48

Delta Required Fractional Digits Minus 'baseStr' Actual Fractional Digits
14 - 5 = 9

### Predicted Worst Case Number of Digits Required by ***go*** Big Float Mantissa to Accurately Represent Calculated Result 

**Total Worst Case Number of Digits for Mantissa**  =
  **Total Additional Digits Through Multiplication** + **Total Base Digits** + **Delta Required Actual Frac Digits** =
   48 + 7 + 9 = 64

**Actual Results (Fractional digits of accuracy = 100)**
186 790 307.754 125 941 479 374 351 002 237 604 814 886 055 744 100 0000000000000000000000000000000000000000000000000000000000
Actual Results produce a total of 51 or 49-significant digits. (Integer Digits= 9 Fractional Digits= 42 or 40 )
We Planned for Worst Case: 64-Digits - No loss of Digits or Accuracy

<u>**⇓ Expected Result to 14-Fractional Digits**</u> 

186 790 307.754 125 941 479 37
186 790 307.754 125 941 479 374 351 002 237 604 814 886 055 744 100 0000000000000000000000000000000000000000000000000000000000
 <u>**⇑ Actual Result**</u>



