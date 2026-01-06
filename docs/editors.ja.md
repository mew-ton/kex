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

## Claude Code

Kex を [Claude Code](https://docs.anthropic.com/ja/docs/agents-and-tools/claude-code/overview) で使用するには、以下のコマンドを実行して直接登録してください:

```bash
claude mcp add kex -- kex start
```

