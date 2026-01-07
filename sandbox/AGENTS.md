# Antigravity Project Guidelines

You have access to the `kex` MCP server in this workspace.
This tool is designed to help you strictly adhere to project guidelines throughout the entire development lifecycle.

**Core Rule:** Before taking action in any phase, you MUST consult the project's documentation using `kex`.

## 1. AI Configuration Guidelines

To ensure consistent behavior, AI agents must be configured with explicit roles and system instructions.

### Role Separation

#### A. Consumer AI (Coding Agent)
- **Goal**: Write or modify code to meet user requirements.
- **Behavior**:
    1.  **Search first**: Queries MCP before writing code.
    2.  **Compliance**: Follows found documents.

#### B. Maintainer AI (Librarian Agent)
- **Goal**: Assist humans in writing or updating guideline documents.
- **Behavior**:
    1.  **Content Focus**: Focuses on writing clear, concise guidelines.
    2.  **Taxonomy**: Places new files in correct structure.


<!-- kex: auto-update start -->
## Documentation Phase

**Core Rule**: Maintain clarity and consistency in documentation.

1.  **Search for style guides** using `search_documents`.
    *   **Keywords**: "documentation", "style", "markdown", "format".
2.  **Read relevant guides** using `read_document`.
3.  **Adhere strictly** to formatting rules.

## Adding New Knowledge

1.  **Search existing structure** to understand where new files belong.
2.  **Check for conflicts** using `Glob`.
3.  **Create Markdown files** with valid Frontmatter (id, title, description, keywords).
4.  **Run `kex check`** to validate.


## Design & Implementation Phase

**Core Rule**: Before proposing changes or writing code, understand existing design and coding standards.

1.  **Search for guidelines** using `search_documents`.
    *   **Keywords**: "architecture", "design", "coding", "style", "naming", "testing".
2.  **Read relevant docs** using `read_document`.
3.  **Validate your plan** against principles found.

**Note**: Use `Glob`/`Read` only for existence checks, not for content search.


<!-- kex: auto-update end -->

## General Usage Note
**Note:** This workspace uses a local binary `./kex` (or global `kex`) for the MCP server.
If the MCP server fails to start, ensure `kex` is built and in your PATH or workspace root.

Always ensure you are following the project's own documentation ("Dogfooding").
