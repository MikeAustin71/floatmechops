# Golang Source Code & Test File Naming Conventions

## Agent Directive
> **MANDATORY**: When creating, organizing, or refactoring Go source files or test files in this project, you must strictly adhere to the hierarchical naming, numbering, prefix, and suffix rules defined in this document.

---

# Part 1: Golang Source Code File Naming Conventions

## Overview
This document describes the naming conventions applied to **Golang** source code files containing types and methods designed to perform a specific action, operation, or calculation. 

This convention establishes a loose file hierarchy consisting of a top-level source code file with public entry points and subsidiary levels containing source code files with private helper methods.

### Objectives and Structure
Source code files are grouped together to perform a single action, operation, or calculation are subdivided into two classes:
1. **Top Level Types, Methods and Functions**: Represent the top level of the hierarchy, containing primary entry points (usually public methods).
2. **Supporting Types, Methods and Functions**: Subsidiary files containing supporting logic (usually private helper methods).

### Source Code File Name Components
All source code file names consist of three components:

1. **File Base Name**: Derived from the subject action, operation, or calculation (e.g., `NaturalLogAGM`).
2. **File Number Component**: Formatted as `_0000` to designate relative hierarchy and order.
3. **File Extension**: Always `.go`.

```bash
[File Base Name] + [File Number Component] + [File Extension]
```

------

## Top Level Types and Methods

*   **File Base Name Component**: Derived from the objective or operational function of the file group (e.g., `NaturalLogAGM`).
*   **File Number Component**: Formatted with 7-digits `_0000000` (7-digits; incremented by 1000).
*   **Structure**: Usually contain public methods or functions serving as entry points into the calculation or operation.

### Example Hierarchy

```bash
naturalLogAGM_0001000.go # Top Level Types and Methods (Public Entry Points)
naturalLogAGM_0002000.go # Supporting Types and Methods
naturalLogAGM_0003000.go # Supporting Types and Methods
```

---

## Supporting Types and Methods

*   **File Base Name Component**: Shares the identical **File Base Name** as the top-level file.
*   **File Numbering Component**: Typically begins at **`_0002000.go`** and increases in increments of 1000 (one-thousand).
*   **Structure**: Usually contain private helper methods called by higher-level types.

### Crowded Numbering & Extended Numbering Components

If a numbering sequence becomes crowded and new files must be inserted between existing files, use an extended numbering component format (`_0000010.go`):

```bash
bigFloatPrecision_0001000.go
decimal_0001000.go
naturalLogAGM_0001000.go 
naturalLogAGM_0002000.go 
naturalLogAGM_0003000.go 
naturalLogAGM_0003010.go # Source File Extended File Number Component for insertions
naturalLogAGM_0003020.go # Source File Extended File Number Component for insertions
naturalLogAGM_0004000.go 
roundingType_0001000.go 
```

---

# Part 2: Golang Test File Naming Conventions

## Overview

Test files contain code designed to test Go source code and are placed in the same directory as the source code files they test.

Test file names consist of **five components**:

1. **File Name Prefix**: Always **`zzzt_`** (ensures test files sort separately from source files).
2. **File Base Name**: Matches the target source subject.
3. **File Number Component**: Formatted as `_0000000` (7-digits; incremented by 1000).
4. **File Test Label**: Always **`_test`**.
5. **File Extension**: Always **`.go`**.

### Test File Naming Formula

```bash
[FileNamePrefix] + [FileBaseName] + [FileNumberComponent] + [FileTestLabel] + [FileExtension]
```

**Example:**

```bash
zzzt_naturalLogAGM_0001000_test.go
zzzt_numberstrkernel_0001000_test.go
zzzt_numberstrkernel_0002000_test.go
```

### Crowded Numbering & Extended Numbering Components

If a numbering sequence becomes crowded and new files must be inserted between existing files, use an extended numbering component format (`_0001010_test.go`) (Increment :

**Example**

```bash
zzzt_naturalLogAGM_0001000_test.go
zzzt_numberstrkernel_0001000_test.go
zzzt_numberstrkernel_0002000_test.go
zzzt_numberstrkernel_0003000_test.go
zzzt_numberstrkernel_0003010_test.go # Test File Extended File Number Component for insertions
zzzt_numberstrkernel_0003020_test.go # Test File Extended File Number Component for insertions
zzzt_numberstrkernel_0003030_test.go # Test File Extended File Number Component for insertions
zzzt_numberstrkernel_0004000_test.go
zzzt_numberstrkernel_0005000_test.go
```

Note that Test file Extended Numbering Components differ from Source Code file Extended Numbering Components.



---

## Complete Project File Layout Example

This example includes both source code and test files. 

```bash
bigintfixeddecimal_0001000.go
bigintfixeddecimal_0002000.go
bigintfixeddecimal_0003000.go
intary_0001000.go
intary_0002000.go
intary_0003000.go
intary_0004000.go
naturalLogAGM_0001000.go 
naturalLogAGM_0002000.go 
naturalLogAGM_0003000.go
naturalLogAGM_0003010.go # Source File Extended File Number Component for insertions
NaturalLogAGM_0003020.go # Source File Extended File Number Component for insertions
NaturalLogAGM_0004000.go 
NaturalLogAGM_0005000.go 
zzzt_NaturalLogAGM_0001000_test.go
zzzt_numberstrkernel_0001000_test.go
zzzt_numberstrkernel_0002000_test.go
zzzt_numberstrkernel_0003000_test.go
zzzt_numberstrkernel_0003010_test.go # Test File Extended File Number Component for insertions
zzzt_numberstrkernel_0003020_test.go # Test File Extended File Number Component for insertions
zzzt_numberstrkernel_0003030_test.go # Test File Extended File Number Component for insertions
zzzt_numberstrkernel_0004000_test.go
zzzt_numberstrkernel_0005000_test.go
```

