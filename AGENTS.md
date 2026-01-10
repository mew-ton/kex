<!-- kex: auto-update start -->
# Project Guidelines (Ref: Kex)

This project uses [Kex](https://github.com/mew-ton/kex) for knowledge management. Before modifying code or Kex-managed documentation, please follow these guidelines.

**Core Rule:** Before taking action in any phase, you MUST consult the project's documentation using `kex`.


## Design & Implementation Phase

**Core Rule**: Before proposing changes or writing code, understand existing design and coding standards.

1.  **Search for guidelines** using `search_documents`.
    *   **Keywords**: "architecture", "design", "coding", "style", "naming", "testing", "strategy".
    *   **Critical:** If `kex` tools are unavailable or fail, **STOP** and report this to the user.
2.  **Read relevant docs** using `read_document`.
3.  **Validate your plan** (or existing code) against principles found.

## Refactoring & Code Audit Phase

**Core Rule**: When asked to review code, find bugs, or check for guideline violations, you MUST NOT rely solely on your internal training.

1.  **Identify the context** (Language, Framework, Component, Function).
2.  **Search for specific rules** using `search_documents`.
    *   **Keywords**: "anti-pattern", "best practice", "forbidden", "required", "error handling", [Context Specific Keywords].
3.  **Verify the code** against the *explicit* rules retrieved from Kex.
    *   *Do not assume* standard conventions apply if Kex has specific overrides.



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


## General Usage Note

**Note**: Use `Glob`/`read_file_content` (or equivalent file system tools) only for existence checks, not for content search. Always rely on the indexed knowledge base via `search_documents`.

<!-- kex: auto-update end -->

Always ensure you are following the project's own documentation ("Dogfooding").