# Golang Coding and Style Guidelines

## Agent Directive
> **MANDATORY**: Writing and/or modifying Go source code must be conducted in strict accordance with the rules defined in this document. Never bypass these rules.

---

## 1. Spacing and Formatting
*   **Statement Spacing**: Insert one blank line before and after each source code statement.
*   **Block Exceptions**: `if`, `for`, and `switch` statements with a single interior line do not need blank lines inside the block. Insert a blank line before the block and after the closing bracket.
*   **Variable Blocks**: Multiple `var` definitions grouped together require only one blank line before the first `var` and one after the last `var`.

---

## 2. Naming Conventions and Capitalization
*   **Packages**: `lowercase` (e.g., `agm`, `taylor`).
*   **Exported (Public)**: `PascalCase` for Constants, Functions, Methods, Types, and Variables.
*   **Unexported (Private)**: `camelCase` for Constants, Functions, Methods, Types, and Variables.
*   **Acronyms**: Always preserve common acronyms in ALL CAPS (e.g., `URLParser`, `HTTPClient`, `ParseJSON()`). Do NOT use `UrlParser` or `ParseJson`.

---

## 3. Custom Types and Object-Orientation Rules

### **A. Pointer Receivers Are Mandatory**

Custom Types MUST always be configured with **pointer receivers** for *both* public and private methods. `func (c *CustomType) MethodName()`.

### **B. Thread Safety (Mutex) on Public Methods**

Public methods on Custom Types act as entry points and MUST be configured with a private `sync.Mutex` object to ensure thread-safe operations.

```go
type StrMech struct {
	stringDataMutex *sync.Mutex // Private mutex for thread safety
}

func (sMech *StrMech) PrintMessage(messageStr string) error {
	if sMech.stringDataMutex == nil {
		sMech.stringDataMutex = new(sync.Mutex)
	}
	sMech.stringDataMutex.Lock()
	defer sMech.stringDataMutex.Unlock()
    
    // ... logic ...
}
```
### **C. Method Calling Constraints**

*   **NO Public-to-Public Calls**: Public Methods configured on a specific Custom Type **MUST NOT** call other Public Methods on that same instance. No exceptions.
*   **Separation of Concerns**: Custom Types should generally contain either all Public Methods or all Private Methods. If a Public Method needs a helper, put the Private Method in a separate Custom Type/file. 
*   **Private Methods**: Private Methods on Custom Types do NOT normally require a Mutex.

---

## 4. Source Code Comments

**Mandatory Requirement**: Source Code Comments are required for all new Functions, Methods, Custom Types, and Interfaces. Explain *What* it is, *Why* it exists, and *Contracts/Assumptions*. Do not explain code line-by-line.

### **A. Type & Interface Comments**

Must begin with the Type/Interface name, form a complete sentence, and describe the purpose/behavior.

```go
// NaturalLogCalculator
//
// Computes natural logarithms using AGM, Newton, or Taylor numerical methods.
type NaturalLogCalculator struct { ... }
```

### **B. Function & Method Comments**

Must follow this exact structure, including explicit "Input Parameters" and "Return Parameters" sections:

```go
// NaturalLog
//
// Returns ln(x). The function requires x > 0.
//
// Input Parameters
// ---------------
//
// x float64
//  'x' is the input value for which the Natural Log will be calculated.
//
// Return Parameters
// -----------------
//
// float64 
//  Returns the Natural Log value for input parameter 'x'.
func NaturalLog(x float64) float64 { ... }
```

---

## 5. Standalone (Global) Functions

*   **Discouraged**: The use of Global Public and Global Private Functions is heavily discouraged. Create methods on Custom Types instead.
*   **If Strictly Required**: If operational circumstances dictate a Global Public Function must be used, it MUST be protected by a dedicated package-level private `sync.Mutex`.
```go
var calculateLCDlock sync.Mutex // Exclusive to CalculateLCD

func CalculateLCD(a, b int) (int, error) {
	
	calculateLCDlock.Lock()
	
    defer calculateLCDlock.Unlock()
    // ... logic ...
}
```

## 6. Types and Methods

The structure of methods on custom types is the standard means for providing source code functionality in ***go*** projects. The use of methods on custom types is preferred to standalone functions.

### Public Methods on Custom Types

Public methods on custom types are typically designed as public entry points to calculations, procedures and operations. This structure is preferred over global public standalone functions.

#### Key Points:

1. Public methods on custom types are configured with mutex objects to protect the integrity of calculations, procedures and operations, as well as guard against race conditions. 

2. Public methods on custom types always use Pointer Receivers.

3. Public methods on custom types are configured as public entry points to calculations, procedures and operations. As such they perform initial input parameter validation and call internal private methods on custom types located in a separate source code file to perform the actual calculation. 

4. Public methods on custom types are treated as ***top level*** methods and located in a source code file with **_0010** file number component.

   Example:

   ```bash
   # Public entry point to lowest common denominator caclulation is located in source code file: LowestCommonDenom_0010.go. 
   
   # The mechanics of the lowest common denominator calculation are configured in a private method on
   # a custom type located in source code file: LowestCommonDenom_0020.go.
   
   ```

   

#### Example Public Method on Custom Type

