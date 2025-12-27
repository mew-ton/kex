---
id: follow-directory-hierarchy
title: Folder Structure
description: >
  Organize documents by Domain > Platform (Optional) > Technology.
keywords:
  - directory
  - structure
  - hierarchy
  - folder
---

## Summary
Organize documents by **Domain**, **Platform** (Optional), and **Technology**.

## Rationale
- **Discoverability**: Predictable paths make it easier to find relevant guidelines alongside code.
- **Scoping**: The directory structure defines the scope (context) in which the rule applies (e.g., `coding/typescript`).

## Guidance
1.  **Domain**: Start with the broad category (e.g., `coding`, `vcs`, `documentation`).
2.  **Platform**: Add execution environment if specific (e.g., `frontend`, `mobile`). *Optional*.
3.  **Technology**: End with the specific language or tool (e.g., `typescript`, `git`, `react`).
4.  **File**: The guideline file itself (`id.md`).

## Examples

### Good
- `coding/typescript/no-any.md`
- `coding/frontend/react/use-hooks.md`
- `vcs/git/conventional-commits.md`
- `documentation/kex/use-imperative-filenames.md`

### Bad
- `universal/coding/rule.md` (Use `coding/rule.md`)
- `typescript/rule.md` (Missing Domain `coding`)
- `react/rule.md` (Missing Domain/Platform)
