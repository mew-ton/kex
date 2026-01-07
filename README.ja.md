# Kex

[English](./README.md) | **日本語**

**Kex** は、**ドキュメントライブラリアン**として実装された**キーワードインデックス型ナレッジエクスチェンジ** (keyword-indexed knowledge exchange) です。"Vibe Coding" を支援するために設計されており、ローカル **MCP サーバー** として動作し、AI エージェント（および人間）が適切なタイミングで適切なドキュメントにアクセスできるように支援します。

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.25+-blue)

## クイックスタート

バイナリのダウンロードや `go install` の手順については **[インストール](docs/ja/installation.md)** を参照してください。

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

完全なドキュメントは `docs/ja/` ディレクトリにあります:

- **[機能解説](docs/ja/feature.md)**: キーワードとスコープについて。
- **[コアコンセプト](docs/ja/concept.md)**: 哲学、アーキテクチャ、Draft/Adopted ステータスについて。
- **[CLI リファレンス](docs/ja/cli.md)**: `init`, `check`, `start`, `generate` の使い方。
- **[設定](docs/ja/configuration.md)**: `.kex.yaml` のリファレンス。
- **[MCP ツール](docs/ja/mcp.md)**: AI エージェントに提供されるツールとクライアント設定。
- **[ドキュメントの執筆](docs/ja/documentation.md)**: Frontmatter スキーマとコンテンツのガイドライン。
- **[ベストプラクティス](docs/ja/best-practice.md)**: ナレッジベースの構造化方法。

## ユースケース

- **[モノレポ / ローカル](docs/ja/usecase-monorepo.md)**: コードと一緒にドキュメントを管理する場合。
- **[中央リポジトリ](docs/ja/usecase-central-repo.md)**: GitHub Pages を介して組織全体のガイドラインを共有する場合。

## 貢献 (Contributing)

[CONTRIBUTING.ja.md](CONTRIBUTING.ja.md) を参照してください。

## ライセンス

MIT
