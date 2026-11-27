<!-- AGENTS.md at project root -->

# Project Workspace Rules

## **Filesystem Boundary** 

- Never create, modify, or delete files outside the active project root directory. Do not write temporary files to `%TEMP%` or `AppData`.
- To Create temporary files -> load skill `go-coding-style-tempory-files`

## Skill Routing
- **Go (`.go` / `go.mod`) tasks**: Load skill `go-127-standards` before writing or refactoring Go code.
- **Go `math` or `math/big` tasks**: Load skill `go-math-packages` when importing or using the `math` or `math/big` packages.

## Skill Triggers

-- Writing or modifying ***go*** custom types or methods  -> load skill `go-coding-style-custom-types-and-methods`

-- Writing or modifying ***go*** enumerations (a.k.a. ***enums***) -> load skill `go-coding-style-enumerations`

-- Writing, modifying or reviewing ***go*** error handling operations -> load skill `go-coding-style-error-handling` 

-- Writing or modifying ***go*** functions, methods, types, interfaces or packages -> load skill `go-coding-style-source-code-comments`

-- Writing or modifying ***go*** stand alone functions -> load skill `go-coding-style-stand-alone-functions`

-- Creating or modifying temporary files -> load skill `go-coding-style-tempory-files`

-- Writing or modifying ***go*** tests functions or test files -> load skill `go-coding-test-guidelines`

 -- Creating, naming or modifying directories and/or files in the project directory  -> load skill `go-file-naming-conventions`

-- Writing, modifying or executing ***python*** scripts -> load skill `python-operations`





