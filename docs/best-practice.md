# Best Practices

## Directory Structure (Scopes)

Organize your documents hierarchically. Kex uses directory names as "Scopes".

**Recommended Hierarchy**:
`Domain` / `Platform` / `Technology`

**Examples**:
- `coding/typescript/no-any.md` -> Scopes: `[coding, typescript]`
- `vcs/git/commit-style.md` -> Scopes: `[vcs, git]`
- `security/api/auth.md` -> Scopes: `[security, api]`

## Atomic Documents

- **Keep it focused**: One rule per file.
- **Explain "Why"**: Rationale is more important for AI than syntax.
- **Use "Adopted" status**: Only set `status: adopted` when the team agrees. Use `draft` for proposals.

## Keywords (Search Strategy)

Kex uses **Exact Match** for keywords. Searching for `tests` will NOT match a document with keyword `test`.

- **Use Singular Form**: Prefer `test`, `bug`, `function` over `tests`, `bugs`, `functions`.
- **Add Synonyms Explicitly**: If a concept is commonly referred to by multiple terms (or plural forms are highly relevant), add them all.
    ```yaml
    keywords: ["test", "tests", "testing", "spec"]
    ```

## AI Instructions

Always include a prompt in your agent configuration telling it to Consult Kex First. See `AGENTS.md` for details.
