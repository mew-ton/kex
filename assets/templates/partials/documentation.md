## Documentation Phase

**Core Rule**: When editing documents under `./contents` (or configured Kex root), maintain clarity and consistency.

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
