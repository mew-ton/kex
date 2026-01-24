# AI Skills

Kex は、互換性のある AI エージェント (現在は Claude Desktop と Claude Code) 向けに **動的スキル (Dynamic Skills)** を生成できます。

## Skills とは?

Skills は、ドキュメントに基づいて Kex が生成する動的なルールファイルです。静的な MCP ルール (固定された指示) とは異なり、Skills は本質的に **検索拡張されたプロンプト** です。

- Kex はドキュメントディレクトリをスキャンします。
- 設定された `keywords` に一致するドキュメントを検索します。
- それらのドキュメントを要約した "Skill" ファイル (例: `.claude/skills/kex/go.md`) を生成します。

## 仕組み

`kex update` が実行されると:

1.  `.kex.yaml` の `ai-skills` 設定を確認します。
2.  `keywords` を使用してナレッジベースを検索します。
3.  関連するドキュメント (タイトル、説明、パス) を取得します。
4.  このリストを Skill テンプレート (例: `assets/templates/.claude/skills/kex/skill.md.template`) に注入します。
5.  結果をリポジトリに保存します (例: `.claude/skills/kex/my-skill.md`)。

## 設定

`.kex.yaml` の例:

```yaml
update:
  ai-skills:
    targets: [claude]        # スキルを生成する対象エージェント
    keywords: [go, testing]  # ナレッジを収集するキーワード
```

この設定は Kex に以下のように指示します: 「`go` または `testing` に関連するすべてのドキュメントを見つけ、それらのガイドラインが存在することを Claude に教えるための Skill を作成せよ」。

## メリット

- **コンテキスト認識**: エージェントはすべてのファイルを読み込むことなく、*どのような* ガイドラインが存在するかを把握できます。
- **トークン効率**: Skill ファイルには要約のみが含まれ、全文は含まれません。
- **発見可能性**: エージェントは必要に応じて `search_documents` MCP ツールを使用して、特定の詳細を読むことができます。

## ベストプラクティス

### バージョン管理 (.gitignore)

Skill ファイルは `kex update` によって自動的に生成されるため、**バージョン管理から除外する** ことを推奨します。これにより、リポジトリをクリーンに保ち、`kex update` を信頼できる唯一の情報源として運用できます。

`.gitignore` に以下を追加してください:

```gitignore
# Kex Generated Skills
.agent/skills/kex.*
.claude/skills/kex.*
```
