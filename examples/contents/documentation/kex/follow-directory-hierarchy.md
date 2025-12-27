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
We **recommend** organizing documents by **Domain**, **Platform** (Optional), and **Technology** (Optional).

This is the default convention for Kex, but you are free to define your own structure (up to any depth). Folder names act as scopes.

## Rationale
- **Discoverability**: Predictable paths make it easier to find relevant guidelines alongside code.
- **Scoping**: The directory structure defines the scope (context) in which the rule applies (e.g., `coding/typescript`).

## Guidance (Recommended Default)
1.  **Domain**: Start with the broad category (e.g., `coding`, `vcs`, `documentation`).
2.  **Platform**: Add execution environment if specific (e.g., `frontend`, `mobile`). *Optional*.
3.  **Technology**: Add specific language or tool (e.g., `typescript`, `git`, `react`). *Optional*.
4.  **File**: The guideline file itself (`id.md`).

## Examples

### Good
- `coding/keep-functions-short.md` (Domain only)
- `coding/typescript/no-any.md` (Domain + Tech)
- `coding/frontend/react/use-hooks.md` (Domain + Platform + Tech)
- `vcs/git/conventional-commits.md` (Domain + Tech)
- `documentation/kex/use-imperative-filenames.md`

### Bad
- `universal/coding/rule.md` (Use `coding/rule.md`)
- `typescript/rule.md` (Missing Domain `coding`)
- `react/rule.md` (Missing Domain/Platform)
