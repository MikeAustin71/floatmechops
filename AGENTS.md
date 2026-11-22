<!-- AGENTS.md at project root -->

# Project Workspace Rules



## Go v1.27 Development Standards

- **Context Grounding**: All generated Go code must strictly adhere to Golang Version 1.27 features, idioms, and standard library conventions.

- **Knowledge Retrieval (Active Search)**: If you are unsure about Go 1.27 specific syntax, or before writing/refactoring complex logic, you must actively query the local knowledge base using your file search or shell tools:

- For incremental feature updates or version differences, search the **`F:/Go_Knowledge/releases/`** directory.

- For core language mechanics, read the **`F:/Go_Knowledge/language-spec.md/language-spec.md`** file.

- For official documentation on the Go math/big package, read the **`F:/Go_Knowledge/packages/math_big_go1.27.0.md`** file.

- For official documentation on the Go math package, read the **`F:/Go_Knowledge/packages/math_go1.27.0.md`** file.

- **Strict Compliance**: Reject and avoid any deprecated Go functions or patterns outlined as obsolete in the local documentation.

- **Tooling First**: Do not guess version-specific behaviors. Rely on reading the specified knowledge paths to guarantee factual accuracy before outputting code.

- The **Go** Version Number in the Project Directory **`go.mod`** have may have a minor release number from that of the **Go** documentation (Example: ***`go.mod` = Go Version 1.27.3*** - **Go Documentation Version 1.27.0**). This discrepancy is NOT significant. You are directed to use the **Go Documentation** located in the **`F:/Go_Knowledge`** directories.

  

## Python Environment
- This repository is primarily written in Go, with a dedicated Python workspace in `py_tools/`.
- Python 3.14.10 is configured in a virtual environment at `py_tools\.venv\`.
- **Always** execute Python and pip using the virtual environment's executable directly from the project root (do not call global `python`, `python3`, or `pip`):
  - Run Python: `.\py_tools\.venv\Scripts\python.exe <script.py>`
  - Run pip: `.\py_tools\.venv\Scripts\python.exe -m pip <command>`
- Place all Python scripts and dependencies (`requirements.txt`) inside `py_tools/`.
- Do not read, search, or modify files inside `py_tools\.venv\`.
