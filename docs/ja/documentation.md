# ドキュメントの執筆

## フォーマット

- **拡張子**: `.md`
- **言語**: AI の処理精度向上のため、たとえ対象読者が日本人であっても、内容は **英語** で記述することを強く推奨します。

## Frontmatter スキーマ

すべてのドキュメントは、YAML形式の Frontmatter から開始する必要があります。

```yaml
---
title: Avoid Magic Numbers # 必須
description: Do not use magic numbers. # 必須: AI がドキュメントを選択するために使用します
keywords: [readability, code-quality] # 必須: 検索用のキーワード
status: adopted # 任意: draft | adopted (デフォルト: adopted)
sources: # 任意: 外部参照へのリンク
  - name: ESLint
    url: https://eslint.org/docs/rules/no-magic-numbers
---
```

### `title` (必須)

- **型**: `string`
- **説明**: 人間が読むためのドキュメントのタイトルです。

### `description` (必須)

- **型**: `string`
- **説明**: ドキュメント内容の簡潔な要約です。これは、AI エージェントが関連するガイドラインを選定する際の重要な判断材料となります。
- **ガイドライン**: **"When [context], use this to [goal]"** 形式を使用してください。
    - **Bad**: `このドキュメントは命名規則についてです。`
    - **Good**: `When naming variables, use this to ensure consistency.`

### `keywords` (必須)

- **型**: `string[]` (文字列のリスト)
- **説明**: 検索インデックスの精度を高めるためのキーワードのセットです。

### `status` (任意)

- **型**: `string`
- **デフォルト**: `adopted`
- **説明**: ガイドラインのライフサイクルステータスです。
  - `adopted`: ガイドラインは有効であり、遵守する必要があります。
  - `draft`: ガイドラインは作成中、または提案段階です。

### `sources` (任意)

- **型**: `object[]` (オブジェクトのリスト)
- **説明**: ガイドラインの根拠となる外部参照や情報ソースのリストです。
  - `name`: ソースの名前 (例: "ESLint")。
  - `url`: ソースへの URL。

## コンテンツ構造

一貫性を保つため、以下の構造を推奨します:

```markdown
## Summary
概要。

## Rationale
なぜこのルールが存在するのか (理由・論拠)。

## Guidance
どのようにルールに従うべきか (具体的な指示)。

## Examples
良いコード例と悪いコード例。
```
