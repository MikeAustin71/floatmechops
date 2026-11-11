# Golang Error Handling Guidelines

## Agent Directive
> **MANDATORY**: When writing or modifying Go source code, all error handling operations must be structured in strict accordance with the rules defined in this document. Never bypass these rules.

---

## 1. Core Error Handling Rules

*   **Handle All Errors**: Every error returned by a function must be explicitly checked and handled.
*   **Never Discard Errors**: You must NEVER discard an error using the blank identifier (`_`). 
    *   *Violation:* `result, _ := someFunction()`
    *   *Requirement:* `result, err := someFunction(); if err != nil { ... }`

---

## 2. The `ErrorGenericFactory`

When writing code to handle errors, you must use the `ErrorGenericFactory` type to add error context and return errors to higher-level methods and functions.

### Import Path
You must import the `errpref` package to access the factory:
```go
import (
	ePref "github.com/MikeAustin71/errpref/v2"
)
```



### Factory Initialization

At the top of a function that handles or returns errors, initialize the config using ***NewErrorGenericConfig***:

```go
// Signature
// func (egc *ErrorGenericFactory) NewErrorGenericConfig(sourceFile string, parentType string, thisFunction string) *ErrorGenericFactory

errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
    "filename.go",      // sourceFile
    "parentTypeName",   // parentType (leave "" if none)
    "functionName()")   // thisFunction
```



## 3. Standardized Error Types (API Reference)

The ErrorGenericFactory provides four specialized custom error generation methods. You must use the appropriate method based on the error context:

### A. Wrapping an Error from another Function Call

Use ***GetErrorFromFunctionCall*** to add context to an error returned by a subordinate function.

```go
// func (egc *ErrorGenericFactory) GetErrorFromFunctionCall(calledFunctionType string, calledFunction string, errorContext string, errorMessage string, returnedCalledFunctionError error) error
```

### B. Parameter Out of Range

Use ***GetErrorInputParameterOutOfRange*** when an input validation fails due to boundaries/limits.

```go
// func (egc *ErrorGenericFactory) GetErrorInputParameterOutOfRange(parameter string, path string, fileName string, url string, errorContext string, errorMessage string) error
```

### C. Parameter is Nil Pointer

Use ***GetErrorInputParameterNilPointer*** when a required pointer argument is nil.

```go
// func (egc *ErrorGenericFactory) GetErrorInputParameterNilPointer(parameter string, path string, fileName string, url string, errorContext string, errorMessage string) error
```

### D. Native/Internal Error

Use ***GetErrorNativeToFunction*** for errors originating in the current function's internal logic.

```go
// func (egc *ErrorGenericFactory) GetErrorNativeToFunction(parameter string, path string, fileName string, url string, errorContext string, errorMessage string) error
```



## 4. Comprehensive Implementation Example

```go
package main

import (
	"fmt"
	"math"

	ePref "github.com/MikeAustin71/errpref/v2"
)

type mathLogic struct{}

// Example 1: Validating inputs and wrapping lower-level errors
func (ml *mathLogic) processValues(value1 *int64, value2 int32) error {
	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"main.go",
		"mathLogic",
		"processValues()")

	// 1. Handle Nil Pointer
	if value1 == nil {
		return errGenericCfg.GetErrorInputParameterNilPointer(
			"value1", "", "", "",
			"Invalid method input parameter",
			"Error: Input parameter 'value1' is a nil pointer")
	}

	// 2. Handle Out of Range
	if *value1 > math.MaxInt32 {
		return errGenericCfg.GetErrorInputParameterOutOfRange(
			fmt.Sprintf("value1=%d", *value1), "", "", "",
			fmt.Sprintf("'value1' exceeds maximum value of type 'int32' (%d)", math.MaxInt32),
			"Error: 'value1' is out of range")
	}

	// 3. Handle and Wrap Subordinate Function Errors
	convertedIntValue1 := int32(*value1)
	_, err := ml.addTwoInt32s(convertedIntValue1, value2)
	
	if err != nil {
		return errGenericCfg.GetErrorFromFunctionCall(
			"mathLogic",
			"addTwoInt32s()",
			"", "",
			err)
	}

	return nil
}

// Example 2: Generating a native logic error
func (ml *mathLogic) addTwoInt32s(value1 int32, value2 int32) (int32, error) {
	errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"main.go",
		"mathLogic",
		"addTwoInt32s()")

	int64Result := int64(value1) + int64(value2)

	// 4. Handle Native Function Logic Errors
	if int64Result > math.MaxInt32 {
		return 0, errGenericCfg.GetErrorNativeToFunction(
			"", "", "", "",
			fmt.Sprintf("Addition result ('%d') exceeds max 'int32' (%d)", int64Result, math.MaxInt32),
			"Addition result out of range")
	}

	return int32(int64Result), nil
}
```
