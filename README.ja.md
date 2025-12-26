# Kex (Knowledge Exchange)

[English](./README.md) | **日本語**

**Kex** は、コーディングガイドラインやプロジェクトドキュメントを管理・提供するために設計された軽量なドキュメントライブラリアンであり、MCP (Model Context Protocol) サーバーです。AI エージェント（および人間）が適切なタイミングで適切なドキュメントにアクセスできるように支援します。

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.25+-blue)

## 機能

-   **MCP サーバー**: Model Context Protocol (JSON-RPC) を介してドキュメントを公開し、AI エージェントがドキュメントを「読み」「検索」できるようにします。
-   **構造の強制**: Frontmatter (ID, タイトル, ステータス, キーワード) を使用して、ドキュメントの厳密なスキーマを強制します。
-   **ライフサイクル管理**: `draft` (下書き) と `adopted` (採用済み) のドキュメントを区別します。
    -   **Draft**: 進行中の作業。エラーは警告として扱われます。
    -   **Adopted**: 正当な情報源 (Source of truth)。エラーがある場合、起動を阻止します。
-   **CLI ツール**:
    -   `init`: 新しいナレッジベースの雛形を作成します。
    -   `check`: すべてのドキュメントを検証し、整合性の問題を報告します。
    -   `start`: Stdio 経由で MCP サーバーを起動します。

## インストール

```bash
# リポジトリをクローン
git clone https://github.com/mew-ton/kex.git
cd kex

# ビルドとインストール
go install ./cmd/kex
```

## はじめに

### 1. リポジトリの初期化

プロジェクトのルートディレクトリに移動し、以下を実行します:

```bash
kex init
```

これにより、`contents/` ディレクトリと `.kex.yaml` 設定ファイルが作成されます。

### 2. ドキュメントの作成

`contents/` ディレクトリ内にマークダウンファイルを作成します。各ファイルには有効な YAML Frontmatter が**必須**です。
また、**AI ファースト** なドキュメント管理を実現するため、ドキュメントの本文は**英語**で記述することを強く推奨します。（多言語対応の AI であっても、英語のほうがコンテキストの理解精度が高いためです）

```markdown
---
id: my-guideline
title: Coding Standards
description: Guidelines for Go development
status: adopted
keywords: [go, style, lint]
---

# Coding Standards

Write your content here...
```

**ステータスの種類:**
-   `draft`: デフォルトではサーバーによって無視されます。進行中の作業に便利です。
-   `adopted`: アクティブでインデックスされます。すべてのチェックに合格する必要があります。

### 3. ドキュメントの検証

`check` コマンドを実行してドキュメントを検証します:

```bash
kex check
```

以下が報告されます:
-   Frontmatter の欠落
-   ファイル名と ID の不一致
-   パースエラー

### 4. サーバーの起動

MCP サーバーを起動して、AI との接続を許可します:

```bash
kex start
```

*注意: これは Stdio 上でインタラクティブな JSON-RPC セッションを開始します。これは (Claude Desktop のような) MCP クライアントによって実行されることを意図しています。*

## 提供される MCP ツール

Kex は接続された AI エージェントに対して以下のツールを公開します:

-   `search_documents(keywords: string[])`: 特定のキーワードに一致するドキュメントを検索します。
-   `read_document(id: string)`: ID を指定してドキュメントの完全な内容を取得します。

## 設定

`.kex.yaml`:

```yaml
root: contents # マークダウンファイルを含むディレクトリ
```

## ライセンス

MIT
