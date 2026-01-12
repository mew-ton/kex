# 設定

## `.kex.yaml`

設定ファイル `kex.yaml` (または `.kex.yaml`) はプロジェクトのルートに配置します。

```yaml
root: contents
baseURL: https://example.com/docs/
```

### `root` (任意)

マークダウン形式のドキュメントファイルが含まれるディレクトリを指定します。

- **型**: `string`
- **デフォルト**: `contents`
- **説明**: Kex はこのディレクトリ内のすべての `.md` ファイルを再帰的にインデックスします。このパスは、`kex start` や `kex check` を実行するディレクトリ (通常はプロジェクトルート) からの相対パスです。

### `baseURL` (任意)

ドキュメントのリモートホスティング先のベース URL を定義します。

- **型**: `string`
- **例**: `https://my-org.github.io/my-repo/`
- **説明**:
    - `kex generate` が `kex.json` インデックス内の絶対 URL を生成するために使用します。
    - 省略した場合、`kex generate` は相対パスを出力します (これはほとんどの静的サイト構成で問題なく機能します)。
    - ドキュメントを異なるドメインから利用する場合など、相対パスが機能しないコンテキストで `kex.json` を使用する場合に必要となります。

### `remoteToken` (任意)

プライベートリポジトリ (例: GitHub Private Pages) 用の認証トークンです。

- **型**: `string`
- **説明**: 設定されている場合、リモートドキュメントの取得時に `Authorization` ヘッダーの Bearer トークンとして送信されます。

### `update` (任意)

`kex update` の動作を設定します。

```yaml
update:
  documents:
    kex: all
  ai-mcp-rules:
    targets: [antigravity, claude]
    scopes: [coding, documentation]
  ai-skills:
    targets: [claude]
    keywords: [go, typescript]
```

- **documents**:
    - `kex`: Kex システムドキュメントを管理します (`all` または `none`)。
- **ai-mcp-rules**:
    - `targets`: 静的ルールを生成するエージェントのリスト (例: `[antigravity, claude]`)。
    - `scopes`: 強制するルールのスコープのリスト (例: `[coding, documentation]`)。
- **ai-skills**:
    - `targets`: 動的スキルを生成するエージェントのリスト (例: `[claude]`)。
    - `keywords`: スキルをフィルタリングするためのキーワードリスト。

> **Note**: `kex init` は選択内容に基づいてこの設定を自動的に生成します。

### `logging` (任意)

サーバーのロギング設定です。

```yaml
logging:
  file: ".kex/kex.log"
  level: "info"
```

- **file**: ログファイルのパス。省略した場合、ログは標準エラー出力 (Stderr) に書き込まれます。
- **level**: ログレベル (デフォルト: `info`)。


## 環境変数 (Environment Variables)

Kex は以下の環境変数をサポートしています:

### `KEX_REMOTE_TOKEN`

- **目的**: プライベートなリモートドキュメントエンドポイントに対する認証。
- **優先順位**: `.kex.yaml` の `remoteToken` よりも優先されます。
- **使用方法**:
  ```bash
  export KEX_REMOTE_TOKEN="your-secret-token"
  kex start https://private.example.com/docs/
  ```


