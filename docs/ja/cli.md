# CLI リファレンス

## `kex init`

新しい Kex リポジトリを初期化します。

```bash
kex init [options]
```

- `.kex.yaml` を作成します
- `contents/` ディレクトリを作成します
- `AGENTS.md` (または `--agent-type=claude` で `CLAUDE.md`) を生成します

## `kex check`

ドキュメントリポジトリの整合性を検証します。

```bash
kex check
```

以下の項目をチェックします:
- 無効な YAML Frontmatter
- ID とファイル名の不一致
- 必須フィールドの欠落
- ID の重複

## `kex start`

MCP サーバーを起動します。

```bash
kex start [options] [path|url]
```

- **ローカルモード**: `kex start ./my-docs`
- **リモートモード**: `kex start https://example.com/docs/` (`kex.json` が必要です)



## `kex generate`

リモートホスティング (GitHub Pages) 用に静的サイト構造を生成します。

```bash
kex generate [output-dir]
```

- すべての "Adopted" ドキュメントを検証します。
- `dist/` ディレクトリを作成します。
- `kex.json` (インデックス) を生成します。
- マークダウンファイルを `dist/` にコピーします。
