# Go Coding Test Guidelines

## Agent Directive

> **MANDATORY**: Writing and/or modifying Go Test Functions must be conducted in strict accordance with the rules defined in this document. Never bypass these rules.

---

## 1. Tests are Required

### Public Functions and Public Methods on Custom Types

All  public functions and public methods on custom types must be covered by tests.



## 2. All Failing Tests Must Be Identified

All failing tests must be identified and reported.



## 3. All ***t.Errorf()*** Calls MUST be Followed by a ***return*** Statement.

### INCORRECT FORMAT

In the example test below, two ***t.Errorf*** statements are configured **without** return statements. This is practice **<u>INCORRECT</u>** and **<u>NOT ALLOWED</u>**.

```go
package main

import (
	"encoding/json"
	"testing"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestUnmarshalUser(t *testing.T) {
	rawJSON := []byte(`{"name": "Alice", "age": 30}`)
	var user User

	// Use t.Fatalf here because if unmarshaling fails, 
	// checking user.Name or user.Age below is pointless.
	if err := json.Unmarshal(rawJSON, &user); err != nil {
		t.Fatalf("json.Unmarshal() failed unexpectedly: %v", err)
	}

	// WRONG! INCORRECT! 
    // missing: return
	if user.Name != "Alice" {
		t.Errorf("got Name = %q, want %q", user.Name, "Alice")
	}
	
	// WRONG! INCORRECT! 
    // missing: return
	if user.Age != 30 {
		t.Errorf("got Age = %d, want %d", user.Age, 30)
	}
}
```

### CORRECT FORMAT

All ***t.Errorf()*** statements <u>**must be followed** by a ***return*** statement</u>.

```go
package main

import (
	"encoding/json"
	"testing"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestUnmarshalUser(t *testing.T) {
	rawJSON := []byte(`{"name": "Alice", "age": 30}`)
	var user User

	// Use t.Fatalf here because if unmarshaling fails, 
	// checking user.Name or user.Age below is pointless.
	if err := json.Unmarshal(rawJSON, &user); err != nil {
		t.Fatalf("json.Unmarshal() failed unexpectedly: %v", err)
	}

	// CORRECT! 
    // Upon Error, Exit function with 'return' statement
	if user.Name != "Alice" {
		t.Errorf("got Name = %q, want %q", user.Name, "Alice")
        return
	}
	
	// CORRECT! 
    // Upon Error, Exit function with 'return' statement
	if user.Age != 30 {
		t.Errorf("got Age = %d, want %d", user.Age, 30)
        return
	}
}
```



## 4. Naming Test Functions

Test Function names should describe the function, method or procedure that function is testing. 

If multiple Test Functions are configured to describe a single function, method or procedure, those Test Function Names should include a numbering component consisting of 5-digits (**`_00000`**). Numbering components should begin with ten (**`_00010`**) and increment by ten (**`_00010`**, **`_00020`**, **`_00030`**, **`_00040`**, **`_00050`**,)  