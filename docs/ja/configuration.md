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


