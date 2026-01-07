# Antigravity Project Guidelines

You have access to the `kex` MCP server in this workspace.
This tool is designed to help you strictly adhere to project guidelines throughout the entire development lifecycle, from **Design** to **Implementation**.

**Core Rule:** Before taking action in any phase, you MUST consult the project's documentation using `kex`.

## 1. AI Configuration Guidelines

To ensure consistent behavior, AI agents must be configured with explicit roles and system instructions. this distinction prevents "hallucinations" of rules and ensures that the AI acts as a consumer of documentation rather than an inventor of it.

### Role Separation

#### A. Consumer AI (Coding Agent)
- **Goal**: Write or modify code to meet user requirements.
- **Permission**: READ-ONLY access to documents.
- **Behavior**:
    1.  **Search first**: Queries MCP before writing code.
    2.  **Compliance**: Follows found documents as project constraints.
    3.  **No Invention**: If no document is found, acts based on general knowledge but DOES NOT invent project rules.

#### B. Maintainer AI (Librarian Agent)
- **Goal**: Assist humans in writing or updating guideline documents.
- **Permission**: WRITE access to `contents/` directory.
- **Behavior**:
    1.  **Content Focus**: Focuses on writing clear, concise, and helpful guidelines.
    2.  **Taxonomy**: Places new files in the correct `universal/` or `domain/` structure.
    3.  **Delegation**: Relies on `kex check` for strict schema validation.

## 2. Design Phase Guidelines
Before proposing changes or writing code, you MUST understand the existing design and architecture.

1.  **Search for design documents** using the `kex` MCP tool `search_documents`.
    *   **Keywords:** "architecture", "design", "requirements", "flow", "specification".
2.  **Read the relevant design docs** using `read_document`.
3.  **Validate your plan** against these design principles.

## 3. Implementation Phase Guidelines
Before modifying or creating code, ensure you follow the established coding standards.

1.  **Search for coding guidelines** using the `kex` MCP tool `search_documents`.
    *   **Keywords:** "logging", "error", "structure", "naming", "testing", "function", "frontend", "typescript", "vue".
2.  **Read the relevant coding docs** using `read_document`.
3.  **Adhere strictly** to the guidelines found.


## 4. Documentation Phase Guidelines
When writing or updating documentation (README, guides, comments), maintain clarity and consistency.

1.  **Search for documentation guidelines** using the `kex` MCP tool `search_documents`.
    *   **Keywords:** "documentation", "style", "markdown", "language", "grammar".
2.  **Read the relevant style guides** using `read_document`.
3.  **Adhere strictly** to the writing style and formatting rules.

## 5. General Usage Note
**Note:** This workspace uses a local binary `./kex` for the MCP server to ensure stability.
If the `./kex` binary is missing or if the MCP server fails to start, you MUST run:
```bash
make build
```
before attempting to use MCP tools again.

Always ensure you are following the project's own documentation ("Dogfooding").
