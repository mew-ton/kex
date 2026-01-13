---
trigger: always_on
---
# Documentation Rules (Ref: Kex)

This project uses [Kex](https://github.com/mew-ton/kex) for knowledge management.

**PRIMARY DIRECTIVE (STRICT):**
Before writing any documentation, planning a refactor, or answering technical questions, you **MUST** first search the knowledge base.

**The "Search Signal" Protocol:**
1.  **ACKNOWLEDGE**: When you receive a documentation task, your *very first* reaction must be "I need to check the rules."
2.  **SEARCH**: Execute `search_documents` immediately.
    *   **Keywords**: Use terms related to the task (e.g., "documentation", "style", "formatting", "naming").
    *   **Scope**: When working in a specific domain, you MUST use `exactScopeMatch: true` with the relevant scope (e.g., `["documentation", "kex"]`).
3.  **VERIFY**: You are forbidden from proceeding until you have read at least one relevant guideline or proven none exist.

**MCP Failure Protocol:**
If `kex` MCP server is unavailable or `search_documents` fails:
1.  **STOP** immediately.
2.  **REPORT** the issue to the user ("Kex MCP is unavailable").
3.  **WAIT** for further instructions. Do **NOT** revert to manual file search (`grep`, `find`, etc.) without explicit permission.

**Forbidden Actions**:
*   Do NOT use `grep`, `cat`, or file system tools to read guidelines in `examples/contents/*. Use Kex tools ONLY.
