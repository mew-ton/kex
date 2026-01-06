# Antigravity Project Guidelines

You have access to the `kex` MCP server in this workspace.
This tool is designed to help you strictly adhere to project guidelines throughout the entire development lifecycle, from **Design** to **Implementation**.

**Core Rule:** Before taking action in any phase, you MUST consult the project's documentation using `kex`.

## 1. Design Phase Guidelines
Before proposing changes or writing code, you MUST understand the existing design and architecture.

1.  **Search for design documents** using the `kex` MCP tool `search_documents`.
    *   **Keywords:** "architecture", "design", "requirements", "flow", "specification".
2.  **Read the relevant design docs** using `read_document`.
3.  **Validate your plan** against these design principles.

## 2. Implementation Phase Guidelines
Before modifying or creating code, ensure you follow the established coding standards.

1.  **Search for coding guidelines** using the `kex` MCP tool `search_documents`.
    *   **Keywords:** "logging", "error", "structure", "naming", "testing", "function", "frontend", "typescript", "vue".
2.  **Read the relevant coding docs** using `read_document`.
3.  **Adhere strictly** to the guidelines found.


## 3. Documentation Phase Guidelines
When writing or updating documentation (README, guides, comments), maintain clarity and consistency.

1.  **Search for documentation guidelines** using the `kex` MCP tool `search_documents`.
    *   **Keywords:** "documentation", "style", "markdown", "language", "grammar".
2.  **Read the relevant style guides** using `read_document`.
3.  **Adhere strictly** to the writing style and formatting rules.

## 4. General Usage Note
**Note:** This workspace uses a local binary `./kex` for the MCP server to ensure stability.
If the `./kex` binary is missing or if the MCP server fails to start, you MUST run:
```bash
make build
```
before attempting to use MCP tools again.

Always ensure you are following the project's own documentation ("Dogfooding").
