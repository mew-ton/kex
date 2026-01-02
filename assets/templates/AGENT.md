# Project Guidelines (Ref: Kex)

This project uses [Kex](https://github.com/mew-ton/kex) for knowledge management.
Before modifying code or documentation, please follow these guidelines.

## 1. Design Phase Guidelines
Before proposing changes or writing code, understand the existing design.

1.  **Search for design documents** using the `kex` MCP tool `search_documents`.
    *   **Keywords:** "architecture", "design", "requirements", "flow", "specification".
2.  **Read the relevant design docs** using `read_document`.
3.  **Validate your plan** against these design principles.

## 2. Implementation Phase Guidelines
Ensure you follow the established coding standards.

1.  **Search for coding guidelines** using the `kex` MCP tool `search_documents`.
    *   **Keywords:** "logging", "error", "structure", "naming", "testing".
2.  **Read the relevant coding docs** using `read_document`.
3.  **Adhere strictly** to the guidelines found.

## 3. Documentation Phase Guidelines
Maintain clarity and consistency in documentation.

1.  **Search for documentation guidelines** using the `kex` MCP tool `search_documents`.
    *   **Keywords:** "documentation", "style", "markdown", "language", "grammar".
2.  **Read the relevant style guides** using `read_document`.
3.  **Adhere strictly** to the formatting rules.

## 4. Adding New Knowledge
To add new knowledge to this project:
1.  Create a Markdown file in the `contents/` directory.
2.  Follow the directory structure examples in `contents/documentation/kex/`.
3.  Ensure the file has valid frontmatter (title, ID).
4.  Run `kex check` to validate your new documents.
