---
trigger: always_on
---
# Documentation Rules (Ref: Kex)

This project uses [Kex](https://github.com/mew-ton/kex) for knowledge management. Before modifying code or Kex-managed documentation, please follow these rules.

**Core Rule (STRICT):** You are FORBIDDEN from writing any content until you have explicitly searched for and read relevant guidelines. "Guessing" or "Reading the room" is considered a violation.

## Documentation Phase

**Core Rule**: When editing documents under `./examples/contents`, maintain clarity and consistency.

1.  **Search for style guides** using `search_documents`.
    *   **Keywords**: "documentation", "style", "markdown", "format".
    *   **Critical:** If `kex` tools are unavailable or fail, **STOP** and report this to the user.
2.  **Read relevant guides** using `read_document`.
3.  **Adhere strictly** to formatting rules.

## Adding New Knowledge

1.  **Search existing structure** to understand where new files belong.
2.  **Check for conflicts** using `Glob`.
3.  **Create Markdown files** with valid Frontmatter (id, title, description, keywords).
4.  **Run `kex check`** to validate.

## Self-Hosting Development (Kex on Kex)

**Core Rule**: When developing Kex itself, you MUST prioritize keeping the existing binary stable to continue using MCP tools.

1.  **Do NOT run `make build`/`make e2e` immediately**: Overwriting `./bin/kex` during development risks breaking your MCP connection (and thus `kex check`).
2.  **Use `go test` for logic**: verify changes using `go test ./internal/...` or `go test ./e2e/...` (without build) first.
3.  **Build Last**: Only run `make build` or full `make e2e` after logic is verified and tests pass.

## General Usage Note

**Note**: Use `Glob`/`read_file_content` (or equivalent file system tools) only for existence checks, not for content search. Always rely on the indexed knowledge base via `search_documents`.
