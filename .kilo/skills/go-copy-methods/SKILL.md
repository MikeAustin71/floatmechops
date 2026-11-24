---
name: go-copy-methods
description: Use this skill ONLY when configuring copy and clone methods on custom types. 
---



# Golang CopyIn(), CopyOut() and Clone() Method Configuration Guidelines

## Agent Directive
> **MANDATORY**: When public methods are configured on custom types with sync.Mutex protection, those custom types will require three separate methods to perform copy operations. Those separate methods will be named CopyIn, CopyOut, and Clone. Never bypass these rules.



## Background - Copy Methods Are Required

When public methods are configured on custom types, and are additionally configured sync.Mutex objects, they cannot be copied using a simple copy operation. 

```go
// Custom Type configured with a sync.Mutex
type Decimal struct {
	mu                     sync.Mutex
	val                    *big.Float
	targetFractionalDigits uint
	prec                   uint
	mode                   big.RoundingMode
}

...

func (d *Decimal) TargetFractionalDigits() uint {

	if d == nil {
		return 0
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	return d.targetFractionalDigits
}
```

 In this configuration, with ***sync.Mutex*** objects , the solution needed to perform copy operations is provided by three separate public methods: **CopyIn()**, ***CopyOut()*** and ***Clone()***.

## CopyIn()

> **NOTE TO AGENT - MANDATORY**: 
> For the CopyIn() method: <u>Never</u> copy pointers to pointers. Always create a new object with equal values and set the receiving pointer to that new object.
>
> ```go
> d.val = new(big.Float).SetMode(copySource.val.Mode()).Set(copySource.val)
> ```
>
> 

```go

// CopyIn
//
// Receives an incoming 'Decimal' object (copySource)
// and copies those structure member variable values into
// the structure member variables of the host, or receiving
// 'Decimal' object.
//
// Input Parameters
// ----------------
//
//  copySource *Decimal
//    Member variable values from this Decimal object
//    will be copied into corresponding member variables
//    in the host, or receiving Decimal object.
//
// Return Parameters
// -----------------
//
//  error 
//  
func (d *Decimal) CopyIn(copySource *Decimal) error {
	
	// NewErrorGenericConfig()
    //
    // func (egc *ErrorGenericFactory) NewErrorGenericConfig(
    //    sourceFile string,
    //    parentType string,
    //    thisFunction string) *ErrorGenericFactory
    //
  	// Add identification information for this function
    errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"decimal_0010.go",
		"Decimal",
		"CopyIn()")
	
	if d == nil {
	
        // func (egc *ErrorGenericFactory) GetErrorInputParameterNilPointer(
        //    parameter string,
        //    path string,
        //    fileName string,
        //    url string,
        //    errorContext string,
        //    errorMessage string) error
        //
        // Formats erorr message and error context for when nil pointers
        // are passed as Function input parameters.
		return errGenericCfg.GetErrorInputParameterNilPointer(
			"d == nil",
			"",
			"",
			"",
			"Invalid Pointer Receiver",
			"Error: Pointer Receiver 'd' is a nil pointer")
	
	}

	if copySource  == nil {
	
        // func (egc *ErrorGenericFactory) GetErrorInputParameterNilPointer(
        //    parameter string,
        //    path string,
        //    fileName string,
        //    url string,
        //    errorContext string,
        //    errorMessage string) error
        //
        // Formats erorr message and error context for when nil pointers
        // are passed as Function input parameters.
		return errGenericCfg.GetErrorInputParameterNilPointer(
			"copySource == nil",
			"",
			"",
			"",
			"copySource is a nil pointer",
			"Error: Input parameter 'copySource' is a nil pointer")
	}


	d.val = new(big.Float).SetMode(copySource.val.Mode()).Set(copySource.val)
    
	d.targetFractionalDigits = copySource.targetFractionalDigits
    
	d.prec = copySource.prec
    
	d.mode = copySource.mode

	return nil
}
```





## CopyOut()

> **NOTE TO AGENT - MANDATORY**: 
> For the CopyOut() method: <u>Never</u> copy pointers to pointers. Always create a new object with equal values and set the receiving pointer to that new object. 
>
> ```go
> copyDestination.val = new(big.Float).SetMode(d.val.Mode()).Set(d.val)
> ```



