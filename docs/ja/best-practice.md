# ベストプラクティス

## ディレクトリ構造 (スコープ)

ドキュメントは階層的に整理してください。Kex はディレクトリ名を「スコープ」として使用します。

**推奨される階層**:
`Domain` (ドメイン) / `Platform` (プラットフォーム) / `Technology` (技術)

**例**:
- `coding/typescript/no-any.md` -> スコープ: `[coding, typescript]`
- `vcs/git/commit-style.md` -> スコープ: `[vcs, git]`
- `security/api/auth.md` -> スコープ: `[security, api]`

## アトミックなドキュメント

- **焦点を絞る**: 1ファイルにつき1つのルール。
- **「なぜ」を説明する**: AI にとって、構文よりも論拠 (Rationale) が重要です。
- **"Adopted" ステータスを使う**: チームが合意した時のみ `status: adopted` を設定してください。提案段階では `draft` を使用します。

## AI への指示

エージェント設定には、必ず「まず Kex に相談する (Consult Kex First)」というプロンプトを含めてください。詳細は `AGENTS.md` を参照してください。
