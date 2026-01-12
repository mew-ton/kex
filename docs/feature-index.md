# Features

Kex provides powerful context retrieval capabilities through **Keywords** and **Scopes**.

## Keywords

Keywords are explicit metadata defined in the frontmatter of your markdown documents. They serve as the primary index for Kex.

### How it works

1.  **Define**: Add relevant tags to the `keywords` field in your document's frontmatter.
2.  **Index**: Kex builds an inverted index maping keywords to documents.
3.  **Search**: When an AI agent (or user) searches for a term, Kex performs a fuzzy/exact match against these keywords.

### Example

```yaml
---
id: "BUTTON_COMPONENT"
title: "Button Component"
keywords: ["ui", "component", "button", "frontend", "design-system"]
---
```

> **Tip (Vibe Coding)**: Use a mix of technical terms ("react", "typescript") and conceptual terms ("auth", "payment") to ensure documents are found in various contexts.

## Scopes

Scopes are implicit context markers derived automatically from your directory structure. They are used to **filter** search results based on the current working context.

### How it works

1.  **Directory Structure**: Organize your `contents/` directory logically.
2.  **Derivation**: Kex treats every directory name in the path as a "scope".
3.  **Filtering**: When searching, if a scope is provided (e.g., derived from the file the user is currently editing), Kex filters the results to only include documents that share at least one scope.

### Example

Given a file at `contents/frontend/components/Button.md`:

- **Path**: `contents/frontend/components/Button.md`
- ** derived Scopes**: `["frontend", "components"]`

If the user is editing a file in `src/frontend/login.tsx` (which implies scope `frontend`), Kex will prioritize/allow documents that have the `frontend` scope.

### Subset Filtering

Search results are strictly filtered using **Subset Filtering**. A document matches only if:
- It matches a **Keyword**.
- AND **all of its Scopes** are present in the query context (or specified scopes).

For example, `contents/frontend/components/Button.md` (Scopes: `[frontend, components]`) will NOT match if the context is just `frontend`. It requires both `frontend` AND `components` in the context. This prevents documents from specific sub-contexts invading broader search contexts.
