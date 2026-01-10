---
paths: "{{.Root}}/**/*.md"
---
# Documentation Rules (Ref: Kex)

This project uses [Kex](https://github.com/mew-ton/kex) for knowledge management. Before modifying code or Kex-managed documentation, please follow these rules.

**Core Rule:** Before taking action in any phase, you MUST consult the project's documentation using `kex`.

## Documentation Phase

**Core Rule**: When editing documents under `{{.Root}}`, maintain clarity and consistency.

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

## General Usage Note

**Note**: Use `Glob`/`read_file_content` (or equivalent file system tools) only for existence checks, not for content search. Always rely on the indexed knowledge base via `search_documents`.

