# ユースケース: 中央リポジトリ (GitHub Pages)

多数のサービス間でガイドラインを共有する組織に最適な構成です。

## セットアップ

1. 専用のリポジトリを作成します (例: `my-org/guidelines`)。
2. `kex init` を実行します。
3. GitHub Actions を設定し、`kex generate dist` を実行して `dist/` ディレクトリを GitHub Pages にデプロイします。

### CI チェック例

Pull Request に対して常に `kex check` を実行し、無効なドキュメントがマージされないようにします。

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

### デプロイワークフロー例

`.github/workflows/deploy.yml` を作成します:

```yaml
name: Deploy Guidelines

on:
  push:
    branches: ["main"]

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
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

      - name: Generate Site
        run: kex generate dist

      - name: Setup Pages
        uses: actions/configure-pages@v5

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: './dist'

      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

## 他のリポジトリでの利用

利用側のリポジトリで、リモート URL を指すように MCP クライアントを設定します:

```bash
kex start https://my-org.github.io/guidelines/
```

## ワークフロー

1. **中央での更新**: ガイドラインリポジトリでルールを更新します。
2. **デプロイ**: GitHub Actions が新しい静的サイトと `kex.json` を公開します。
3. **利用**: 接続されているすべてのエージェントは、即座に新しいルールを参照できるようになります。

## メリット/デメリット

- **メリット**: 唯一の信頼できる情報源 (SSOT) となり、組織全体へのグローバルな更新が容易です。
- **デメリット**: コードとルールベースのバージョン乖離 (コードがルールより遅れる可能性がある) が発生する可能性があります。また、ネットワーク接続に依存します。
