# ユースケース: モノレポ / ローカル

最もシンプルな構成です。ドキュメントは、それが規定するコードと一緒に存在します。

## セットアップ

1. プロジェクトルートで `kex init` を実行します。
   - これにより `.kex.yaml` と `AGENTS.md` が作成されます。
2. `contents/` にドキュメントを書きます。
3. エディタ設定で `kex start .` (または単に `kex start`) を実行するように設定します。

## ワークフロー

### 1. エージェント指示のセットアップ (`AGENTS.md`)

生成された `AGENTS.md` は2つの目的で使用します:
1.  **ドキュメントメンテナンス**: `contents/` 内のファイルを追加・編集する方法の指示。
2.  **コーディングガイドライン**: エージェント設定 (例: `.cursorrules`) に追加するプロンプト。これにより、エージェントはコーディング中にガイドラインを検索する方法を知ることができます。

**アクション**: `AGENTS.md` の関連するセクションを、エージェントのシステムプロンプトまたは設定にコピーしてください。

### 2. 日々のワークフロー

1.  **コードを書く**: AI エージェントは、追加されたプロンプトに基づいて `kex` にガイドラインを問い合わせます。
2.  **ドキュメント更新**: ルールが変更された場合、同じ PR 内でマークダウンファイルを更新します。
3.  **CI/CD**: CI パイプラインで `kex check` を実行し、ドキュメントの整合性を保証します。

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

- **メリット**: コードとバージョン管理される、遅延ゼロ、セットアップが簡単。
- **デメリット**: 複数のリポジトリ間でルールを共有するのが難しい。
