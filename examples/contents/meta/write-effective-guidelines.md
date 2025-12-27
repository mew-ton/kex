---
id: write-effective-guidelines
title: Write effective guidelines for AI and Humans
description: >
  Standard for creating kex guidelines: structure, filename conventions, and keyword strategy for discoverability.
keywords:
  - documentation
  - guideline
  - meta
  - keyword-strategy
  - writing
sources:
  - name: Diátaxis (Reference)
    url: https://diataxis.fr/reference-documentation/
---

## Summary
Follow this standard to create guidelines that are easily discoverable by AI agents and useful for humans.

## Filename Convention
- Use **Imperative Mood**: `verb-adjective-object.md` (e.g., `keep-functions-short.md`, `handle-errors-explicitly.md`).
- Must match the `id` field in frontmatter.

## Keyword Strategy (The 3 Dimensions)
To ensure guidelines are found during various stages of development, provide keywords covering three dimensions:
1.  **Object**: The concrete entity being modified (e.g., `function`, `variable`, `file`, `import`).
2.  **Symptom/Problem**: What the user observes or suffers from (e.g., `length`, `complexity`, `large`, `bug`, `error`).
3.  **Concept/Quality**: The theoretical principle (e.g., `readability`, `maintainability`, `atomic-design`, `dry`).

## Structure
1.  **Summary**: 1-2 sentences. Use imperative command. "Do X to achieve Y."
2.  **Rationale**: Why is this important? "Without this, Z happens."
3.  **Guidance**: Actionable steps. Use numbered lists for sequences.
4.  **Examples**:
    - **Bad**: A concise anti-pattern.
    - **Good**: The corrected implementation.
