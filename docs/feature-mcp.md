# MCP Tools

Kex exposes the following tools to connected MCP Clients.

## `search_documents`

Searches the knowledge base for relevant documents.

- **Arguments**:
  - `keywords` (string[]): List of keywords to search for.
- **Returns**: A list of document summaries (ID, Title, Description, Path).

## `read_document`

Retrieves the full content of a specific document.

- **Arguments**:
  - `id` (string): The ID of the document to read.
- **Returns**: The full markdown content of the document.

## Client Configuration

To use Kex with your AI editor, you need to configure the MCP settings.

### Claude Desktop / VSCode (RoDev)

Add this to your `claude_desktop_config.json` or MCP settings:

```json
{
  "mcpServers": {
    "kex": {
      "command": "kex",
      "args": ["start", "/absolute/path/to/your/repo"]
    }
  }
}
```

**Note**: Use absolute paths or `~/` for the repository argument.

### Claude Code

To use Kex with [Claude Code](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/overview), run the following command to register it directly:

```bash
# Add kex to Claude Code
claude mcp add kex -- kex start /absolute/path/to/your/repo
```

**Note**: Use absolute paths or `~/` for the repository argument.

## Usage Guidelines

How you use Kex depends on your goal.

### 1. For Coding (Integration)

If you are using Kex to assist with coding in an existing repository, you should add the following instructions to your agent's system prompt (e.g., `.cursorrules` or `demos/AGENTS.md` equivalent):

> **Design Phase:**
> 1. **Search for design documents** using the `kex` MCP tool `search_documents`.
>    * **Keywords:** "architecture", "design", "requirements", "flow", "specification".
> 2. **Read the relevant design docs** using `read_document`.
> 3. **Validate your plan** against these design principles.
>
> **Implementation Phase:**
> 1. **Search for coding guidelines** using the `kex` MCP tool `search_documents`.
>    * **Keywords:** "logging", "error", "structure", "naming", "testing", "component".
> 2. **Read the relevant coding docs** using `read_document`.
> 3. **Adhere strictly** to the guidelines found.

### 2. For Documentation (Authoring)

If you are writing documentation that Kex will serve, `kex init` generates an `AGENTS.md` file designed for this purpose. You can simply point your AI agent to this file to help it understand how to add or modify knowledge in the repository.
