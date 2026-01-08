# Writing Documentation

## Format

- **File Extension**: `.md`
- **Language**: English is recommended for AI processing (even if humans read Japanese).

## Frontmatter Schema

Every document **MUST** start with YAML frontmatter.

```yaml
---
title: Avoid Magic Numbers # Required: Human readable title
description: Do not use magic numbers. # Required: Used by AI for selection
keywords: [readability, code-quality] # Required: Key terms for search
status: adopted # Optional: adopted | draft (Default: adopted)
sources: # Optional: External references
  - name: ESLint
    url: https://eslint.org/docs/rules/no-magic-numbers
---
```

### `title` (Required)

- **Type**: `string`
- **Description**: The human-readable title of the document.

### `description` (Required)

- **Type**: `string`
- **Description**: A brief summary of the document's content. This is used by the AI agent to select relevant guidelines.

### `keywords` (Required)

- **Type**: `string[]` (List of strings)
- **Description**: A set of keywords to facilitate search indexing.

### `status` (Optional)

- **Type**: `string`
- **Default**: `adopted`
- **Description**: The lifecycle status of the guideline.
  - `adopted`: The guideline is active and should be followed.
  - `draft`: The guideline is a work in progress.

### `sources` (Optional)

- **Type**: `object[]` (List of objects)
- **Description**: A list of external references or sources supporting the guideline.
  - `name`: Name of the source (e.g., "ESLint").
  - `url`: URL to the source.

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
