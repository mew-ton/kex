# エディタ設定

## Claude Desktop

Kex を [Claude Desktop](https://claude.ai/download) で使用するには、クライアントの設定ファイル (`claude_desktop_config.json`) に以下の設定を追加してください:

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