```go

// CopyOut
//
// Receives an outgoing 'Decimal' object (copyDestination)
// and copies the host Decimal structure member variable values into
// the structure member variables of this copyDestination 'Decimal'
// object.
//
// Input Parameters
// ----------------
//
//  copyDestination *Decimal
//    Member variable values from the host (source) Decimal
//    object will be copied into corresponding member variables
//    of the copyDestination Decimal object.
//
// Return Parameters
// -----------------
//
//  error 
//  
func (d *Decimal) CopyOut(copyDestination *Decimal) error {
	
	// NewErrorGenericConfig()
    //
    // func (egc *ErrorGenericFactory) NewErrorGenericConfig(
    //    sourceFile string,
    //    parentType string,
    //    thisFunction string) *ErrorGenericFactory
    //
  	// Add identification information for this function
    errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"decimal_0010.go",
		"Decimal",
		"CopyOut()")
	
	if d == nil {
	
        // func (egc *ErrorGenericFactory) GetErrorInputParameterNilPointer(
        //    parameter string,
        //    path string,
        //    fileName string,
        //    url string,
        //    errorContext string,
        //    errorMessage string) error
        //
        // Formats erorr message and error context for nil pointers
        // passed as Function input parameters.
		return errGenericCfg.GetErrorInputParameterNilPointer(
			"d == nil",
			"",
			"",
			"",
			"Invalid Pointer Receiver",
			"Error: Pointer Receiver 'd' is a nil pointer")
	
	}

	if copyDestination  == nil {
	
        // func (egc *ErrorGenericFactory) GetErrorInputParameterNilPointer(
        //    parameter string,
        //    path string,
        //    fileName string,
        //    url string,
        //    errorContext string,
        //    errorMessage string) error
        //
        // Formats erorr message and error context for when nil pointers
        // are passed as Function input parameters.
		return errGenericCfg.GetErrorInputParameterNilPointer(
			"copyDestination == nil",
			"",
			"",
			"",
			"copyDestination is a nil pointer",
			"Error: Input parameter 'copyDestination' is a nil pointer")
	}


    copyDestination.val = new(big.Float).SetMode(d.val.Mode()).Set(d.val) 
    
	copyDestination.targetFractionalDigits = d.targetFractionalDigits
    
	copyDestination.prec = d.prec
    
	copyDestination.mode = d.mode

	return nil
}
```



## Clone()

> **NOTE TO AGENT - MANDATORY**: 
> For the Clone() method: <u>Never</u> copy pointers to pointers. Always create a new object with equal values and set the receiving pointer to that new object. 
>
> ```go
> clonedDecimal.val = new(big.Float).SetMode(d.val.Mode()).Set(d.val)
> ```
>
> 

```go
// Clone
//
// Returns a newly created *Decimal object which is an
// exact clone of the host Decimal object.
//
// Input Parameters
// ----------------
//
//  None
//
//
// Return Parameters
// -----------------
//
//  *Decimal
//    Pointer to a newly created Decimal object containing duplicate 
//    member variable values to those in the host Decimal object.
//
//  error
//
func (d *Decimal) Clone() (*Decimal, error) {
	
	// NewErrorGenericConfig()
    //
    // func (egc *ErrorGenericFactory) NewErrorGenericConfig(
    //    sourceFile string,
    //    parentType string,
    //    thisFunction string) *ErrorGenericFactory
    //
  	// Add identification information for this function
    errGenericCfg := new(ePref.ErrorGenericFactory).NewErrorGenericConfig(
		"decimal_0010.go",
		"Decimal",
		"CopyOut()")
	
	if d == nil {
	
        // func (egc *ErrorGenericFactory) GetErrorInputParameterNilPointer(
        //    parameter string,
        //    path string,
        //    fileName string,
        //    url string,
        //    errorContext string,
        //    errorMessage string) error
        //
        // Formats erorr message and error context for when nil pointers
        // are passed as Function input parameters.
        return &Decimal{}, errGenericCfg.GetErrorInputParameterNilPointer(
			"d == nil",
			"",
			"",
			"",
			"Invalid Pointer Receiver",
			"Error: Pointer Receiver 'd' is a nil pointer")
	
	}

    clonedDecimal := new(Decimal)
        
	clonedDecimal.val = new(big.Float).SetMode(d.val.Mode()).Set(d.val)

	clonedDecimal.targetFractionalDigits = d.targetFractionalDigits
    
	clonedDecimal.prec = d.prec
    
	clonedDecimal.mode = d.mode

	return clonedDecimal, nil
}
```

