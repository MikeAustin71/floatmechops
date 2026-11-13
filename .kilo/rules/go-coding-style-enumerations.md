```text
## Agent Directive
> **MANDATORY**: Writing and/or modifying Go source code must be conducted in strict accordance with the rules defined in this document. Never bypass these rules.

---
```



# Go Coding Style Rule: Enumerations (Enums)

## Objective
When generating or refactoring Go code, AI Agents must adhere to the following conventions for declaring, initializing, and managing enumerated types (enums).

---

## Rule Specifications

### 1. Custom Type Definition
* Always define a custom underlying integer type (typically `int` or `uint8`) for the enumeration:
  ```go
  type <TypeName> int
  ```

### 2. Numbering Scheme
* **Always begin enumerations with a value of `1`** using `iota + 1`.
* The `0` value is strictly reserved for an uninitialized, default, or invalid state.

### 3. Constant Naming Conventions
* **Exported (Public) Constants**: Use **PascalCase**.
  
  ```go
  const MaxPacketSize = 512
  ```
* **Unexported (Private) Constants**: Use **lowerCamelCase**.
  ```go
  const maxBufferSize = 1024
  ```
* Enum constants should follow this visibility standard depending on whether they need to be accessed outside their package.

### 4. Accompanying String Array
* Every enum declaration **must** be accompanied by a fixed-size string array (or slice) containing the textual representations of the enum constants.
* Because enums start at index `1`, index `0` of the string array **must always be `"INVALID"`**.
* Implement a `String()` method on the enum type implementing `fmt.Stringer` that safely indexes into this array and returns `"INVALID"` (index 0) for out-of-range values.

---

## Reference Example

Below is the required implementation pattern for a public enum named `RoundingType` containing the elements `AwayFromZero`, `BankersRounding`, and `Truncation`.

```go
package mathutil

import "fmt"

// RoundingType defines the strategy used for numeric rounding.
type RoundingType int

// Exported RoundingType enum values starting at 1.
const (
	AwayFromZero RoundingType = iota + 1
	BankersRounding
	Truncation
)

// RoundingTypeNames maps RoundingType values to their string representations.
// Index 0 represents the uninitialized / invalid state.
var RoundingTypeNames = [...]string{
	"INVALID",
	"AwayFromZero",
	"BankersRounding",
	"Truncation",
}

// IsValid
//
// Returns a boolean value. 
// If the RoundingType value is valid, returns 'true'.
// If the RoundingType value is invalid, returns 'false'. 
func (r RoundingType) IsValid() bool {
    
 	if int(r) < 1 || int(r) >= len(RoundingTypeNames) {
		return false
	}
   
    return true
}


// String implements the fmt.Stringer interface for RoundingType.
func (r RoundingType) String() string {
	if int(r) < 1 || int(r) >= len(RoundingTypeNames) {
		return RoundingTypeNames[0]
	}
	return RoundingTypeNames[r]
}
```

---

## Agent Enforcement Checklist
- [ ] Enum starts at `1` using `iota + 1` (never starts at `0`).
- [ ] Exported constants use PascalCase; unexported constants use lowerCamelCase.
- [ ] A corresponding string array is defined.
- [ ] Index `0` of the string array is explicitly set to `"INVALID"`.
- [ ] A safe `String()` method bounds-checks the value and falls back to index `0`.