---
trigger: always_on
---
# Coding Rules (Ref: Kex)

This project uses [Kex](https://github.com/mew-ton/kex) for knowledge management. Before modifying code or Kex-managed documentation, please follow these rules.

**Core Rule (STRICT):** You are FORBIDDEN from writing any code until you have explicitly searched for and read relevant guidelines. "Guessing" or "Reading the room" is considered a violation.

## Design & Implementation Phase

**Core Rule**: Before planning or coding, you must prove you have consulted the knowledge base.

1.  **Mandatory Search** using `search_documents`.
    *   **Keywords**: "architecture", "design", "coding", "style", "naming", "testing", "strategy".
    *   **STOP**: Do not proceed until you have successfully retrieved and read at least one relevant document or confirmed none exist.
2.  **Read relevant docs** using `read_document`.
3.  **Validate your plan** (or existing code) against principles found.

## Refactoring & Code Audit Phase

**Core Rule**: When asked to review code, find bugs, or check for guideline violations, you MUST NOT rely solely on your internal training.

1.  **Identify the context** (Language, Framework, Component, Function).
2.  **Search for specific rules** using `search_documents`.
    *   **Keywords**: "anti-pattern", "best practice", "forbidden", "required", "error handling", [Context Specific Keywords].
    *   **Tip**: When focusing on a specific layer or component (e.g., `domain`, `infrastructure`), use `search_documents` with `exactScopeMatch: true` and the component name as the keyword to see all rules for that scope.
3.  **Verify the code** against the *explicit* rules retrieved from Kex.
    *   *Do not assume* standard conventions apply if Kex has specific overrides.

## Self-Hosting Development (Kex on Kex)

**Core Rule**: When developing Kex itself, you MUST prioritize keeping the existing binary stable to continue using MCP tools.

1.  **Do NOT run `make build`/`make e2e` immediately**: Overwriting `./bin/kex` during development risks breaking your MCP connection.
2.  **Use `go test` for logic**: verify changes using `go test ./internal/...` or `go test ./e2e/...` (without build) first.
3.  **Build Last**: Only run `make build` or full `make e2e` after logic is verified and tests pass.

## General Usage Note

**Note**: Use `Glob`/`read_file_content` (or equivalent file system tools) only for existence checks, not for content search. Always rely on the indexed knowledge base via `search_documents`.
