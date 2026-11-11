---
name: go-project-directory-structure
description: Use this skill ONLY when scaffolding, creating, or setting up a new Golang project directory structure.
---

# Golang Project Directory Structure Guidelines

## Agent Directive
> **MANDATORY**: When creating a new Golang Project, all directories for source code and test files must be structured in strict accordance with the rules defined in this document. Never bypass these rules.
>
> **STOP AND ASK**: If the User fails to provide the (1) **Project Parent Directory**, (2) **Project Base Directory**, or (3) **Source Code Directory** names, you MUST stop and ask the User for these names before generating any folders or files.

---

## 1. Structure Overview

The Golang Project Directory Structure consists of seven (7) component directories. Three (3) are variable names supplied by the user, and four (4) are constant, automatically generated directories.

### Abstract Layout
```text
[Project Parent Directory] (Variable - User Supplied)
|
|_ [Project Base Directory] (Variable - User Supplied)
  |
  |_ AppTest           (Constant)
  |_ assets            (Constant - MUST be lowercase)
  |_ Docs              (Constant)
  |_ [source code dir] (Variable - User Supplied - MUST be lowercase)
  |_ TechNotes         (Constant)
```
## 2. Variable Directory Names (User Supplied)

You must obtain these three names from the user before scaffolding:

1. 
2. **Project Parent Directory**: The existing parent path where the new project will reside.
3. **Project Base Directory**: The root folder for the specific new Golang project.
4. **Source Code Directory**: The directory that will house the actual .go source and test files.*CRITICAL RULE*: The Source Code Directory name **MUST BE ALL LOWERCASE**.

------



## 3. Constant Directory Names (Auto-Generated)

These four directories must be created inside the **Project Base Directory** exactly as spelled below:

1. 
2. **AppTest**: Stores the main.go file. (Capital 'A' and 'T').
3. **assets**: Stores images and diagrams. (*CRITICAL RULE*: **MUST BE ALL LOWERCASE**).
4. **Docs**: Stores project documentation. (Capital 'D').
5. **TechNotes**: Stores technical notes and documents. (Capital 'T' and 'N').

*Note on Test Files: Test files do not get their own directory. They must always be stored in the same directory as the source code files they are testing.*

------



## 4. Concrete Scaffold Example

If the user provides the following variables:

- Project Parent Directory = F:\GoProjects\Arenero73\X01-SandBox\golang_code\LCD
- Project Base Directory = LCD_01
- Source Code Directory = lcd

You must generate the following exact file system tree:

```text
F:\GoProjects\Arenero73\X01-SandBox\golang_code\LCD\
|
|_ LCD_01\
  |
  |_ AppTest\
  |
  |_ assets\
  |
  |_ Docs\
  |
  |_ lcd\
  |
  |_ TechNotes\
```
