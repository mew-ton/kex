# 設定

## `.kex.yaml`

設定ファイル `kex.yaml` (または `.kex.yaml`) はプロジェクトのルートに配置します。

```yaml
root: contents
baseURL: https://example.com/docs/
```

### `source` (任意)

マークダウン形式のドキュメントファイルが含まれるディレクトリを指定します。

- **型**: `string`
- **デフォルト**: `contents`
- **説明**: Kex はこのディレクトリ内のすべての `.md` ファイルを再帰的にインデックスします。このパスは、`kex start` や `kex check` を実行するディレクトリ (通常はプロジェクトルート) からの相対パスです。

**設定例**:
```yaml
source: contents
```

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
  strategies:
    "contents/documentation/kex/*": "overwrite"
    "AGENTS.md": "marker-update"
```

- **strategies**: Glob パターンと更新戦略のマップです。
  - `overwrite`: ファイルをテンプレートで完全に置換します。
  - `marker-update`: マーカー間のコンテンツのみを更新します (`AGENTS.md` のようなファイル向け)。
  - `append`: コンテンツが欠落している場合、末尾に追加します 既存の内容は変更しません。
  - `skip`: 何もしません。

### `agent` (任意)

AI エージェント用のガイドライン生成を設定します。

```yaml
agent:
  type: general # または "claude"
  scopes:
    - coding
    - documentation
```

- **type**: 生成するエージェントガイドラインのタイプ。
    - `general`: 標準的な `AGENTS.md`。
    - `claude`: Anthropic Claude 向けの `CLAUDE.md`。
- **scopes**: ガイドラインに含めるセクション (スコープ) のリスト。
    - `coding`: 設計および実装フェーズのガイドライン。
    - `documentation`: ドキュメント作成フェーズのルール。

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


