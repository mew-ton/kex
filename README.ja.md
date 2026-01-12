# Kex

[English](./README.md) | **日本語**

**Kex** は、キーワード検索を通じて **最小限かつ最適なコーディングガイドライン** を AI に提供する、**司書ツール (MCP / Skills 管理)** です。「Vibe Coding」を支援するために設計されており、AI エージェント (および人間) が必要な瞬間に適切なドキュメントへアクセスできるようサポートします。

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.25+-blue)

## 名前について

主機能である `keyword-indexed knowledge exchange` に由来しています

## クイックスタート

バイナリのダウンロードや `go install` の手順については、**[インストール](docs/ja/installation.md)** を参照してください。

### 1. 初期化
```bash
kex init
```

### 2. 起動 (ローカル)
```bash
kex start .
```

### 3. 起動 (リモート / GitHub Pages)
```bash
kex start https://my-org.github.io/guidelines/
```

## ドキュメント

詳細なドキュメントは `docs/ja/` ディレクトリに格納されています:

- **[機能解説](docs/ja/feature-index.md)**: キーワードとスコープの仕組みについて。
- **[コアコンセプト](docs/ja/concept.md)**: 設計哲学、アーキテクチャ、Draft/Adopted ステータスについて。
- **[CLI リファレンス](docs/ja/cli.md)**: `init`, `check`, `start`, `generate` コマンドの使い方。
- **[設定](docs/ja/configuration.md)**: `.kex.yaml` 設定ファイルのリファレンス。
- **[MCP ツール](docs/ja/feature-mcp.md)**: AI エージェントに提供されるツール群とクライアント設定。
- **[AI Skills](docs/ja/feature-skills.md)**: AI エージェント向けの動的ナレッジ生成。
- **[ドキュメントの執筆](docs/ja/documentation.md)**: Frontmatter スキーマとコンテンツのガイドライン。
- **[ベストプラクティス](docs/ja/best-practice.md)**: ナレッジベースの効果的な構造化方法。

## ユースケース

- **[モノレポ / ローカル](docs/ja/usecase-monorepo.md)**: コードとドキュメントを一元管理する場合。
- **[中央リポジトリ](docs/ja/usecase-central-repo.md)**: GitHub Pages 等を通じて、組織全体のガイドラインを一箇所で共有する場合。

## 貢献 (Contributing)

[CONTRIBUTING.ja.md](CONTRIBUTING.ja.md) を参照してください。

## ライセンス

MIT