```go
// This method located in source code file LowestCommonDenom_0010.go


import (
    // Required for error handling
	ePref "github.com/MikeAustin71/errpref/v2"
)

type LCDCalculation struc {
	mu *sync.Mutex    
}

// CalculateLCD
//
// Calculates the Lowest Common Denominator for input
// parameters, 'a' and 'b'.
//
// Input Parameters
// ----------------
//
// a int
//  One of two 'int' values which be used to compute the
//  Lowest Common Denominator. If 'a' is zero, an error
//  is returned. 
//
// b int
//  The second of two 'int' values which will be used to
//  compute Lowest Common Denominator. If 'b' is zero, an
//  error is returned.
//
// Return Parameters
// -----------------
//
// int
//  The calculated Lowest Common Denominator based on input
//  parameters 'a' and 'b'.
func (LCDcalc *LCDCalculation) CalculateLCD(a, b int) (int, error) {
 
    LCDcalc.mu.Lock()
    
    defer LCDcalc.mu.Unlock()
    
	// NewErrorGenericConfig()
    //
    // func (egc *ErrorGenericFactory) NewErrorGenericConfig(
    //    sourceFile string,
    //    parentType string,
    //    thisFunction string) *ErrorGenericFactory
    //
  	// Add identification information for this function
    errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"CalculateLCD_0010.go",
		"",
		"CalculateLCD()")
    
    if a == 0 {
        
        
		//func (egc *ErrorGenericFactory) GetErrorInputParameterOutOfRange(
		//  parameter string,
		//  path string,
		//  fileName string,
		//  url string,
		//  errorContext string,
		//  errorMessage string) error
		err := errGenericCfg.GetErrorInputParameterOutOfRange(
			fmt.Sprintf("a=%d", a),
			"",
			"",
			"",
			"a is equal to zero",
			"Error: 'a' is invalid and out of range")

		return -1, err
	}
    
    if b == 0 {
        
        
		//func (egc *ErrorGenericFactory) GetErrorInputParameterOutOfRange(
		//  parameter string,
		//  path string,
		//  fileName string,
		//  url string,
		//  errorContext string,
		//  errorMessage string) error
		err := errGenericCfg.GetErrorInputParameterOutOfRange(
			fmt.Sprintf("b=%d", b),
			"",
			"",
			"",
			"b is equal to zero",
			"Error: 'b' is invalid and out of range")

		return -1, err
	}
    
    
    lcdValue, err := new(lcdCalcMechanics).calcLCD(a, b)
    
    if err != nil {
        
        // func (egc *ErrorGenericFactory) GetErrorFromFunctionCall(
        //    calledFunctionType string,
        //    calledFunction string,
        //    errorContext string,
        //    errorMessage string,
        //    returnedCalledFunctionError error)
        //
		// Adds error context to errors returned by other 
        // Functions.
  		err2 := errGenericCfg.GetErrorFromFunctionCall(
			"lcdCalcMechanics",
			"calcLCD()",
			"",
			"",
			err)
       
        return -1, err2
    }
        
    return lcdValue, nil
}


```

###  Private Methods on Custom Types

Public methods on custom types are typically designed as private entry points to source code which actually performs calculations, procedures and operations. This structure is preferred over global private standalone functions.

#### Key Points

1. Private methods on custom types are typically **not** configured with mutex objects.

2. Private methods on custom types always use Pointer Receivers.

3. Private methods on custom types do the actual heavy lifting for calculations, procedures and operations. The internal mechanics of calculations, procedures and operations are typically performed by private methods on custom types. 

4. Public methods on custom types typically call private methods on custom types.

5. When private methods on custom types need to call helper methods, they generally other, lower level, private methods on custom types. 

6. Private methods on custom types are treated as ***lower level*** methods and located in a source code file with file number components beginning with **_0020**.

   

#### Example Private Method on Custom Type

```go
// This type and method are located in source code file LowestCommonDenom_0020.go


import (
    // Required for error handling
	ePref "github.com/MikeAustin71/errpref/v2"
)

// Type lcdCalcMechanics provides methods designed to 
// calculate the lowest common denominator.
type lcdCalcMechanics struc {}

// calculateLCD
//
// Calculates the Lowest Common Denominator for input
// parameters, 'a' and 'b'.
//
// Input Parameters
// ----------------
//
// a int
//  One of two 'int' values which be used to compute the
//  Lowest Common Denominator.
//
// b int
//  The second of two 'int' values which will be used to
//  compute Lowest Common Denominator.
//
// Return Parameters
// -----------------
//
// int
//  The calculated Lowest Common Denominator based on input
//  parameters 'a' and 'b'.
func (lcdCalcMech *lcdCalcMechanics) calcLCD(a, b) (int, error){
	
    // NewErrorGenericConfig()
    //
    // func (egc *ErrorGenericFactory) NewErrorGenericConfig(
    //    sourceFile string,
    //    parentType string,
    //    thisFunction string) *ErrorGenericFactory
    //
  	// Add identification information for this function
    errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"CalculateLCD_0020.go",
		"lcdCalcMechanics",
		"calcLCD()")
 
	if a == 0 || b == 0 {
		return 0
	}
    
	absA := a
    
	if absA < 0 {
		absA = -absA
	}
    
	absB := b
    
	if absB < 0 {
		absB = -absB
	}
    
    // call to private method on a custom type
    // Source located in LowestCommonDenom_0030.go
    gcdValue, err := new(GCDHelper).gcd(absA, absB)
    
    if err != nil {
       
        // func (egc *ErrorGenericFactory) GetErrorFromFunctionCall(
        //    calledFunctionType string,
        //    calledFunction string,
        //    errorContext string,
        //    errorMessage string,
        //    returnedCalledFunctionError error)
        //
		// Adds error context to errors returned by other 
        // Functions.
  		err2 := errGenericCfg.GetErrorFromFunctionCall(
			"GCDHelper",
			"gcd()",
			"",
			"",
			err)
       
        return -1, err2
    }
    
	return (absA * absB) / gcdValue, nil
}

```

