# ベストプラクティス

## ディレクトリ構造 (スコープ)

ドキュメントは階層的に整理してください。Kex はディレクトリ名を「スコープ」として扱います。

**推奨される階層**:
`Domain` / `Platform` / `Technology`

**例**:
- `coding/typescript/no-any.md` -> スコープ: `[coding, typescript]`
- `vcs/git/commit-style.md` -> スコープ: `[vcs, git]`
- `security/api/auth.md` -> スコープ: `[security, api]`

## アトミックなドキュメント

- **焦点を絞る**: 1ファイルにつき1つのルールに限定します。
- **「なぜ」を説明する**: AI にとって、構文そのものよりも、なぜそうするのかという論拠 (Rationale) が重要です。
- **"Adopted" ステータスを使う**: チームで合意が取れたものにのみ `status: adopted` を設定してください。提案段階のものは `draft` を使用します。

## キーワード (検索戦略)

Kex のキーワード検索は **完全一致 (Exact Match)** です。`tests` で検索しても、キーワード `test` しか持たないドキュメントはヒットしません。

- **単数形を使用する**: `tests`, `bugs` ではなく `test`, `bug` を優先してください。
- **同義語を明示的に追加する**: ひとつの概念が複数の用語で呼ばれる一般的な場合や、複数形での検索が強く想定される場合は、それらをすべてキーワードに追加してください。
    ```yaml
    keywords: ["test", "tests", "testing", "spec"]
    ```

## AI への指示

エージェント設定には、必ず「まず Kex に相談する (Consult Kex First)」という指示を含めてください。詳細は `AGENTS.md` を参照してください。
