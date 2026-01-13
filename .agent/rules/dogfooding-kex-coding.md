---
trigger: always_on
---
# Coding Rules (Ref: Kex)

This project uses [Kex](https://github.com/mew-ton/kex) for knowledge management.

**PRIMARY DIRECTIVE (STRICT):**
Before writing any code, planning a refactor, or answering technical questions, you **MUST** first search the knowledge base.

**The "Search Signal" Protocol:**
1.  **ACKNOWLEDGE**: When you receive a coding task, your *very first* reaction must be "I need to check the rules."
2.  **SEARCH**: Execute `search_documents` immediately.
    *   **Keywords**: Use terms related to the task (e.g., "coding", "style", "testing", "naming", "architecture", "function", "component").
    *   **Scope**: When working in a specific language/domain (e.g., Go, TypeScript), you MUST use `exactScopeMatch: true` with the relevant scope (e.g., `["coding", "go"]`, `["coding", "react", "component"]`, `["coding", "utils", "function"]`).
3.  **VERIFY**: You are forbidden from proceeding until you have read at least one relevant guideline or proven none exist.

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
