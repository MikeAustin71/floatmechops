# Golang Coding and Style Guidelines: Temporary or Scratch Files

## Agent Directive

> **MANDATORY**: Creating, writing and/or modifying temporary or scratch files must be conducted in strict accordance with the rules defined in this document. Never bypass these rules.

---



## Temporary or Scratch Files 

Often the need arises to create one or more Temporary or Scratch files in order to further the goals and objectives of the assigned task. Temporary or Scratch files might be used to test algorithms before final implementation, validate calculation outcomes without modifying exiting source code or to experiment, in a safe and controlled manner, with new ideas,  algorithms or proposed solutions. 

As such, Temporary or Scratch files are typically created for an immediate specific purpose, and later deleted or overwritten when they no longer serve a useful function. However, when a specific need or project benefit arises, Temporary or Scratch files my be retained indefinitely.



## Directories Reserved for Temporary or Scratch Files

Temporary or Scratch Files may be created in the following Project Directories two directories:

**[Project Base Directory]/AppTest**

**[Project Base Directory]/py_tools**



**NEVER** create, modify, or delete temporary or scratch files in directories outside the **[Project Base Directory]**.



### [Project Base Directory]/AppTest

Typically ***main.go*** and other related test source code files may be placed in this directory. 

As shown in the following example, configured with the proper ***import*** statement, the ***Base_Project_Directory/AppTest/main.go*** file can call any public method or function in the primary ***source code directory***.

In this example, the file ***floatmechops/AppTest/main.go*** contains a ***main*** function which calls a ***test*** function which in turn calls a public method in the ***source code directory***. 

```go
package main

// Location of 'main.go': /MikeAustin71/floatmechops/AppTest/main.go
// Source Code Directory: 'floatmech'

// 'import' statement allows testSetup() to call public
// methods and functions stored in source code files
// in the directory: /MikeAustin71/floatmechops/floatmech
// These public methods and functions may now be used to 
// test and evaluate live code.
// ----------------------------------------------------
import "github.com/MikeAustin71/floatmechops/floatmech"

func main() {
	err := testSetup()
	
	if err != nil {
		fmt.Printf("testSetup err: %v\n", err)
		return
	}
	
	fmt.Println("testSetup success")
	
	return
}

func testSetup() error {
	var vNumStr string
	var targetFractionalDigits uint
	var mode big.RoundingMode

	vNumStr = "6.9857"


	targetFractionalDigits = 8
	mode = big.ToNearestAway

	decimal, err := new(floatmech.Decimal).NewFromFloatNumStr(vNumStr, targetFractionalDigits, mode)

	if err != nil {
		err2:= fmt.Errorf("NewFromFloatNumStr() failed: %s", err.Error())

		return err2
	}

	fmt.Printf("decimal = %v\n", decimal.String())

	return nil
}


```



You are authorized to create temporary or scratch files in the ***AppTest*** directory or create new subdirectories and files underneath the ***AppTest*** directory.



### [Project Base Directory]/py_tools

The ***py_tools*** directory is configured for ***python*** source code and programming operations. You are authorized to create new temporary or scratch files as well as permanent python files in this directory. In addition, you may also create subdirectories in this directory.  