<!-- AGENTS.md at project root -->
# Go v1.27 Development Standards

- **Context Grounding**: All generated Go code must strictly adhere to Golang Version 1.27 features, idioms, and standard library conventions.
- **Knowledge Retrieval (Active Search)**: If you are unsure about Go 1.27 specific syntax, or before writing/refactoring complex logic, you must actively query the local knowledge base using your file search or shell tools:
- For incremental feature updates or version differences, search the `F:/Go_Knowledge/releases/` directory.
- For core language mechanics, read the `F:/Go_Knowledge/language-spec.md/language-spec.md` file.
- **Strict Compliance**: Reject and avoid any deprecated Go functions or patterns outlined as obsolete in the local documentation.
- **Tooling First**: Do not guess version-specific behaviors. Rely on reading the specified knowledge paths to guarantee factual accuracy before outputting code.
