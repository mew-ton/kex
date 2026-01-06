# Editor Configuration

## Claude Desktop

To use Kex with [Claude Desktop](https://claude.ai/download), add the following configuration to your client's settings file (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "kex": {
      "command": "kex",
      "args": ["start"]
    }
  }
}
```

## Claude Code

To use Kex with [Claude Code](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/overview), run the following command to register it directly:

```bash
# Initialize with Claude Code support (generates CLAUDE.md)
kex init --agent-type=claude

# Add kex to Claude Code
claude mcp add kex -- kex start
```


