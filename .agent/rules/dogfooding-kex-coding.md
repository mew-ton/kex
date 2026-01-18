---
trigger: always_on
---

# Coding Guidelines (Ref: Kex)

This project uses [Kex](https://github.com/mew-ton/kex) for knowledge management.

**PRIMARY DIRECTIVE (STRICT):**
Before writing any code, planning a refactor, or answering technical questions, you **MUST** first search the knowledge base using **Kex MCP**.

**The "Search Signal" Protocol:**
1.  **ACKNOWLEDGE**: When you receive a coding task, your *very first* reaction must be "I need to check the guidelines."
2.  **SEARCH**: Execute `search_documents` immediately.
    *   **Scope Strategy**: Set the scope based on the **Language** or **Framework** you are using.
        *   Examples: `["go"]`, `["typescript"]`, `["react"]`, `["frontend"]`
    *   **Keyword Strategy**: Combine keywords from these 3 dimensions:
        *   **Object**: Target entity (e.g., "function", "variable", "test")
        *   **Symptom**: Context/Problem (e.g., "large", "error", "complex", "refactor")
        *   **Concept**: Quality/Goal (e.g., "clean-architecture", "safety", "naming")
3.  **VERIFY**: You are forbidden from proceeding until you have read at least one relevant guideline or proven none exist.

**Forbidden Actions**:
*   Do NOT use `grep`, `cat`, or file system tools to read guidelines in `contents/`. Use Kex tools ONLY.

**MCP Failure Protocol:**
If `kex` MCP server is unavailable or `search_documents` fails:
1.  **STOP** immediately.
2.  **REPORT** the issue to the user ("Kex MCP is unavailable").
3.  **WAIT** for further instructions. Do **NOT** revert to manual file search (`grep`, `find`, etc.) without explicit permission.

**Self-Hosting Development (Kex on Kex)**
*   **Safety First**: Do NOT run `make build` or overwrite `./bin/kex` early. It breaks your tools.
*   **Test Logic**: Use `go test` to verify changes before building.

**Forbidden Actions**:
*   Do NOT use `grep`, `cat`, or file system tools to read guidelines in `examples/contents/`. Use Kex tools ONLY.