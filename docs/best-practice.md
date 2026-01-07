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

## AI Instructions

Always include a prompt in your agent configuration telling it to Consult Kex First. See `AGENTS.md` for details.
