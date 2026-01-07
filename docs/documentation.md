# Writing Documentation

## Format

- **File Extension**: `.md`
- **Language**: English is recommended for AI processing (even if humans read Japanese).

## Frontmatter Schema

Every document **MUST** start with YAML frontmatter.

```yaml
---
id: avoid-magic-numbers  # Required: Must match filename (avoid-magic-numbers.md)
title: Avoid Magic Numbers # Required
description: Do not use magic numbers. # Required: Used by AI to select docs
keywords: [readability, code-quality] # Required: Key terms for search
status: adopted # Optional: draft | adopted (Default: adopted)
sources: # Optional: Links to external references
  - name: ESLint
    url: https://eslint.org/docs/rules/no-magic-numbers
---
```

## Content Structure

We recommend the following structure for consistency:

```markdown
## Summary
Brief explanation.

## Rationale
Why this rule exists.

## Guidance
How to follow the rule.

## Examples
Positive and negative code examples.
```
