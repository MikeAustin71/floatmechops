

# How to Configure ***Python*** in ***Go*** Project



## My Configuration

### Operating System

- **Edition	Windows 11 Pro**
- **Version	25H2**

#### AI Software

- **Ollama Version: *0.40.17* Running Local LLM Models**
- **Agent: *Kilo* (Integrated in JetBrains Editors)**

#### Editors

- **JetBrains Editors** - ***Kilo*** is integrated in the editors. 

#### Local LLM's

- ***qwen3.8-27b*: Model Derivatives with 32kb and 64kb contexts.**
- ***gemma4:31b***: Model Derivatives with 32kb and 64kb contexts.



## Background

My ***Kilo*** agent has tried on several occasions to use python and has been unable to do so claiming it is not configured. ***Kilo*** operations are limited to the project directory where all of the code is written in the  ***go*** programming language.

**Python 3.14.10** is installed on my computer. 

I have a separate python project set up in a separate directory configured automatically by my ***JetBrains*** editor.  The ***PyCharm*** editor automatically configured  ***.venv*** in the top level directory of my python projects.

I would like to set up python for use by my agent. 

I was thinking that I should set up a separate directory inside of my ***go*** Project Directory structure. However, I'm guessing I would also have to set up ***.venv***



## Question

Please provide the steps necessary to manually set up and configure python in a separate directory inside my ***go*** Project Directory structure?



# Gemini 4 Hydrogen - Response No 1 - 2026-11-21 15:31-Hours

### Why This Happens

PyCharm automatically creates and attaches a `.venv` when you open a Python project. When you open a **Go** project in a JetBrains IDE (such as GoLand), no Python SDK or `.venv` is created inside that directory. 

When **Kilo** runs inside your Go project directory on Windows 11, it typically fails or reports that Python is not configured for two reasons:

1. **No `.venv` exists in the workspace**, and bare `python` / `python3` commands on Windows often aren't on `PATH` or hit Windows' built-in Microsoft Store stub (`WindowsApps\python.exe`).
2. **Local LLMs (`qwen3.8-27b`, `gemma4:31b`) only see your Go files** in the workspace, so unless a Python directory and an explicit instruction file tell Kilo where the interpreter lives, the agent assumes Python isn't set up.

---

### Recommended Directory Layout

Placing Python inside its own subfolder (for example, `py_tools/`) keeps your Go codebase clean. Go's toolchain (`go build ./...`, `go test ./...`) safely ignores subdirectories that contain no `.go` files as well as dot-prefixed folders like `.venv`.

```text
your-go-project/
├── AGENTS.md              <-- Tells Kilo where Python & .venv are
├── .gitignore             <-- Excludes .venv and cache files from Git
├── go.mod
├── main.go
└── py_tools/              <-- Dedicated Python directory
    ├── .venv/             <-- Isolated Python 3.14.10 environment
    └── requirements.txt   <-- Tracks installed Python packages
```

---

### Step-by-Step Setup

Open the **Terminal** in your JetBrains IDE (or Windows PowerShell) inside the root of your Go project.

#### 1. Locate Your Installed Python 3.14.10 Command

Depending on how Python 3.14.10 was installed on Windows 11 (standalone installer vs. Python Install Manager), verify which command invokes it:

```powershell
# Try the Windows Python launcher (either syntax depending on launcher version)
py -3.14 --version
# OR
py -V:3.14 --version
# OR
python --version
```

* **If `py` says Python isn't found, or `python` opens the Microsoft Store:** Check where PyCharm found your Python interpreter (usually `C:\Users\<YourName>\AppData\Local\Programs\Python\Python314\python.exe` or `C:\Users\<YourName>\AppData\Local\Python\pythoncore-3.14-64\python.exe`). You can use that full path in Step 3.

#### 2. Create the Python Subdirectory

From your Go project root:

```powershell
mkdir py_tools
cd py_tools
```

#### 3. Create the `.venv` Virtual Environment

Inside `py_tools`, create `.venv` using the working command from Step 1:

```powershell
# Use whichever command worked in Step 1:
py -3.14 -m venv .venv
```

*(Or `py -V:3.14 -m venv .venv`, or `"C:\full\path\to\python.exe" -m venv .venv`)*

#### 4. Verify `.venv` and Upgrade `pip`

You do **not** need to activate the virtual environment to use it—calling `.venv\Scripts\python.exe` directly is the most reliable method on Windows (both for you and for AI agents):

```powershell
.\.venv\Scripts\python.exe --version
.\.venv\Scripts\python.exe -m pip install --upgrade pip
```

It should output `Python 3.14.10`.

*(Optional)* If you also want to activate `.venv` manually in your interactive PowerShell terminal:

```powershell
.\.venv\Scripts\Activate.ps1
```

> **Windows 11 Tip:** If PowerShell blocks `Activate.ps1` with an execution policy error, run this once for your user account and re-try:
>
> ```powershell
> Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
> ```

Create an empty `requirements.txt` (or populate it after installing packages):

```powershell
New-Item -Path requirements.txt -ItemType File
```

---

### Step 5: Protect Git and JetBrains from `.venv` Bloat

1. **Update `.gitignore`** (in your Go project root) so you don't commit thousands of virtual environment files:

   ```gitignore
   # Python environment
   py_tools/.venv/
   __pycache__/
   *.pyc
   ```

2. **Exclude `.venv` in your JetBrains IDE:**

   * In the **Project** tool window, expand `py_tools`.
   * Right-click `.venv` → **Mark Directory as** → **Excluded**.
   * This prevents JetBrains from indexing all of `site-packages` (and prevents Go/IDE search clutter).
   * *(Note: If you use IntelliJ IDEA Ultimate with the Python plugin rather than GoLand, you can also go to **File → Project Structure → SDKs** and point a Python SDK to `py_tools\.venv\Scripts\python.exe`. In GoLand, simply marking `.venv` as Excluded is all you need.)*

---

### Step 6: Tell Kilo How to Use Python (`AGENTS.md`)

Kilo automatically reads an `AGENTS.md` file placed at the root of your project. Because your agent runs on local models (`qwen3.8-27b` and `gemma4:31b`), giving it the exact Windows executable path stops it from guessing `python` or `python3` and complaining that Python isn't configured.

Create a file named `AGENTS.md` in the **root of your Go project** with the following content:

```markdown
# Project Workspace Rules

## Python Environment
- This repository is primarily written in Go, with a dedicated Python workspace in `py_tools/`.
- Python 3.14.10 is configured in a virtual environment at `py_tools\.venv\`.
- **Always** execute Python and pip using the virtual environment's executable directly from the project root (do not call global `python`, `python3`, or `pip`):
  - Run Python: `.\py_tools\.venv\Scripts\python.exe <script.py>`
  - Run pip: `.\py_tools\.venv\Scripts\python.exe -m pip <command>`
- Place all Python scripts and dependencies (`requirements.txt`) inside `py_tools/`.
- Do not read, search, or modify files inside `py_tools\.venv\`.
```

*(Optional)* If you use Kilo's codebase indexing or want to ensure Kilo never scans `.venv` files, also add `py_tools/.venv/` to a `.kilocodeignore` file in the project root.

---

### Step 7: Test It with Kilo

1. Start a **new Kilo task/session** (so it picks up the new `AGENTS.md` file).

2. Ask Kilo:

   > *"Run a quick check using the project's Python environment to print the Python version and executable path."*

3. Kilo should run `.\py_tools\.venv\Scripts\python.exe` and confirm it is executing inside `...\your-go-project\py_tools\.venv\Scripts\python.exe` on Python `3.14.10`.