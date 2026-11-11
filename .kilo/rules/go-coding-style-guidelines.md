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
var calculateLCDlock *sync.Mutex // Exclusive to CalculateLCD

func CalculateLCD(a, b int) (int, error) {
	if calculateLCDlock == nil {
		calculateLCDlock = new(sync.Mutex)
	}
	calculateLCDlock.Lock()
	defer calculateLCDlock.Unlock()
    // ... logic ...
}
```
