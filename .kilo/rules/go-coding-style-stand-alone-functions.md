# Golang Coding and Style Guidelines - Standalone Functions

## Agent Directive

> **MANDATORY**: Writing and/or modifying Go source code must be conducted in strict accordance with the rules defined in this document. Never bypass these rules.

------



## 1. Usage 

- The use of **Global Public Functions** is discouraged. As a substitute for **Global Public Functions**, the practice of creating public methods on custom types is recommended. 
- The use of **Private Functions** is discouraged. As a substitute for **Private Functions**, the practice of creating private methods on custom types is recommended. 

------



## 2. Naming Functions

Use descriptive variable names when dealing with complex numerical algorithms, even when they're slightly longer than typical ***Go*** code.

------



## 3. Global Public Functions

In those cases where the use of **Global Public Functions** is dictated by operational circumstances or objectives, such **Global Public Functions** should be configured with a private mutex object as shown in the  following example:

```go
import (
    // Required for error handling
	ePref "github.com/MikeAustin71/errpref/v2"
)



// This sync.Mutex object is reserved for the 
// exclusive use of Global Public Function 'CalculateLCD'
calculateLCDLock sync.Mutex

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
func CalculateLCD(a, b int) (int, error) {
  
	calculateLCDLock.Lock()

	defer calculateLCDLock.Unlock()
    
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

------



## 4. Private Functions

In those cases where the use of **Private Functions** is dictated by operational circumstances or objectives, such **Private Functions** should be configured as shown in the  following example:

```go
import (
    // Required for error handling
	ePref "github.com/MikeAustin71/errpref/v2"
)

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
func calculateLCD(a int, b int) (int, error) {
	
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
		"",
		"calculateLCD()")
 
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

Note that the use of private sync.Mutex objects is not normally required for **Private Functions**.

