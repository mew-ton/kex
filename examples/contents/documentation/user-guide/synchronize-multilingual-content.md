---

title: Synchronize multilingual content
description: >
  When maintaining multilingual documentation, ensure content synchronization and prevent version drift.
keywords:
  - documentation
  - language
  - translation
  - sync
  - japanese
  - english
---

## Summary
Ensure that content in all supported languages is kept in sync. Updates to one language must be reflected in the others immediately.

## Rationale
- **Consistency**: Users should receive the same information regardless of language preference.
- **Trust**: Outdated translations erode trust in the documentation.
- **Maintainability**: Letting language versions drift apart makes future synchronization significantly harder and more error-prone.

## Guidance
1.  **Atomic Updates**: When modifying documentation, update all supported language files (e.g., `docs/feature-index.md` and `docs/ja/feature-index.md`) in the same commit.
2.  **No Lag**: Do not leave one language "to be updated later".
3.  **Verify**: Ensure the meaning and technical details are identical across versions.
