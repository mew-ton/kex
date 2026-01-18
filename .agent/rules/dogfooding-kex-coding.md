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
    *   **Scope Strategy**: Construct the scope hierarchy: `Domain` / `Platform` (Optional) / `Technology` (Optional).
        *   **Rule**: You **MUST** include the `Domain` (e.g., `coding`, `security`). You **SHOULD** include Platform/Technology if relevant to the specific task (e.g., adding an API -> include `api`).
        *   **Manual Construction (Best for Verification/Planning)**:
            *   `["coding", "api", "go"]` (Domain + Platform + Tech)
            *   `["security", "web"]` (Domain + Platform)
        *   **File Inference (Fast Path for Coding)**: Use `filePath` to automatically infer scopes when editing existing files.
            *   `filePath: "main.go"` implies `["coding", "go"]`.
    *   **Keyword Strategy**: Combine keywords from these 3 dimensions (Object / Symptom / Concept).

    > [!TIP]
    > **Scope Cheatsheet (Common Examples)**
    > *   **Domain**: `coding`, `security`, `architecture`, `vcs` (git)
    > *   **Platform**: `web`, `api`, `cli`, `backend`, `frontend`, `mobile`
    > *   **Technology**: `go`, `typescript`, `react`, `docker`, `sql`
    >
    > **Construct the scope combining**: `[Domain, Platform, Technology]`
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