<!-- kex: auto-update start -->
## 1. Design Phase Guidelines

Before proposing changes or writing code, understand the existing design. **Always use MCP tools for content search.**

1. **Search for design documents** using the `kex` MCP tool `search_documents`.
   * **Keywords:** "architecture", "design", "requirements", "flow", "specification".
2. **Read the relevant design docs** using `read_document`.
3. **Validate your plan** against these design principles.

**Note:** Use `Glob`/`Read` only for file existence checks or duplication verification, not for content search.

## 2. Implementation Phase Guidelines

Ensure you follow the established coding standards. **Always use MCP tools for content search.**

1. **Search for coding guidelines** using the `kex` MCP tool `search_documents`.
   * **Keywords:** "logging", "error", "structure", "naming", "testing", "component", "frontend", "backend".
2. **Read the relevant coding docs** using `read_document`.
3. **Adhere strictly** to the guidelines found.

**Note:** Use `Glob`/`Read` only for file existence checks or duplication verification, not for content search.

## 3. Documentation Phase Guidelines

Maintain clarity and consistency in documentation. **Always use MCP tools for content search.**

1. **Search for documentation guidelines** using the `kex` MCP tool `search_documents`.
   * **Keywords:** "documentation", "style", "markdown", "language", "grammar", "frontmatter", "format".
2. **Read the relevant style guides** using `read_document`.
3. **Adhere strictly** to the formatting rules found.

**Note:** Use `Glob`/`Read` only for file existence checks or duplication verification, not for content search.

## 4. Adding New Knowledge

To add new knowledge to this project:

1. **Search for existing guidelines** using the `kex` MCP tool `search_documents` to understand format and structure.
   * **Keywords:** Related to the topic you're documenting, plus "format", "structure", "template".
2. **Read relevant examples** using `read_document` to understand the expected format.
3. **Check for file conflicts** using `Glob` to ensure no duplicate filenames exist.
4. Create a Markdown file in the appropriate `contents/` subdirectory.
5. Follow the directory structure and format patterns from existing documents.
6. Ensure the file has valid frontmatter (id, title, description, keywords).
7. Run `kex check` to validate your new documents.
<!-- kex: auto-update end -->

Always ensure you are following the project's own documentation ("Dogfooding").