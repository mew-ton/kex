---
trigger: always_on
---
# Documentation Rules (Ref: Kex)

This project uses [Kex](https://github.com/mew-ton/kex) for knowledge management.

**PRIMARY DIRECTIVE (STRICT):**
Before writing any documentation, planning a refactor, or answering technical questions, you **MUST** first search the knowledge base.

**The "Search Signal" Protocol:**
1.  **ACKNOWLEDGE**: When you receive a documentation task, your *very first* reaction must be "I need to check the style guidelines."
2.  **SEARCH**: Execute `search_documents` immediately.
    *   **Scope Strategy**: Construct the scope hierarchy: `Domain` / `Platform` (Optional) / `Technology` (Optional).
        *   **Rule**: You **MUST** include the `Domain` (e.g., `kex`, `api`).
        *   **Manual Construction (Primary)**:
            *   `["kex", "documentation"]` (Domain + Tech)
            *   `["api", "documentation"]` (Domain + Tech)
        *   **File Inference (Helper)**: `filePath` automatically infers scopes, but is secondary to understanding the correct domain.
            *   `filePath: "docs/api.md"` infers `["api", "documentation"]`.
    *   **Keyword Strategy**: Combine keywords from these 3 dimensions:
        *   **Object**: Target entity (e.g., "frontmatter", "filename")
        *   **Symptom**: Context/Problem (e.g., "missing", "invalid")
        *   **Concept**: Quality/Goal (e.g., "consistency", "style")

    > [!TIP]
    > **Scope Cheatsheet (Documentation Context)**
    > *   **Domain**: `kex`, `api`, `cli`, `sdk` (Target Product/Area)
    > *   **Platform**: `user-guide`, `reference`, `architecture`
    > *   **Technology**: `markdown`, `mermaid`, `godoc`
    >
    > **Construct the scope**: `[Domain, Platform, Technology]`
3.  **VERIFY**: You are forbidden from proceeding until you have read at least one relevant guideline or proven none exist.

**MCP Failure Protocol:**
If `kex` MCP server is unavailable or `search_documents` fails:
1.  **STOP** immediately.
2.  **REPORT** the issue to the user ("Kex MCP is unavailable").
3.  **WAIT** for further instructions. Do **NOT** revert to manual file search (`grep`, `find`, etc.) without explicit permission.

**Forbidden Actions**:
*   Do NOT use `grep`, `cat`, or file system tools to read guidelines in `examples/contents/*. Use Kex tools ONLY.
