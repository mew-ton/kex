---
id: include-valid-frontmatter
title: Include valid frontmatter
description: >
  All guideline documents must have valid YAML frontmatter with required fields.
keywords:
  - frontmatter
  - yaml
  - metadata
  - id
  - title
  - description
---

## Summary
Every guideline document must begin with YAML frontmatter containing required metadata fields.

## Rationale
- **Discoverability**: Metadata enables efficient search and filtering by AI and humans.
- **Validation**: Frontmatter allows automated checks for document consistency.
- **Context**: Description and keywords provide quick understanding without reading full content.

## Guidance
1. **Required fields**: `id`, `title`, `description`
2. **Optional fields**: `keywords`
3. **ID matching**: The `id` must match the filename (without `.md` extension).
4. **Format**: Use YAML syntax enclosed by `---` delimiters.

## Examples

```markdown
---
id: keep-functions-short
title: Keep functions short
description: >
  Functions should be concise and focused on a single responsibility.
keywords:
  - function
  - readability
  - maintainability
---
```

## Field Specifications

- **id**: Lowercase kebab-case matching filename (e.g., `avoid-using-any` for `avoid-using-any.md`)
- **title**: Human-readable title in sentence case
- **description**: Brief summary (can use `>` for multiline)
- **keywords**: Array of searchable terms for discoverability
