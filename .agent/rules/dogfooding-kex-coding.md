---
trigger: always_on
---
# Coding Rules (Ref: Kex)

This project uses [Kex](https://github.com/mew-ton/kex) for knowledge management. Before modifying code or Kex-managed documentation, please follow these rules.

**Core Rule:** Before taking action in any phase, you MUST consult the project's documentation using `kex`.

## Design & Implementation Phase

**Core Rule**: Before proposing changes or writing code, understand existing design and coding standards.

1.  **Search for rules** using `search_documents`.
    *   **Keywords**: "architecture", "design", "coding", "style", "naming", "testing", "strategy".
    *   **Critical:** If `kex` tools are unavailable or fail, **STOP** and report this to the user.
2.  **Read relevant docs** using `read_document`.
3.  **Validate your plan** (or existing code) against principles found.

## Refactoring & Code Audit Phase

**Core Rule**: When asked to review code, find bugs, or check for guideline violations, you MUST NOT rely solely on your internal training.

1.  **Identify the context** (Language, Framework, Component, Function).
2.  **Search for specific rules** using `search_documents`.
    *   **Keywords**: "anti-pattern", "best practice", "forbidden", "required", "error handling", [Context Specific Keywords].
    *   **Requirement**: When focusing on a specific layer or component (e.g., `go`, `typescript`, `frontend`), you **MUST** use `search_documents` with `exactScopeMatch: true` and the scope names as keywords. This ensures you see *all* rules for that scope, protecting against "unknown unknowns".
3.  **Verify the code** against the *explicit* rules retrieved from Kex.
    *   *Do not assume* standard conventions apply if Kex has specific overrides.

## Self-Hosting Development (Kex on Kex)

**Core Rule**: When developing Kex itself, you MUST prioritize keeping the existing binary stable to continue using MCP tools.

1.  **Do NOT run `make build`/`make e2e` immediately**: Overwriting `./bin/kex` during development risks breaking your MCP connection.
2.  **Use `go test` for logic**: verify changes using `go test ./internal/...` or `go test ./e2e/...` (without build) first.
3.  **Build Last**: Only run `make build` or full `make e2e` after logic is verified and tests pass.

## Critical Anti-Patterns (Forbidden Actions)

**Restriction**: You MUST NOT use file system tools (`list_dir`, `view_file`, `read_file_content`, `grep`) to read, search, or discover guidelines in `examples/contents` (or any configured Kex root).

- **Why?**: Accessing files directly bypasses the "Dogfooding" process. We must test the *retrieval accuracy* of Kex itself. If you read the files directly, you are not testing the product.
- **Correct Action**: Use `search_documents` and `read_document`.
- **Exception**: You may read these files ONLY if your task is specifically to *edit* the text of the guideline itself (after retrieving it via Kex) or to debug the *parser logic*.
