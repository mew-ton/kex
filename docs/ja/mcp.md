# MCP ツール

Kex は接続された MCP クライアントに対して以下のツールを提供します。

## `search_documents`

知識ベースから関連するドキュメントを検索します。

- **引数**:
  - `keywords` (string[]): 検索するキーワードのリスト。
- **戻り値**: ドキュメントの概要リスト (ID, Title, Description, Path)。

## `read_document`

特定のドキュメントの完全な内容を取得します。

- **引数**:
  - `id` (string): 読み込むドキュメントの ID。
- **戻り値**: ドキュメントの完全なマークダウンコンテンツ。

## クライアント設定

AI エディタで Kex を使用するには、MCP 設定を行う必要があります。

### Claude Desktop / VSCode (RoDev)

`claude_desktop_config.json` または MCP 設定に以下を追加してください:

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

**注意**: リポジトリ引数には絶対パス、または `~/` から始まるパスを使用してください。

### Claude Code

Kex を [Claude Code](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/overview) で使用するには、以下のコマンドを実行して直接登録します:

```bash
# Kex を Claude Code に追加

## 利用ガイドライン

Kex の利用方法は、目的によって異なります。

### 1. コーディング (既存リポジトリへの組み込み)

既存のリポジトリでコーディング支援として Kex を使用する場合、エージェントのシステムプロンプト (例: `.cursorrules` やそれに相当するもの) に以下の指示を追加することをお勧めします:

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

### 2. ドキュメント執筆 (管理)

Kex が提供するドキュメント自体を執筆・管理する場合、`kex init` によって生成される `AGENTS.md` がその目的に特化しています。AI エージェントにこのファイルを参照させることで、リポジトリ内の知識をどのように追加・修正すべきかを理解させることができます。
