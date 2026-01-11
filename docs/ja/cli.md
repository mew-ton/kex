# CLI リファレンス

## `kex init`

新しい Kex リポジトリを初期化します。

```bash
kex init [options]
```

- **対話モード**: フラグを指定せずに実行すると、AI エージェント (Antigravity, Cursor, Claude) とスコープ (Coding, Documentation) を対話形式で尋ねられます。
- **フラグ**:
    - `--agents=<agent>`: AI エージェントを選択します (複数指定可)。
    - `--scopes=<scope>`: スコープを選択します (複数指定可)。
- 指定された設定で `.kex.yaml` を作成します。
- `contents/` ディレクトリを作成します。
- 選択内容に基づいて、エージェントルールファイル (例: `.antigravity/rules/kex-coding.md`) を生成します。

## `kex check`

ドキュメントリポジトリの整合性を検証します。

```bash
kex check
```

以下の項目をチェックします:
- YAML Frontmatter の形式が正しいか
- ID とファイル名が一致しているか
- 必須フィールドが含まれているか
- ID に重複がないか

## `kex start`

MCP サーバーを起動します。

```bash
kex start [options] [path|url]
```

- **ローカルモード**: `kex start ./my-docs`
- **リモートモード**: `kex start https://example.com/docs/` (`kex.json` が配置されている必要があります)
- **フラグ**:
    - `--log-file=<path>`: ログを標準エラー出力ではなく、指定したファイルに書き込みます。
    - `--root=<path>`: ガイドラインのコンテンツルートを指定します (ローカルモードのみ有効)。



## `kex generate`

リモートホスティング (GitHub Pages など) 用に静的サイト構造を生成します。

```bash
kex generate [output-dir]
```

- すべての "adopted" (採用済み) ドキュメントを検証します。
- `dist/` ディレクトリを作成します。
- `kex.json` (インデックスファイル) を生成します。
- マークダウンファイルを `dist/` にコピーします。

## `kex update`

既存のリポジトリ内の Kex システムドキュメントとエージェント設定を更新します。

```bash
kex update [options]
```

- **システムドキュメント (`contents/documentation/kex/*`)**: Kex バイナリのバージョンに合わせて内容を更新します (上書きされます)。
- **エージェントルール**: `.kex.yaml` の戦略に基づいて更新します (上書き)。
- 動作のカスタマイズについては `.kex.yaml` の設定を参照してください。
