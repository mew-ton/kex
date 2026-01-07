# 設定

## `.kex.yaml`

設定ファイル `kex.yaml` (または `.kex.yaml`) はプロジェクトのルートに配置します。

```yaml
root: contents
baseURL: https://example.com/docs/
```

### `root` (Optional)

マークダウン形式のドキュメントファイルが含まれるディレクトリを指定します。

- **型**: `string`
- **デフォルト**: `contents`
- **説明**: Kex はこのディレクトリ内のすべての `.md` ファイルを再帰的にインデックスします。このパスは、`kex start` や `kex check` を実行するディレクトリ (通常はプロジェクトルート) からの相対パスです。

### `baseURL` (Optional)

ドキュメントのリモートホスティング先のベース URL を定義します。

- **型**: `string`
- **例**: `https://my-org.github.io/my-repo/`
- **説明**:
    - `kex generate` が `kex.json` インデックス内の絶対 URL を生成するために使用します。
    - 省略した場合、`kex generate` は相対パスを出力します (これはほとんどの静的サイト構成で問題なく機能します)。
    - 相対パスが機能しない異なるドメインやコンテキストから `kex.json` を利用する場合に必要となります。

### `remoteToken` (Optional)

プライベートリポジトリ (例: GitHub Private Pages) 用の認証トークンです。

- **型**: `string`
- **説明**: 設定されている場合、リモートドキュメントの取得時に `Authorization` ヘッダーの Bearer トークンとして送信されます。

### `update` (Optional)

`kex update` の動作を設定します。

```yaml
update:
  strategies:
    "contents/documentation/kex/*": "overwrite"
    "AGENTS.md": "marker-update"
```

- **strategies**: Glob パターンと更新戦略のマップです。
  - `overwrite`: ファイルをテンプレートで置換します。
  - `marker-update`: マーカー間のコンテンツを更新します (`AGENTS.md` 用)。
  - `append`: コンテンツが欠落している場合、末尾に追加します。
  - `skip`: 何もしません。

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


