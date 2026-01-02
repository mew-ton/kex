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
