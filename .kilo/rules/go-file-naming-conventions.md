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
*   **File Numbering Component**: Typically begins with **`_0010.go`** and increases in increments of ten (10).
*   **Structure**: Usually contain public methods or functions serving as entry points into the calculation or operation.

### Example Hierarchy

```bash
NaturalLogAGM_0010.go # Top Level Types and Methods (Public Entry Points)
NaturalLogAGM_0020.go # Supporting Types and Methods
NaturalLogAGM_0030.go # Supporting Types and Methods
```

---

## Supporting Types and Methods

*   **File Base Name Component**: Shares the identical **File Base Name** as the top-level file.
*   **File Numbering Component**: Typically begins at **`_0020.go`** and increases in increments of ten (10).
*   **Structure**: Usually contain private helper methods called by higher-level types.

### Crowded Numbering & Extended Components

If a numbering sequence becomes crowded and new files must be inserted between existing files, use an extended numbering component format (`_000_010.go`):

```bash
[FILE BASE NAME]_0010.go 
[FILE BASE NAME]_0020.go 
[FILE BASE NAME]_0021.go 
[FILE BASE NAME]_0021_010.go # Extended File Number Component for insertions
[FILE BASE NAME]_0022.go 
[FILE BASE NAME]_0023.go 
```

---

# Part 2: Golang Test File Naming Conventions

## Overview

Test files contain code designed to test Go source code and are placed in the same directory as the source code files they test.

Test file names consist of **five components**:

1. **File Name Prefix**: Always **`zzzt_`** (ensures test files sort separately from source files).
2. **File Base Name**: Matches the target source subject.
3. **File Number Component**: Formatted as `_0000` (incremented by 10).
4. **File Test Label**: Always **`_test`**.
5. **File Extension**: Always **`.go`**.

### Test File Naming Formula

```bash
[FileNamePrefix] + [FileBaseName] + [FileNumberComponent] + [FileTestLabel] + [FileExtension]
```

**Example:**

```bash
zzzt_NaturalLogAGM_0010_test.go
zzzt_numberstrkernel_0010_test.go
zzzt_numberstrkernel_0020_test.go
```

### Crowded Numbering & Extended Components

If a numbering sequence becomes crowded and new files must be inserted between existing files, use an extended numbering component format (`_000_010.go`):

**Example**

```bash
zzzt_numberstrkernel_0010_test.go
zzzt_numberstrkernel_0020_test.go
zzzt_numberstrkernel_0021_test.go
zzzt_numberstrkernel_0021_010_test.go # Extended File Number Component for insertions
zzzt_numberstrkernel_0022_test.go
zzzt_numberstrkernel_0030_test.go
```



---

## Complete Project File Layout Example

This example includes both source code and test files. 

```bash
bigintfixeddecimal_0010.go
bigintfixeddecimal_0020.go
bigintfixeddecimal_0030.go
intary_0010.go
intary_0020.go
intary_0030.go
intary_0040.go
NaturalLogAGM_0010.go 
NaturalLogAGM_0020.go 
NaturalLogAGM_0030.go
NaturalLogAGM_0030_010.go # Extended File Number Component for insertions
NaturalLogAGM_0031.go 
NaturalLogAGM_0040.go 
NaturalLogAGM_0050.go 
zzzt_NaturalLogAGM_0010_test.go
zzzt_numberstrkernel_0010_test.go
zzzt_numberstrkernel_0020_test.go
```

