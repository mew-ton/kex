# ドキュメントの執筆

## フォーマット

- **拡張子**: `.md`
- **言語**: AI の処理精度のため、英語が推奨されます (人間が日本語を読む場合でも)。

## Frontmatter スキーマ

すべてのドキュメントは YAML frontmatter で開始する**必要**があります。

```yaml
---
id: avoid-magic-numbers  # 必須: ファイル名と一致させる必要があります (avoid-magic-numbers.md)
title: Avoid Magic Numbers # 必須
description: Do not use magic numbers. # 必須: AI がドキュメントを選択するために使用します
keywords: [readability, code-quality] # 必須: 検索用のキーワード
status: adopted # オプション: draft | adopted (デフォルト: adopted)
sources: # オプション: 外部参照へのリンク
  - name: ESLint
    url: https://eslint.org/docs/rules/no-magic-numbers
---
```

## コンテンツ構造

一貫性のため、以下の構造を推奨します:

```markdown
## Summary
概要。

## Rationale
なぜこのルールが存在するのか。

## Guidance
どのようにルールに従うべきか。

## Examples
良いコード例と悪いコード例。
```
