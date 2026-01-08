# ユースケース: モノレポ / ローカル

最もシンプルな構成です。ドキュメントは、それが適用されるコードと同じ場所に存在します。

## セットアップ

1. プロジェクトルートで `kex init` を実行します。
   - これにより `.kex.yaml` と `AGENTS.md` が作成されます。
2. `contents/` にドキュメントを作成します。
3. エディタ設定で `kex start .` (または単に `kex start`) を実行するように設定します。

## ワークフロー

### 1. エージェント指示のセットアップ (`AGENTS.md`)

生成された `AGENTS.md` には、主に2つの目的があります:
1.  **ドキュメントメンテナンス**: `contents/` 内のファイルを追加・編集する方法についての指示。
2.  **コーディングガイドライン**: エージェント設定 (例: `.cursorrules`) に追加すべきプロンプト。これにより、エージェントはコーディング作業中に適切にガイドラインを検索できるようになります。

**アクション**: `AGENTS.md` の関連するセクションを、エージェントのシステムプロンプトまたは設定ファイルにコピーしてください。

### 2. 日々のワークフロー

1.  **コードを書く**: AI エージェントは、設定されたプロンプトに従って `kex` にガイドラインを問い合わせます。
2.  **ドキュメント更新**: ルールを変更・追加する場合、それに関連するコード変更と同じ Pull Request でマークダウンファイルを更新します。
3.  **CI/CD**: CI パイプラインで `kex check` を実行し、ドキュメントの整合性を担保します。

### 3. CI チェック例

`.github/workflows/check.yml` を作成します:

```yaml
name: Check Guidelines

on:
  pull_request:

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install Kex
        run: go install github.com/mew-ton/kex/cmd/kex@latest

      - name: Check Docs
        run: kex check
```

## メリット/デメリット

- **メリット**: コードとセットでバージョン管理され、変更のラグが発生しません。セットアップも非常に簡単です。
- **デメリット**: 複数のリポジトリ間でルールを共有・同期するのは難しくなります。
