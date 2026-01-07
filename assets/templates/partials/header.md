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
