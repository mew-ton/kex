# CLI リファレンス

## `kex init`

新しい Kex リポジトリを初期化します。

```bash
kex init [options]
```

- **対話モード**: フラグが指定されていない場合、エージェントタイプとスコープ (Coding, Documentation) を尋ねます。
- **フラグ**: `--agent-type=<type>` (例: `general`, `claude`) でプロンプトをスキップできます。
- 選択された設定で `.kex.yaml` を作成します。
- `contents/` ディレクトリを作成します。
- 選択されたスコープに基づいて動的な内容で `AGENTS.md` (または `CLAUDE.md`) を生成します。

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

## `kex update`

既存のリポジトリ内の Kex システムドキュメントとエージェント設定を更新します。

```bash
kex update [options]
```

- **システムドキュメント (`contents/documentation/kex/*`)**: 現在のバイナリバージョンに合わせて更新します (上書き)。
- **`AGENTS.md`**: `<!-- kex: auto-update start -->` マーカー間のガイドラインを安全に更新します。
- 動作をカスタマイズするには `.kex.yaml` 設定を参照してください。
